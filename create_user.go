package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/OmarJarbou/Chirpy/internal/auth"
	"github.com/OmarJarbou/Chirpy/internal/database"
)

type createUserRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type createUserSuccessResponseBody struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(response_writer http.ResponseWriter, req *http.Request) {
	errorResBody := errorResponseBody{}
	var jsonResBody []byte

	hashed, err := auth.HashPassword(req.Context().Value("password").(string))
	if err != nil {
		errorResBody.Error = err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err2, 400)
		return
	}

	db_user := database.CreateUserParams{
		Email:          req.Context().Value("email").(string),
		HashedPassword: hashed,
	}
	user, err3 := cfg.DBQueries.CreateUser(req.Context(), db_user)
	if err3 != nil {
		errorResBody.Error = "Error while creating user: " + err3.Error()
		jsonResBody, err4 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err4, 500)
		return
	}

	successResBody := createUserSuccessResponseBody{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	jsonResBody, err5 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err5, 201)
}
