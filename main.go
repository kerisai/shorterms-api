package main

import (
	"github.com/kerisai/shorterms-api/config"
	"github.com/kerisai/shorterms-api/http"
)

func main() {
	config := config.LoadConfig()

	http.RunServer(config)
}