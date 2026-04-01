package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rendi-hendra/resful-api/internal/config"
	"github.com/rendi-hendra/resful-api/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	app      *fiber.App
	db       *gorm.DB
	log      *logrus.Logger
	validate *validator.Validate
	vConfig  *viper.Viper
)

func init() {
	vConfig = config.NewViper()
	log = config.NewLogger(vConfig)
	db = config.NewDatabase(vConfig, log)
	validate = config.NewValidator(vConfig)
	app = config.NewFiber(vConfig)

	config.Bootstrap(&config.BootstapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   vConfig,
	})
}

func ClearAll() {
	err := db.Exec("DELETE FROM users").Error
	if err != nil {
		log.Fatalf("Failed to clear database: %+v", err)
	}
}

func CreateUser(t *testing.T) {
	requestBody := model.RegisterUserRequest{
		ID:       "user1",
		Password: "password123",
		Name:     "User One",
		Email:    "user1@example.com",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func GetTokens(t *testing.T) *model.TokenResponse {
	requestBody := model.LoginUserRequest{
		ID:       "user1",
		Password: "password123",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.WebResponse[*model.TokenResponse]
	payload, _ := io.ReadAll(response.Body)
	json.Unmarshal(payload, &responseBody)
	return responseBody.Data
}

func TestRegisterSuccess(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		ID:       "user1",
		Password: "password123",
		Name:     "User One",
		Email:    "user1@example.com",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.WebResponse[*model.UserResponse]
	payload, _ := io.ReadAll(response.Body)
	json.Unmarshal(payload, &responseBody)
	assert.Equal(t, requestBody.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()
	CreateUser(t)

	requestBody := model.RegisterUserRequest{
		ID:       "user1",
		Password: "password123",
		Name:     "User One",
		Email:    "user1@example.com",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestRegisterValidationError(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		ID:       "",
		Password: "",
		Name:     "",
		Email:    "invalid-email",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestLoginSuccess(t *testing.T) {
	ClearAll()
	CreateUser(t)

	requestBody := model.LoginUserRequest{
		ID:       "user1",
		Password: "password123",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.WebResponse[*model.TokenResponse]
	payload, _ := io.ReadAll(response.Body)
	json.Unmarshal(payload, &responseBody)
	assert.NotEmpty(t, responseBody.Data.AccessToken)
	assert.NotEmpty(t, responseBody.Data.RefreshToken)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	CreateUser(t)

	requestBody := model.LoginUserRequest{
		ID:       "user1",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestLoginNotFound(t *testing.T) {
	ClearAll()

	requestBody := model.LoginUserRequest{
		ID:       "notfound",
		Password: "password123",
	}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestRefreshSuccess(t *testing.T) {
	ClearAll()
	CreateUser(t)
	tokens := GetTokens(t)

	// Refresh token
	refreshReq := httptest.NewRequest(http.MethodPost, "/refresh-token", nil)
	refreshReq.Header.Set("Authorization", "Bearer "+tokens.RefreshToken)

	refreshResp, _ := app.Test(refreshReq)
	assert.Equal(t, http.StatusOK, refreshResp.StatusCode)

	var newTokens model.WebResponse[*model.TokenResponse]
	newPayload, _ := io.ReadAll(refreshResp.Body)
	json.Unmarshal(newPayload, &newTokens)
	assert.NotEmpty(t, newTokens.Data.AccessToken)
}

func TestRefreshInvalidToken(t *testing.T) {
	ClearAll()

	refreshReq := httptest.NewRequest(http.MethodPost, "/refresh-token", nil)
	refreshReq.Header.Set("Authorization", "Bearer invalidtoken")

	refreshResp, _ := app.Test(refreshReq)
	assert.Equal(t, http.StatusUnauthorized, refreshResp.StatusCode)
}

func TestGetCurrentSuccess(t *testing.T) {
	ClearAll()
	CreateUser(t)
	tokens := GetTokens(t)

	// Get Current User
	currentReq := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	currentReq.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	currentResp, _ := app.Test(currentReq)
	assert.Equal(t, http.StatusOK, currentResp.StatusCode)

	var userResp model.WebResponse[*model.UserResponse]
	userPayload, _ := io.ReadAll(currentResp.Body)
	json.Unmarshal(userPayload, &userResp)
	assert.Equal(t, "user1", userResp.Data.ID)
}

func TestGetCurrentUnauthorized(t *testing.T) {
	ClearAll()

	currentReq := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	// No Authorization header

	currentResp, _ := app.Test(currentReq)
	assert.Equal(t, http.StatusUnauthorized, currentResp.StatusCode)
}

func TestUpdateSuccess(t *testing.T) {
	ClearAll()
	CreateUser(t)
	tokens := GetTokens(t)

	// Update User
	updateBody := model.UpdateUserRequest{
		Name:     "Updated Name",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}
	updatePayload, _ := json.Marshal(updateBody)
	updateReq := httptest.NewRequest(http.MethodPatch, "/api/users/_current", bytes.NewBuffer(updatePayload))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	updateResp, _ := app.Test(updateReq)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	var userResp model.WebResponse[*model.UserResponse]
	userPayload, _ := io.ReadAll(updateResp.Body)
	json.Unmarshal(userPayload, &userResp)
	assert.Equal(t, "Updated Name", userResp.Data.Name)
	assert.Equal(t, "updated@example.com", userResp.Data.Email)
}

func TestUpdateValidationError(t *testing.T) {
	ClearAll()
	CreateUser(t)
	tokens := GetTokens(t)

	// Update User with excessively long name
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "a"
	}
	updateBody := model.UpdateUserRequest{
		Name: longName,
	}
	updatePayload, _ := json.Marshal(updateBody)
	updateReq := httptest.NewRequest(http.MethodPatch, "/api/users/_current", bytes.NewBuffer(updatePayload))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	updateResp, _ := app.Test(updateReq)
	assert.Equal(t, http.StatusBadRequest, updateResp.StatusCode)
}

func TestUpdateUnauthorized(t *testing.T) {
	ClearAll()

	updateReq := httptest.NewRequest(http.MethodPatch, "/api/users/_current", nil)
	// No Authorization header

	updateResp, _ := app.Test(updateReq)
	assert.Equal(t, http.StatusUnauthorized, updateResp.StatusCode)
}
