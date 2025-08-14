package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type chirpValidatorRequestBody struct {
	Body string `json:"body"`
}

type errorResponseBody struct {
	Error string `json:"error"`
}

type chirpValidatorSuccessResponseBody struct {
	CleanedBody string `json:"cleaned_body"`
}

func chirpValidator(response_writer http.ResponseWriter, req *http.Request) {
	reqBody := chirpValidatorRequestBody{}
	errorResBody := errorResponseBody{}
	successResBody := chirpValidatorSuccessResponseBody{}
	var jsonResBody []byte
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

	successResBody.CleanedBody = filtered_chirp
	jsonResBody, err4 := json.Marshal(successResBody)

	writeJSONResponse(response_writer, jsonResBody, err4, 200)
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
