package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	endpoint = "http://localhost:8080/chat/messages"
	duration = 10 * time.Second
)

type ChatRequest struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	Content    string `json:"content"`
}

type Stats struct {
	TotalRequests   int64
	SuccessRequests int64
	ErrorRequests   int64
	TotalLatency    int64 // microseconds
	MinLatency      int64 // microseconds
	MaxLatency      int64 // microseconds
}

var (
	stats Stats
	mu    sync.Mutex
)

func main() {
	fmt.Println("🚀 Starting load test...")
	fmt.Printf("📡 Endpoint: %s\n", endpoint)
	fmt.Printf("⏱️  Duration: %v\n", duration)
	fmt.Println("🔥 Spamming with maximum concurrency...")
	fmt.Println()

	// Initialize min latency to a large value
	atomic.StoreInt64(&stats.MinLatency, 999999999)

	startTime := time.Now()
	endTime := startTime.Add(duration)

	// Use a wait group to track all goroutines
	var wg sync.WaitGroup

	// Channel to limit concurrency (optional - remove if you want unlimited)
	// For maximum load, we'll use a high number of concurrent workers
	maxWorkers := 1000 // Adjust based on your machine
	semaphore := make(chan struct{}, maxWorkers)

	// Start workers
	workerCount := 0
	for time.Now().Before(endTime) {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		go func(id int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			for time.Now().Before(endTime) {
				sendRequest(id)
			}
		}(workerCount)
		workerCount++
		// Small delay to avoid creating too many goroutines at once
		time.Sleep(1 * time.Millisecond)
	}

	// Wait for all workers to finish
	wg.Wait()

	// Calculate final stats
	total := atomic.LoadInt64(&stats.TotalRequests)
	success := atomic.LoadInt64(&stats.SuccessRequests)
	errors := atomic.LoadInt64(&stats.ErrorRequests)
	totalLatency := atomic.LoadInt64(&stats.TotalLatency)
	minLatency := atomic.LoadInt64(&stats.MinLatency)
	maxLatency := atomic.LoadInt64(&stats.MaxLatency)

	avgLatency := float64(0)
	if success > 0 {
		avgLatency = float64(totalLatency) / float64(success) / 1000.0 // Convert to milliseconds
	}

	actualDuration := time.Since(startTime)
	reqPerSec := float64(total) / actualDuration.Seconds()
	successPerSec := float64(success) / actualDuration.Seconds()

	fmt.Println()
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 60)) + "=")
	fmt.Println("📊 LOAD TEST RESULTS")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 60)) + "=")
	fmt.Printf("⏱️  Duration:        %v\n", actualDuration.Round(time.Millisecond))
	fmt.Printf("👥 Workers:         %d\n", workerCount)
	fmt.Printf("📈 Total Requests:  %d\n", total)
	fmt.Printf("✅ Success:         %d (%.2f%%)\n", success, float64(success)/float64(total)*100)
	fmt.Printf("❌ Errors:          %d (%.2f%%)\n", errors, float64(errors)/float64(total)*100)
	fmt.Printf("🚀 Requests/sec:    %.2f req/s\n", reqPerSec)
	fmt.Printf("✅ Success/sec:     %.2f req/s\n", successPerSec)
	fmt.Printf("⚡ Avg Latency:     %.2f ms\n", avgLatency)
	if minLatency < 999999999 {
		fmt.Printf("🏃 Min Latency:     %.2f ms\n", float64(minLatency)/1000.0)
		fmt.Printf("🐌 Max Latency:     %.2f ms\n", float64(maxLatency)/1000.0)
	}
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 60)) + "=")
}

func sendRequest(workerID int) {
	// Create request payload
	reqBody := ChatRequest{
		SenderID:   fmt.Sprintf("user-%d", workerID%100),
		ReceiverID: fmt.Sprintf("user-%d", (workerID+1)%100),
		Content:    fmt.Sprintf("Load test message from worker %d at %d", workerID, time.Now().UnixNano()),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		atomic.AddInt64(&stats.ErrorRequests, 1)
		atomic.AddInt64(&stats.TotalRequests, 1)
		return
	}

	// Measure latency
	start := time.Now()
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	latency := time.Since(start).Microseconds()

	atomic.AddInt64(&stats.TotalRequests, 1)

	if err != nil {
		atomic.AddInt64(&stats.ErrorRequests, 1)
		return
	}
	defer resp.Body.Close()

	// Read response body (optional, but good practice)
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		atomic.AddInt64(&stats.SuccessRequests, 1)
		atomic.AddInt64(&stats.TotalLatency, latency)

		// Update min/max latency
		for {
			oldMin := atomic.LoadInt64(&stats.MinLatency)
			if latency >= oldMin {
				break
			}
			if atomic.CompareAndSwapInt64(&stats.MinLatency, oldMin, latency) {
				break
			}
		}

		for {
			oldMax := atomic.LoadInt64(&stats.MaxLatency)
			if latency <= oldMax {
				break
			}
			if atomic.CompareAndSwapInt64(&stats.MaxLatency, oldMax, latency) {
				break
			}
		}
	} else {
		atomic.AddInt64(&stats.ErrorRequests, 1)
	}
}

