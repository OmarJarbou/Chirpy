package main

import (
	"net/http"
)

type errorResponseBody struct {
	Error string `json:"error"`
}

func writeJSONResponse(response_writer http.ResponseWriter, jsonResBody []byte, err error, statusCode int) {
	if err != nil {
		response_writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response_writer.WriteHeader(500)
		response_writer.Write([]byte("Error while marshaling response body " + err.Error()))
		return
	}

	response_writer.Header().Set("Content-Type", "application/json")
	response_writer.WriteHeader(statusCode)
	response_writer.Write(jsonResBody)
}
