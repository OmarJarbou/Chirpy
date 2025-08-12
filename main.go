package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

const PORT string = "8080"

func main() {
	serve_mux := http.NewServeMux()

	server := http.Server{
		Handler: serve_mux,
		Addr:    ":" + PORT,
	}

	api_config := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// option 1:
	// Adjust project's file structure to be the same as file server's path structure. (both index.html and assets folder moved to folder called app)
	// That way, you can access anything inside /app/ "ONLY" by writing "http://localhost:8080/app/[..path continue..]".
	// But you cant access any thing outside app/
	// serve_mux.Handle("/app/", http.FileServer(http.Dir("."))) // http.Dir(".") means "serve files from the current directory where the Go program is running (typically root).

	// option 2:
	// Keep program structure as is. (index.html at root, and assets at root)
	// That way, "/app" prefex will be stripped out of the recieved url:
	// "http://localhost:8080/app/[..path continue..]" (request url) => "http://localhost:8080/[..path continue..]" (stripped url)
	// That way, the stripped url will exactly point to what found in root
	serve_mux.Handle("/app/", api_config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// Why we realy want app/ at the path and changed from / to it? To solve handling/routing conflicts for the incoming req:
	// Initially, we had:
	// serve_mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	// This is a very broad rule: "If any request comes in that doesn't match
	// a more specific rule, send it to the FileServer." This means that any
	// request, including one for /healthz, would first hit this FileServer rule.
	// The FileServer would then look for a file named healthz in your filepathRoot.
	// Since there isn't one, it would likely return a 404 Not Found error.

	serve_mux.HandleFunc("/healthz", readinessHandler)

	serve_mux.HandleFunc("/metrics/", api_config.numberOfRequestsEncountered)
	serve_mux.HandleFunc("/reset/", api_config.resetFileServerHits)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error while listening on port '"+PORT+"':", err)
		os.Exit(1)
	}
}
