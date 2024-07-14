package locations_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	server "github.com/alexPavlikov/gora_driver_location_service/cmd"
	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
)

func TestDriverPostCord(t *testing.T) {

	var cfg = config.Config{
		Env:     "local",
		Timeout: 5 * time.Second,
		Server: config.Server{
			Path: "localhost",
			Port: 8001,
		},
		LogLevel: "local",
		Kafka: config.Server{
			Path: "localhost",
			Port: 9092,
		},
		KafkaTopic: "drivers",
	}
	srv, close, err := server.ServerLoad(&cfg)
	if err != nil {
		t.Error(err)
	}

	defer close()

	recorder := httptest.NewRecorder()

	bodyReader := strings.NewReader(`{"longitude": 4040, "latitude": 3030}`)

	r := httptest.NewRequest("POST", "http://localhost:8001/v1/locations/5434", bodyReader)

	srv.ServeHTTP(recorder, r)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}
}
