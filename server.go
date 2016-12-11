package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	cache  *cache.Cache
	client *Client
	config *Configuration
	logger *log.Logger
}

func NewServer(conf *Configuration, logger *log.Logger) *Server {
	return &Server{
		cache:  cache.New(5*time.Minute, 30*time.Second),
		client: NewClient(conf),
		config: conf,
		logger: logger,
	}
}

func (s *Server) postStatsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		panic(err)
	}

	s.logger.Println("request:", req.URL)

	var stats *wordpressPostStats
	key := fmt.Sprintf("stats/post/%d", id)
	cachedStats, found := s.cache.Get(key)

	if found {
		s.logger.Println("cache hit")
		stats = cachedStats.(*wordpressPostStats)
	} else {
		s.logger.Println("cache missed")
		stats, err = s.client.getPostStats(id)
		if err != nil {
			// TODO: Better error handling
			s.logger.Println("error:", err)
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(struct{}{})
			return
		}
		s.cache.Set(key, stats, cache.DefaultExpiration)
	}

	s.logger.Println("views:", stats.Views)

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
