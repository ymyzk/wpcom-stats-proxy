package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	client *Client
	config *Configuration
}

func NewServer(conf *Configuration) *Server {
	return &Server{
		client: NewClient(conf),
		config: conf,
	}
}

func (s *Server) postStatsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		panic(err)
	}

	stats := wordpressPostStats{}
	err = s.client.getPostStats(&stats, id)
	if err != nil {
		// TODO: Better error handling
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(struct{}{})
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (s *Server) start() {
	r := mux.NewRouter()
	r.HandleFunc("/stats/post/{id:[0-9]+}", s.postStatsHandler)

	http.Handle("/", r)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
