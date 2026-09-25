package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/valkey-io/valkey-go"
)

const notificationsTopic = "notifications"

var (
	valkeyClient  valkey.Client
	kafkaProducer sarama.SyncProducer
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build()).ToInt64()
	if err != nil {
		http.Error(w, "id 생성 실패", http.StatusInternalServerError)
		return
	}
	podName := os.Getenv("HOSTNAME")

	if _, _, err := kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: notificationsTopic,
		Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":%q}`, id, podName)),
	}); err != nil {
		log.Printf("Kafka 메시지 전송 실패: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     podName,
		"source":  "ci-argocd-e2e-test",
		"variant": "canary-live-demo-v1",
	})
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

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)

	port := "8080"
	fmt.Printf("notiflex-api listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
