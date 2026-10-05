package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const notificationsTopic = "notifications"

var (
	valkeyClient  valkey.Client
	kafkaProducer sarama.SyncProducer
	tracer        = otel.Tracer("notiflex-api")

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "핸들러별 HTTP 응답 상태 코드 누적 카운트",
	}, []string{"path", "status"})
)

// statusRecorder는 핸들러가 WriteHeader로 응답한 실제 상태 코드를 가로채
// http_requests_total 메트릭에 반영하기 위한 http.ResponseWriter 래퍼다.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func instrument(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)
		httpRequestsTotal.WithLabelValues(path, strconv.Itoa(rec.status)).Inc()
	}
}

// initTracer는 OTEL_EXPORTER_OTLP_ENDPOINT가 설정된 경우 Tempo로 트레이스를 전송하는
// TracerProvider를 초기화한다. 엔드포인트가 없으면 트레이싱 없이 nil을 반환한다.
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		log.Printf("OTEL_EXPORTER_OTLP_ENDPOINT 미설정, 트레이싱 비활성화")
		return nil, nil
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("OTLP exporter 생성 실패: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("notiflex-api"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("OTel resource 생성 실패: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	log.Printf("트레이싱 활성화: OTEL_EXPORTER_OTLP_ENDPOINT=%s", endpoint)
	return tp, nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "healthHandler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "idHandler")
	defer span.End()

	id, err := incrCounter(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "id 생성 실패")
		http.Error(w, "id 생성 실패", http.StatusInternalServerError)
		return
	}
	podName := os.Getenv("HOSTNAME")

	publishNotification(ctx, id, podName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     podName,
		"source":  "ci-argocd-e2e-test",
		"variant": "canary-live-demo-v1",
	})
}

// incrCounter는 Valkey INCR 호출을 별도 span으로 감싸, 트레이스에서
// Valkey 응답 지연을 API 핸들러 처리 시간과 구분해서 볼 수 있게 한다.
func incrCounter(ctx context.Context) (int64, error) {
	ctx, span := tracer.Start(ctx, "valkey.incr")
	defer span.End()
	span.SetAttributes(attribute.String("db.system", "valkey"))

	return valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build()).ToInt64()
}

// publishNotification은 Kafka 메시지 발행을 별도 span으로 감싼다. 실패해도
// /id 응답 자체는 계속 진행되므로 에러는 span에 기록만 하고 그대로 반환하지 않는다.
func publishNotification(ctx context.Context, id int64, podName string) {
	_, span := tracer.Start(ctx, "kafka.produce")
	defer span.End()
	span.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination", notificationsTopic),
	)

	if _, _, err := kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: notificationsTopic,
		Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":%q}`, id, podName)),
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Kafka 메시지 전송 실패")
		log.Printf("Kafka 메시지 전송 실패: %v", err)
	}
}

func valkeyPassword() string {
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			return string(data)
		}
	}
	return os.Getenv("VALKEY_PASSWORD")
}

func newValkeyClient() (valkey.Client, error) {
	var client valkey.Client
	var err error
	for i := 0; i < 10; i++ {
		client, err = valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{os.Getenv("VALKEY_ADDR")},
			Password:    valkeyPassword(),
		})
		if err == nil {
			return client, nil
		}
		log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return nil, err
}

func newKafkaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0
	return cfg
}

func newKafkaProducer(broker string) (sarama.SyncProducer, error) {
	cfg := newKafkaConfig()
	cfg.Producer.Return.Successes = true

	var producer sarama.SyncProducer
	var err error
	for i := 0; i < 10; i++ {
		producer, err = sarama.NewSyncProducer([]string{broker}, cfg)
		if err == nil {
			return producer, nil
		}
		log.Printf("Kafka Producer 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return nil, err
}

func consumeNotifications(broker string) {
	cfg := newKafkaConfig()

	var consumer sarama.Consumer
	var err error
	for i := 0; i < 10; i++ {
		consumer, err = sarama.NewConsumer([]string{broker}, cfg)
		if err == nil {
			break
		}
		log.Printf("Kafka Consumer 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Printf("Kafka Consumer 생성 실패: %v", err)
		return
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(notificationsTopic)
	if err != nil {
		log.Printf("Kafka 파티션 목록 조회 실패: %v", err)
		return
	}

	done := make(chan struct{})
	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(notificationsTopic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Kafka 파티션(%d) 구독 실패: %v", partition, err)
			continue
		}
		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for msg := range pc.Messages() {
				log.Printf("Kafka 메시지 수신: partition=%d offset=%d value=%s", msg.Partition, msg.Offset, string(msg.Value))
			}
		}(partitionConsumer)
	}
	<-done
}

func main() {
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		log.Printf("OpenTelemetry 초기화 실패, 트레이싱 없이 계속 진행: %v", err)
	}
	if tp != nil {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tp.Shutdown(shutdownCtx); err != nil {
				log.Printf("TracerProvider shutdown 실패: %v", err)
			}
		}()
	}

	client, err := newValkeyClient()
	if err != nil {
		log.Fatalf("Valkey 연결 실패: %v", err)
	}
	valkeyClient = client
	defer valkeyClient.Close()

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	producer, err := newKafkaProducer(kafkaBroker)
	if err != nil {
		log.Fatalf("Kafka Producer 연결 실패: %v", err)
	}
	kafkaProducer = producer
	defer kafkaProducer.Close()

	go consumeNotifications(kafkaBroker)

	http.HandleFunc("/health", instrument("/health", healthHandler))
	http.HandleFunc("/id", instrument("/id", idHandler))
	http.Handle("/metrics", promhttp.Handler())

	port := "8080"
	fmt.Printf("notiflex-api listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
