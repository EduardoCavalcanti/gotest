package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type APIResponse struct {
	Message string
}

const maxHandlers = 10

var sem = make(chan struct{}, maxHandlers)

func service(done chan string) {
	defer func() { <-sem }()

	// AQUI EU FAÇO A REQUEST
	time.Sleep(time.Second * 2)

	log.Println("inside service fn")

	done <- "API Response Goes Here"
}

func worker(ctx context.Context, done chan string) {
recvLoop:
	for {
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			break recvLoop
		}

		go service(done)
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	log.Println("Request accepted")
	done := make(chan string)
	go worker(ctx, done)

	response := APIResponse{
		Message: <-done,
	}

	log.Println("Request end")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", api)
	log.Fatal(http.ListenAndServe(":8002", nil))
}
