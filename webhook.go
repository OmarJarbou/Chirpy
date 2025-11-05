package main

import (
	"encoding/json"
	"net/http"

	"github.com/OmarJarbou/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type webhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) webhookHandler(response_writer http.ResponseWriter, req *http.Request) {
	reqBody := webhookRequest{}
	errorResBody := errorResponseBody{}

	api_key, err := auth.GetAPIKey(req.Header)
	if err != nil {
		errorResBody.Error = err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)

		writeJSONResponse(response_writer, jsonResBody, err2, 401)
		return
	}

	if cfg.PolkaKey != api_key {
		errorResBody.Error = "You are not allowed to do this action!"
		jsonResBody, err3 := json.Marshal(errorResBody)

		writeJSONResponse(response_writer, jsonResBody, err3, 401)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		errorResBody.Error = "Error while decoding request's json " + err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)

		writeJSONResponse(response_writer, jsonResBody, err2, 500)
		return
	}

	switch reqBody.Event {
	case "user.upgraded":
		{
			user_uuid, err3 := uuid.Parse(reqBody.Data.UserID)
			if err3 != nil {
				errorResBody.Error = "Error while parsing user id to uuid: " + err3.Error()
				jsonResBody, err4 := json.Marshal(errorResBody)

				writeJSONResponse(response_writer, jsonResBody, err4, 400)
				return
			}

			_, err5 := cfg.DBQueries.UpgradeUserTOChirpyRed(req.Context(), user_uuid)
			if err5 != nil {
				errorResBody.Error = "Error while upgrading user to chirpy red: " + err5.Error()
				jsonResBody, err6 := json.Marshal(errorResBody)

				writeJSONResponse(response_writer, jsonResBody, err6, 404)
				return
			}

			response_writer.WriteHeader(204)
		}
	default:
		{
			response_writer.WriteHeader(204)
			return
		}
	}
}
