package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/models"
)

// fakeUserService mocke l’interface business.UserService pour les tests
type fakeUserService struct{}

func (f fakeUserService) CreateUser(_ context.Context, _ models.User) error {
	return nil
}

type failingUserService struct{}

func (f failingUserService) CreateUser(_ context.Context, _ models.User) error {
	return assert.AnError
}

type failingDBService struct{}

func (f failingDBService) CreateUser(_ context.Context, _ models.User) error {
	return assert.AnError // Simule une erreur métier (ex : DB down)
}

type failingMQService struct{}

func (f failingMQService) CreateUser(_ context.Context, _ models.User) error {
	return assert.AnError // Simule une erreur lors de la publication
}


func TestCreateUserHandler_Mock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(fakeUserService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	payload := map[string]string{
		"username": "lahoucine",
		"email":    "lahoucine@example.com",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestCreateUserHandler_ServiceFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(failingUserService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	payload := map[string]string{
		"username": "lahoucine",
		"email":    "fail@example.com",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "could not create user", response["error"])
}


func TestCreateUserHandler_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(failingDBService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	payload := map[string]string{
		"username": "failuser",
		"email":    "failuser@example.com",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "could not create user", response["error"])
}

func TestCreateUserHandler_RabbitMQError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(failingMQService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	payload := map[string]string{
		"username": "mqfail",
		"email":    "mqfail@example.com",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "could not create user", response["error"])
}

func TestCreateUserHandler_MissingEmail(t *testing.T) {
	handler := NewHandler(fakeUserService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	payload := map[string]string{
		"username": "testuser",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateUserHandler_InvalidEmail(t *testing.T) {
	handler := NewHandler(fakeUserService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	payload := map[string]string{
		"username": "testuser",
		"email":    "not-an-email",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateUserHandler_MissingUsername(t *testing.T) {
	handler := NewHandler(fakeUserService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	payload := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	handler := NewHandler(fakeUserService{})
	router := gin.Default()
	router.POST("/users", handler.CreateUserHandler)

	badJSON := `{"username": "ok", "email": }`

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(badJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
