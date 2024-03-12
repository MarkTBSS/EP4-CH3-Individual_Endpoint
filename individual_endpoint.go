package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "AnuchitO", Age: 18},
}

func usersHandler(write http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		log.Println("GET")
		b, err := json.Marshal(users)
		if err != nil {
			write.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(write, "error: %v", err)
			return
		}
		write.Header().Set("Content-Type", "application/json")
		write.Write(b)
	}
}

func healthHandler(write http.ResponseWriter, request *http.Request) {
	write.WriteHeader(http.StatusOK)
	write.Write([]byte("OK"))
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(write, request)
		log.Printf("Server http middleware: %s %s %s %s",
			request.RemoteAddr,
			request.Method,
			request.URL,
			time.Since(start))
	})
}

func main() {
	http.HandleFunc("/users", logMiddleware(usersHandler))
	http.HandleFunc("/health", logMiddleware(healthHandler))

	log.Println("Server started at :2565")
	log.Fatal(http.ListenAndServe(":2565", nil))
	log.Println("bye bye!")
}
