package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/consul-client/client"

	"github.com/gorilla/mux"
)

// RunServer Creates http server to subscribe to consul
func RunServer(server client.Service) {
	r := mux.NewRouter()
	r.HandleFunc(server.HealthCheck, func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods("GET")
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", server.URL, server.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
