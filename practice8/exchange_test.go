package practice8

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRateSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/convert" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		if from != "USD" || to != "EUR" {
			t.Fatalf("unexpected query params: from=%s to=%s", from, to)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.92}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	rate, err := service.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rate != 0.92 {
		t.Fatalf("expected 0.92, got %v", rate)
	}
}

func TestGetRateAPIBusinessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"invalid currency pair"}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("AAA", "BBB")
	if err == nil {
		t.Fatal("expected api error, got nil")
	}
}

func TestGetRateMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestGetRateServerError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":"internal server error"}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected server error, got nil")
	}
}

func TestGetRateEmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error for empty body, got nil")
	}
}

func TestGetRateTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.92}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	service.Client.Timeout = 50 * time.Millisecond

	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected timeout/network error, got nil")
	}
}
