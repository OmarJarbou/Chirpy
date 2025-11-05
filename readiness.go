package main

import (
	"net/http"
)

func readinessHandler(response_writer http.ResponseWriter, req *http.Request) {
	response_writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response_writer.WriteHeader(200)
	response_writer.Write([]byte("OK"))
}
