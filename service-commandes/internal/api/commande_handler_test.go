package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/models"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/business"
)

// fakeCommandService simule une réussite
type fakeCommandService struct{}

func (f fakeCommandService) CreateCommande(_ context.Context, _ models.Commande) error {
	return nil
}

func (f fakeCommandService) GetAllCommandes(_ context.Context) ([]models.Commande, error) {
	return []models.Commande{
		{
			ID:        "1",
			UserID:    "user1",
			Product:   "Product 1",
			Amount:    10.0,
			Status:    "en_attente",
			CreatedAt: time.Now().UTC(),
		},
	}, nil
}

func (f fakeCommandService) GetCommandeByID(_ context.Context, id string) (*models.Commande, error) {
	return &models.Commande{
		ID:        id,
		UserID:    "user-id",
		Product:   "Produit test",
		Amount:    25.5,
		Status:    "livrée",
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (f fakeCommandService) UpdateCommande(_ context.Context, id string, _ models.Commande) error {
	return nil
}

func (f fakeCommandService) DeleteCommande(_ context.Context, id string) error {
	return nil
}

// failingCommandService simule une erreur
type failingCommandService struct{}

func (f failingCommandService) CreateCommande(_ context.Context, _ models.Commande) error {
	return errors.New("insert error")
}

func (f failingCommandService) GetAllCommandes(_ context.Context) ([]models.Commande, error) {
	return nil, errors.New("get all error")
}

func (f failingCommandService) GetCommandeByID(_ context.Context, id string) (*models.Commande, error) {
	return nil, errors.New("not found")
}

func (f failingCommandService) UpdateCommande(_ context.Context, id string, _ models.Commande) error {
	return errors.New("update failed")
}

func (f failingCommandService) DeleteCommande(_ context.Context, id string) error {
	return errors.New("delete failed")
}

// ---------------------------- TESTS ----------------------------

func TestCreateCommandeHandler_Success(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	payload := map[string]interface{}{
		"user_id": uuid.New().String(),
		"product": "Test produit",
		"amount":  49.99,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/commandes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestCreateCommandeHandler_ServiceError(t *testing.T) {
	router := setupRouterWith(failingCommandService{})

	payload := map[string]interface{}{
		"user_id": uuid.New().String(),
		"product": "Erreur produit",
		"amount":  19.99,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/commandes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestGetAllCommandesHandler_Success(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	req, _ := http.NewRequest(http.MethodGet, "/commandes", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetCommandeByIDHandler_Success(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	validID := uuid.New().String()


	req, _ := http.NewRequest(http.MethodGet, validID, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestUpdateCommandeHandler_Success(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	payload := map[string]interface{}{
		"user_id": uuid.New().String(),
		"product": "Produit modifié",
		"amount":  59.99,
	}
	body, _ := json.Marshal(payload)

	validID := uuid.New().String()

	req, _ := http.NewRequest(http.MethodPut, validID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteCommandeHandler_Success(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	validID := uuid.New().String()

	req, _ := http.NewRequest(http.MethodDelete, validID, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateCommandeHandler_MissingFields(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	payload := map[string]interface{}{
		"user_id": uuid.New().String(),
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/commandes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateCommandeHandler_InvalidJSON(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	badJSON := `{"user_id": "ok", "product": }`

	req, _ := http.NewRequest("POST", "/commandes", bytes.NewBufferString(badJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// setupRouterWith installe un router Gin avec le handler simulé
func setupRouterWith(service business.CommandeService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewHandler(service)

	r.POST("/commandes", h.CreateCommandeHandler)
	r.GET("/commandes", h.GetAllCommandesHandler)
	r.GET("/commandes/:id", h.GetCommandeByIDHandler)
	r.PUT("/commandes/:id", h.UpdateCommandeHandler)
	r.DELETE("/commandes/:id", h.DeleteCommandeHandler)

	return r
}
