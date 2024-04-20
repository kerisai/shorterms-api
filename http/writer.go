package http

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrNotFound            = errors.New("page not found")
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrInternalServerError = errors.New("something went wrong on our side")
)

func WriteHttpBodyJson(w http.ResponseWriter, status int, body any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		WriteHttpError(w, http.StatusInternalServerError, ErrInternalServerError)
	}
}

func WriteHttpError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
}

func WriteHttpErrors(w http.ResponseWriter, status int, errs map[string]error) {
	res := make(map[string]string)
	for field, err := range errs {
		res[field] = err.Error()
	}
	WriteHttpBodyJson(w, status, map[string]map[string]string{"message": res})
}

func WriteHttpInternalServerError(w http.ResponseWriter) {
	WriteHttpError(w, http.StatusInternalServerError, ErrInternalServerError)
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "OK!"})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteHttpError(w, http.StatusNotFound, ErrNotFound)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	WriteHttpError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
}
