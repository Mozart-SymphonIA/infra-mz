package httpx

import (
	"fmt"
	"net/http"
)

const (
	contentTypePlain = "text/plain; charset=utf-8"
	contentTypeJSON  = "application/json; charset=utf-8"
	contentTypeTSV   = "text/plain; charset=utf-8"
)

func writeText(w http.ResponseWriter, status int, contentType, body string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func BadRequest(w http.ResponseWriter, message string) {
	writeText(w, http.StatusBadRequest, contentTypePlain, message)
}

func NotFound(w http.ResponseWriter, message string) {
	writeText(w, http.StatusNotFound, contentTypePlain, message)
}

func InternalError(w http.ResponseWriter, err error) {
	writeText(w, http.StatusInternalServerError, contentTypePlain, "An unexpected error occurred.")
	fmt.Println("error:", err)
}

func OkJSON(w http.ResponseWriter, body string) {
	writeText(w, http.StatusOK, contentTypeJSON, body)
}

func OkTSV(w http.ResponseWriter, body string) {
	writeText(w, http.StatusOK, contentTypeTSV, body)
}
