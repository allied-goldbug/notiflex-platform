package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/valkey-io/valkey-go"
)

var valkeyClient valkey.Client

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     podName,
		"source":  "ci-argocd-e2e-test",
		"variant": "bluegreen-live-demo-v4",
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

func main() {
	client, err := newValkeyClient()
	if err != nil {
		log.Fatalf("Valkey 연결 실패: %v", err)
	}
	valkeyClient = client
	defer valkeyClient.Close()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)

	port := "8080"
	fmt.Printf("notiflex-api listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
