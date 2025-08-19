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

type loginRequestBody struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
}

const DEFAULT_TOKEN_EXP_TIME = 3600

type createUserORLoginSuccessResponseBody struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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

	successResBody := createUserORLoginSuccessResponseBody{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	jsonResBody, err5 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err5, 201)
}

func (cfg *apiConfig) handleLogin(response_writer http.ResponseWriter, req *http.Request) {
	reqBody := loginRequestBody{}
	errorResBody := errorResponseBody{}
	var jsonResBody []byte
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		errorResBody.Error = "Error while decoding request's json " + err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)

		writeJSONResponse(response_writer, jsonResBody, err2, 500)
		return
	}

	user, err3 := cfg.DBQueries.GetUserByEmail(req.Context(), reqBody.Email)
	if err3 != nil {
		errorResBody.Error = "Error while fetching user by this email: " + err3.Error()
		jsonResBody, err4 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err4, 401)
		return
	}

	err5 := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err5 != nil {
		errorResBody.Error = err5.Error()
		jsonResBody, err6 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err6, 401)
		return
	}

	var expirationSeconds int
	if reqBody.ExpiresInSeconds == nil || *reqBody.ExpiresInSeconds < 0 || *reqBody.ExpiresInSeconds > 3600 {
		expirationSeconds = DEFAULT_TOKEN_EXP_TIME
	} else {
		expirationSeconds = *reqBody.ExpiresInSeconds
	}

	duration := time.Duration(expirationSeconds) * time.Second
	token, err7 := auth.MakeJWT(user.ID, cfg.ChirpySecretKey, duration)
	if err7 != nil {
		errorResBody.Error = err7.Error()
		jsonResBody, err8 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err8, 400)
		return
	}

	successResBody := createUserORLoginSuccessResponseBody{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	jsonResBody, err9 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err9, 200)
}
