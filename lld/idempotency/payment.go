/*
In-memory implementation
*/
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

type PaymentRequest struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}

type PaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Amount    int64  `json:"amount"`
}

type Record struct {
	RequestHash string
	Status      string
	Response    []byte
	StatusCode  int
}

type IdempotencyStore struct {
	mu      sync.Mutex
	records map[string]*Record
}

func NewStore() *IdempotencyStore {
	return &IdempotencyStore{records: make(map[string]*Record)}
}

func (s *IdempotencyStore) Reserve(key, reqHash string) (*Record, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if record, exists := s.records[key]; exists {
		copy := *record
		copy.Response = append([]byte(nil), record.Response...)
		return &copy, false
	}
	s.records[key] = &Record{
		RequestHash: reqHash,
		Status:      "IN_PROGRESS",
	}
	return nil, true
}

func (s *IdempotencyStore) Complete(
	key string, response []byte, statusCode int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.records[key]
	record.Status = "COMPLETED"
	record.Response = append([]byte(nil), response...)
	record.StatusCode = statusCode
}

var PaymentCounter atomic.Uint64

func ProcessPayment(req PaymentRequest) PaymentResponse {
	id := PaymentCounter.Add(1)
	return PaymentResponse{
		PaymentID: fmt.Sprintf("pay_%d", id),
		Status:    "SUCCESS",
		Amount:    req.Amount,
	}
}

func RequestHash(req PaymentRequest) string {
	data, _ := json.Marshal(req)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func PaymentHandler(s *IdempotencyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		key := r.Header.Get("Idempotency-key")
		if key == "" || len(key) > 255 {
			http.Error(w, "invalid idempotency key", http.StatusBadRequest)
			return
		}

		var req PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if req.UserID == "" || req.Amount < 0 {
			http.Error(w, "invalid payment request", http.StatusBadRequest)
			return
		}

		scopedKey := req.UserID + ":" + key
		hash := RequestHash(req)

		record, reserved := s.Reserve(scopedKey, hash)
		if !reserved {
			if record.RequestHash != hash {
				http.Error(w, "idempotent key reused with different payload", http.StatusConflict)
				return
			}
			if record.Status == "IN_PROGRESS" {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "payment processing", http.StatusConflict)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(record.StatusCode)
			w.Write(record.Response)
			return
		}
		result := ProcessPayment(req)
		body, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "response encoding failed", http.StatusInternalServerError)
			return
		}
		s.Complete(scopedKey, body, http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	}
}

func main() {
	store := NewStore()
	http.HandleFunc("/payments", PaymentHandler(store))
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
curl -i -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: payment-123" \
  -d '{"user_id":"user-1","amount":1000}'
*/
