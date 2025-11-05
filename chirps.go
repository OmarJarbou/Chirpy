package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/OmarJarbou/Chirpy/internal/database"
	"github.com/google/uuid"
)

type createChirpRequestBody struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(response_writer http.ResponseWriter, req *http.Request) {
	errorResBody := errorResponseBody{}
	var jsonResBody []byte
	id := req.Context().Value("user_id").(uuid.UUID)

	createChirpParams := database.CreateChirpParams{
		Body:   req.Context().Value("filtered_chirp").(string),
		UserID: id,
	}

	chirp, err3 := cfg.DBQueries.CreateChirp(req.Context(), createChirpParams)
	if err3 != nil {
		errorResBody.Error = "Error while creating chirp: " + err3.Error()
		jsonResBody, err4 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err4, 500)
		return
	}

	successResBody := Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	jsonResBody, err5 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err5, 201)
}

func (cfg *apiConfig) handleGetAllChirps(response_writer http.ResponseWriter, req *http.Request) {
	errorResBody := errorResponseBody{}
	var jsonResBody []byte

	var chirps []database.Chirp
	var err error
	author_id := req.URL.Query().Get("author_id")
	sortq := req.URL.Query().Get("sort")
	if author_id == "" {
		chirps, err = cfg.DBQueries.GetAllChirps(req.Context())
		if err != nil {
			errorResBody.Error = "Error while fetching chirps from database: " + err.Error()
			jsonResBody, err2 := json.Marshal(errorResBody)
			writeJSONResponse(response_writer, jsonResBody, err2, 500)
			return
		}
	} else {
		author_uuid, err3 := uuid.Parse(author_id)
		if err3 != nil {
			errorResBody.Error = "Invalid UUID: " + err3.Error()
			jsonResBody, err4 := json.Marshal(errorResBody)
			writeJSONResponse(response_writer, jsonResBody, err4, 401)
			return
		}
		chirps, err = cfg.DBQueries.GetAllChirpsForAUser(req.Context(), author_uuid)
		if err != nil {
			errorResBody.Error = "Error while fetching chirps for a user from database: " + err.Error()
			jsonResBody, err5 := json.Marshal(errorResBody)
			writeJSONResponse(response_writer, jsonResBody, err5, 500)
			return
		}
	}
	if sortq == "asc" {
		sort.SliceStable(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else if sortq == "desc" {
		sort.SliceStable(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	successResBody := []Chirp{}

	for _, chirp := range chirps {
		successResBody = append(successResBody, Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		})
	}

	jsonResBody, err3 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err3, 200)
}

func (cfg *apiConfig) handleGetChirpByID(response_writer http.ResponseWriter, req *http.Request) {
	chirp_id := req.PathValue("chirpID")
	errorResBody := errorResponseBody{}
	var jsonResBody []byte
	parsed_chirp_id, err := uuid.Parse(chirp_id)
	if err != nil {
		errorResBody.Error = "Invalid UUID: " + err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err2, 400)
		return
	}

	chirp, err3 := cfg.DBQueries.GetChirpById(req.Context(), parsed_chirp_id)
	if err3 != nil {
		errorResBody.Error = "Error while fetching chirp with id '" + chirp_id + "' from database: " + err3.Error()
		jsonResBody, err4 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err4, 404)
		return
	}

	successResBody := Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	jsonResBody, err5 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err5, 200)
}
