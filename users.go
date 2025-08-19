package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/OmarJarbou/Chirpy/internal/auth"
	"github.com/OmarJarbou/Chirpy/internal/database"
	"github.com/google/uuid"
)

type userRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type loginRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

const DEFAULT_TOKEN_EXP_TIME = 3600
const DEFAULT_REFRESH_TOKEN_EXP_TIME = 60 * 3600

type updateUserSuccessResponseBody struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type createUserORLoginSuccessResponseBody struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type refreshSuccessResponseBody struct {
	Token string `json:"token"`
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

	duration := time.Duration(DEFAULT_TOKEN_EXP_TIME) * time.Second
	token, err7 := auth.MakeJWT(user.ID, cfg.ChirpySecretKey, duration)
	if err7 != nil {
		errorResBody.Error = err7.Error()
		jsonResBody, err8 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err8, 400)
		return
	}

	refreshTokenString, err9 := auth.MakeRefreshToken()
	if err9 != nil {
		errorResBody.Error = err9.Error()
		jsonResBody, err10 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err10, 400)
		return
	}

	refresh_token_duration := time.Duration(DEFAULT_REFRESH_TOKEN_EXP_TIME) * time.Second
	refresh_token_expiration_time := time.Now().Add(refresh_token_duration)
	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: refresh_token_expiration_time,
	}

	refreshToken, err11 := cfg.DBQueries.CreateRefreshToken(req.Context(), createRefreshTokenParams)
	if err11 != nil {
		errorResBody.Error = "Error while creating new refresh token: " + err11.Error()
		jsonResBody, err12 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err12, 400)
		return
	}

	successResBody := createUserORLoginSuccessResponseBody{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	}
	jsonResBody, err13 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err13, 200)
}

func (cfg *apiConfig) handleRefreshToken(response_writer http.ResponseWriter, req *http.Request) {
	errorResBody := errorResponseBody{}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResBody.Error = err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err2, 400)
		return
	}

	user, err3 := cfg.DBQueries.GetUserFromRefreshToken(req.Context(), tokenString)
	if err3 != nil {
		errorResBody.Error = "Error while fetching user by this token: " + err3.Error()
		jsonResBody, err4 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err4, 401)
		return
	}

	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid {
		errorResBody.Error = "Your refresh token is invalid, (expired or revoked)!"
		jsonResBody, err5 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err5, 401)
		return
	}

	duration := time.Duration(DEFAULT_TOKEN_EXP_TIME) * time.Second
	new_token, err6 := auth.MakeJWT(user.ID, cfg.ChirpySecretKey, duration)
	if err6 != nil {
		errorResBody.Error = err6.Error()
		jsonResBody, err7 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err7, 400)
		return
	}

	successResBody := refreshSuccessResponseBody{
		Token: new_token,
	}

	jsonResBody, err8 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err8, 200)
}

func (cfg *apiConfig) handleRevokeToken(response_writer http.ResponseWriter, req *http.Request) {
	errorResBody := errorResponseBody{}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResBody.Error = err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err2, 400)
		return
	}

	err3 := cfg.DBQueries.SetRefreshTokenAsRevoked(req.Context(), tokenString)
	if err3 != nil {
		errorResBody.Error = "Error while revoking user's token: " + err3.Error()
		jsonResBody, err4 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err4, 400)
		return
	}

	response_writer.WriteHeader(204)
}

func (cfg *apiConfig) handleUpdateUser(response_writer http.ResponseWriter, req *http.Request) {
	errorResBody := errorResponseBody{}
	var jsonResBody []byte

	email := req.Context().Value("email").(string)
	password := req.Context().Value("password").(string)
	user_id := req.Context().Value("user_id").(uuid.UUID)

	hashed, err := auth.HashPassword(password)
	if err != nil {
		errorResBody.Error = err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err2, 400)
		return
	}

	db_user := database.UpdateUserInfoParams{
		Email:          email,
		HashedPassword: hashed,
		ID:             user_id,
	}
	user, err5 := cfg.DBQueries.UpdateUserInfo(req.Context(), db_user)
	if err5 != nil {
		errorResBody.Error = "Error while creating user: " + err5.Error()
		jsonResBody, err6 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err6, 500)
		return
	}

	successResBody := updateUserSuccessResponseBody{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	jsonResBody, err13 := json.Marshal(successResBody)
	writeJSONResponse(response_writer, jsonResBody, err13, 200)
}

func (cfg *apiConfig) handleDeleteChirp(response_writer http.ResponseWriter, req *http.Request) {
	errorResBody := errorResponseBody{}

	user_id := req.Context().Value("user_id").(uuid.UUID)

	chirp_id := req.PathValue("chirpID")

	chirp_uuid, err := uuid.Parse(chirp_id)
	if err != nil {
		errorResBody.Error = "Error while parsing chirp id string to uuid: " + err.Error()
		jsonResBody, err2 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err2, 400)
		return
	}

	chirp, err3 := cfg.DBQueries.GetChirpById(req.Context(), chirp_uuid)
	if err3 != nil {
		errorResBody.Error = "No chirp with this id: " + err3.Error()
		jsonResBody, err4 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err4, 404)
		return
	}

	if chirp.UserID != user_id {
		errorResBody.Error = "You don't have access to do anything on this chirp: "
		jsonResBody, err5 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err5, 403)
		return
	}

	err6 := cfg.DBQueries.DeleteChirpById(req.Context(), chirp.ID)
	if err6 != nil {
		errorResBody.Error = "Error while deleting chirp: " + err6.Error()
		jsonResBody, err7 := json.Marshal(errorResBody)
		writeJSONResponse(response_writer, jsonResBody, err7, 403)
		return
	}

	response_writer.WriteHeader(204)
}
