package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	serverPort = 9999
)

func NewFakeUserServer(port int) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		content := `
[
{"name": "user1","email": "user1@example.com"},
{"name": "user2","email": "user2@example.com"},
{"name": "user3","email": "user3@example.com"}
]
`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(content))
		log.Println("GET /api/users")
	}).Methods("GET")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	return srv
}

func main() {
	log.Println("Starting user double server")
	userSrv := NewFakeUserServer(serverPort)
	err := userSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error received: %v", err)
	}
	log.Println("Exiting user double server")
}
