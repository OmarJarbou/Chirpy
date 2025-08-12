package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// will hold any stateful, in-memory data we'll need to keep track of.
type apiConfig struct {
	fileserverHits atomic.Int32 // atomic.Int32 type is a really cool standard-library type that allows us to safely increment and read an integer value across multiple goroutines (HTTP requests).
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response_writer http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(response_writer, req)
	})
}

func (cfg *apiConfig) numberOfRequestsEncountered(response_writer http.ResponseWriter, req *http.Request) {
	response_writer.Header().Set("Content-Type", "text/html")
	response_writer.WriteHeader(200)
	response_writer.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetFileServerHits(response_writer http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
	response_writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response_writer.WriteHeader(200)
	response_writer.Write([]byte("File server hits has been reset successfully"))
}
