package writer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CacheValue represents the structure to be stored in the cache
type CacheValue struct {
	ID string `json:"id"`
}

// Writer handles writing random UUID key-value pairs to birb-nest
type Writer struct {
	baseURL    string
	httpClient *http.Client
	interval   time.Duration
	keys       []string
	keysMu     sync.Mutex
}

// New creates a new Writer instance
func New(baseURL string, interval time.Duration) *Writer {
	return &Writer{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		interval: interval,
		keys:     make([]string, 0, 100),
	}
}

// Start begins the write loop, continuously writing random UUIDs to birb-nest
// This method blocks until the context is cancelled
func (w *Writer) Start(ctx context.Context) error {
	log.Printf("Starting UUID writer with %v interval", w.interval)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Write immediately on start
	if err := w.performRandomOperation(ctx); err != nil {
		log.Printf("Warning: initial operation failed: %v", err)
	}

	// Continue operations on interval
	for {
		select {
		case <-ctx.Done():
			log.Println("Writer stopped by context cancellation")
			return ctx.Err()
		case <-ticker.C:
			if err := w.performRandomOperation(ctx); err != nil {
				log.Printf("Warning: operation failed: %v", err)
			}
		}
	}
}

// performRandomOperation randomly decides to either write or read (if keys exist)
func (w *Writer) performRandomOperation(ctx context.Context) error {
	w.keysMu.Lock()
	hasKeys := len(w.keys) > 0
	w.keysMu.Unlock()

	// If we have keys, randomly choose between read and write
	if hasKeys && rand.Intn(2) == 0 {
		return w.readRandomKey(ctx)
	}
	return w.writeRandomUUID(ctx)
}

// addKey adds a key to the slice, maintaining max size of 100
func (w *Writer) addKey(key string) {
	w.keysMu.Lock()
	defer w.keysMu.Unlock()

	w.keys = append(w.keys, key)
	if len(w.keys) > 100 {
		w.keys = w.keys[1:]
	}
}

// getRandomKey returns a random key from the slice
func (w *Writer) getRandomKey() string {
	w.keysMu.Lock()
	defer w.keysMu.Unlock()

	if len(w.keys) == 0 {
		return ""
	}
	return w.keys[rand.Intn(len(w.keys))]
}

// writeRandomUUID generates two random UUIDs and writes them to birb-nest
func (w *Writer) writeRandomUUID(ctx context.Context) error {
	// Generate random UUID for key
	key := uuid.New().String()

	// Generate random UUID for value and create CacheValue struct
	cacheValue := CacheValue{
		ID: uuid.New().String(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(cacheValue)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	// Write to birb-nest using POST /v1/cache/:key
	url := fmt.Sprintf("%s/v1/cache/%s", w.baseURL, key)
	log.Printf("Writing: key=%s value=%s", key, string(jsonData))

	// Create request with JSON as body
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully wrote key=%s", key)

	// Add key to our slice
	w.addKey(key)

	return nil
}

// readRandomKey reads a random key from our stored keys
func (w *Writer) readRandomKey(ctx context.Context) error {
	key := w.getRandomKey()
	if key == "" {
		return fmt.Errorf("no keys available to read")
	}

	// Read from birb-nest using GET /v1/cache/:key
	url := fmt.Sprintf("%s/v1/cache/%s", w.baseURL, key)
	log.Printf("Reading: key=%s", key)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the value
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Successfully read key=%s value=%s", key, string(body))
	return nil
}

// Close closes the HTTP client
func (w *Writer) Close() error {
	log.Println("Closing writer")
	w.httpClient.CloseIdleConnections()
	return nil
}