package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	serverPort = 10000
)

func NewNotifServer(port int) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/api/notif", func(w http.ResponseWriter, r *http.Request) {
		buff := []byte(`{"status":"error", "message":"Not implemented"}`)
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Content-Type", "application/json")
		w.Write(buff)
	}).Methods(http.MethodPost)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	return srv
}

func main() {
	log.Println("Starting user notification service")
	userSrv := NewNotifServer(serverPort)
	err := userSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error received: %v", err)
	}
	log.Println("Exiting user notification service")
}
