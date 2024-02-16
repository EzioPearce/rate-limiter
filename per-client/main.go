package main

import (
	"encoding/json"
	"net/http"
)

type message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func perClientrateLimiter() http.Handler {
	clients := make(map[string]*rateLimiter)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if _, found := clients[clientIP]; !found {
			clients[clientIP] = newRateLimiter()
		}
		clients[clientIP].rateLimit(w, r)
	})
}

func endpointHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "Hello There",
	}
	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		return
	}
	//writer.Write([]byte(`{"status":"ok","body":"Hello, World!"}`))
}

func main() {
	http.Handle("/ping", perClientrateLimiter(endpointHandler))
	http.ListenAndServe(":8080", nil)
}
