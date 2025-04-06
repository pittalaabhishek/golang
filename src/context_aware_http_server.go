package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create a new HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/process", processHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server
	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// processHandler handles the /process endpoint with context-aware processing
func processHandler(w http.ResponseWriter, r *http.Request) {
	// Create a context with 5-second timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel() // Important to avoid context leaks

	// Channel to receive the processing result
	resultCh := make(chan string)

	// Start the long-running task in a goroutine
	go func() {
		resultCh <- longRunningTask(ctx)
	}()

	// Wait for either the task to complete or context to be cancelled
	select {
	case result := <-resultCh:
		// Task completed successfully
		fmt.Fprintf(w, "Processing completed: %s\n", result)
		log.Println("Request completed successfully")
	case <-ctx.Done():
		// Context was cancelled (timeout or client disconnect)
		err := ctx.Err()
		switch err {
		case context.DeadlineExceeded:
			msg := "Processing cancelled due to timeout"
			http.Error(w, msg, http.StatusGatewayTimeout)
			log.Println(msg)
		case context.Canceled:
			msg := "Processing cancelled by client"
			http.Error(w, msg, http.StatusBadRequest)
			log.Println(msg)
		default:
			msg := fmt.Sprintf("Processing cancelled: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			log.Println(msg)
		}
	}
}

// longRunningTask simulates a task that takes a long time to complete
func longRunningTask(ctx context.Context) string {
	log.Println("Long running task started")

	// Simulate work that takes 8 seconds (longer than our timeout)
	select {
	case <-time.After(2 * time.Second):
		return "Task result"
	case <-ctx.Done():
		log.Println("Task cancelled:", ctx.Err())
		return "" // This won't be used as context is already done
	}
}