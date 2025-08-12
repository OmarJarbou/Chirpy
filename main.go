package main

import (
	"fmt"
	"net/http"
)

const PORT string = "8080"

func main() {
	serve_mux := http.NewServeMux()

	server := http.Server{
		Handler: serve_mux,
		Addr:    ":" + PORT,
	}

	serve_mux.Handle("/", http.FileServer(http.Dir(".")))

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error while listening on port '"+PORT+"':", err)
	}
}
