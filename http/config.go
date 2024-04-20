package http

var allowedOrigins []string
var allowedMethods []string
var allowedHeaders []string

func Configure(clientUrl, env string) {
	allowedOrigins = []string{clientUrl}
	allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	if env == "local" || env == "development" {
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
