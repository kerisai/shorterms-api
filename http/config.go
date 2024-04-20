package http

import "github.com/kerisai/shorterms-api/config"

var allowedOrigins []string
var allowedMethods []string
var allowedHeaders []string

func Configure(c config.Config) {
	allowedOrigins = []string{c.ClientUrl}
	allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	if c.Env == "local" || c.Env == "development" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000", "http://localhost:3001")
		allowedMethods = append(allowedMethods, "HEAD", "TRACE")
	}

	allowedHeaders = []string{
		"Accept",
		"Authorization",
		"X-Forwarded-Authorization",
		"Content-Type",
	}
}
