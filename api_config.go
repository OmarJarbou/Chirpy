package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"encoding/json"

	"github.com/OmarJarbou/Chirpy/internal/database"
)

// will hold any stateful, in-memory data we'll need to keep track of.
type apiConfig struct {
	fileserverHits  atomic.Int32 // atomic.Int32 type is a really cool standard-library type that allows us to safely increment and read an integer value across multiple goroutines (HTTP requests).
	DBQueries       *database.Queries
	ChirpySecretKey string
}

type resetSuccessResponseBody struct {
	Message string `json:"message"`
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
	errResBody := errorResponseBody{}
	var jsonResBody []byte
	err := cfg.DBQueries.DeleteAllUsers(req.Context())
	if err != nil {
		errResBody.Error = "Error while deleting users from database: " + err.Error()
		jsonResBody, err2 := json.Marshal(errResBody)
		writeJSONResponse(response_writer, jsonResBody, err2, 500)
		return
	}

	cfg.fileserverHits = atomic.Int32{}
	resetSuccessResBody := resetSuccessResponseBody{
		Message: "File server hits has been reset successfully",
	}
	jsonResBody, err3 := json.Marshal(resetSuccessResBody)
	writeJSONResponse(response_writer, jsonResBody, err3, 200)
}
