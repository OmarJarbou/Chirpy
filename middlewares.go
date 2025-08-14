package main

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response_writer http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(response_writer, req)
	})
}

func middlewareValidateChirp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response_writer http.ResponseWriter, req *http.Request) {
		reqBody := createChirpRequestBody{}
		errorResBody := errorResponseBody{}
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			errorResBody.Error = "Error while decoding request's json " + err.Error()
			jsonResBody, err2 := json.Marshal(errorResBody)

			writeJSONResponse(response_writer, jsonResBody, err2, 500)
			return
		}

		if len(reqBody.Body) > 140 {
			errorResBody.Error = "Chirp is too long, it must be 140 character long or less"
			jsonResBody, err3 := json.Marshal(errorResBody)

			writeJSONResponse(response_writer, jsonResBody, err3, 400)
			return
		}
		// good solution but it results in a string with all it's characters small:
		// filtered_chirp := strings.ReplaceAll(strings.ToLower(reqBody.Body), "kerfuffle", "****")
		// filtered_chirp = strings.ReplaceAll(filtered_chirp, "sharbert", "****")
		// filtered_chirp = strings.ReplaceAll(filtered_chirp, "fornax", "****")

		banned := []string{"kerfuffle", "sharbert", "fornax"}

		filtered_chirp := reqBody.Body
		for _, word := range banned {
			// `(?i)` makes the regex case-insensitive
			re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word))
			filtered_chirp = re.ReplaceAllString(filtered_chirp, "****")
		}

		ctx := context.WithValue(req.Context(), "filtered_chirp", filtered_chirp)
		ctx = context.WithValue(ctx, "user_id", reqBody.UserID)

		// if chirp is valid
		next.ServeHTTP(response_writer, req.WithContext(ctx))
	})
}
