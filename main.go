package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	ClientID     string `required:"true" split_words:"true"`
	ClientSecret string `required:"true" split_words:"true"`
	RedirectURL  string `split_words:"true"`
	Token        string `required:"true"`
	Site         string `required:"true"`
	Host         string
	Port         int `default:"8000"`
}

// https://developer.wordpress.com/docs/api/1.1/get/sites/%24site/stats/post/%24post_id/
type wordpressPostStats struct {
	Views int `json:"views"`
}

func main() {
	var conf Configuration
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	err := envconfig.Process("WPCOM_STATS_PROXY", &conf)
	if err != nil {
		panic(err)
	}
	s := NewServer(&conf, logger)
	s.start()
}
