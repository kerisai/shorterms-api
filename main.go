package main

import (
	"github.com/kerisai/shorterms-api/config"
	"github.com/kerisai/shorterms-api/db"
	"github.com/kerisai/shorterms-api/http"
)

func main() {
	config := config.LoadConfig()
	_ = db.CreateConnPool(config)

	http.RunServer(config)
}
