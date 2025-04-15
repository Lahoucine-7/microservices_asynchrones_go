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

var mockID = "11111111-1111-1111-1111-111111111111"

// fakeCommandService simule un service fonctionnel
type fakeCommandService struct{}

func (f fakeCommandService) CreateCommande(_ context.Context, _ models.Commande) error {
	return nil
}

func (f fakeCommandService) GetAllCommandes(_ context.Context) ([]models.Commande, error) {
	return []models.Commande{
		{
			ID:        mockID,
			UserID:    "123e4567-e89b-12d3-a456-426614174000",
			Product:   "Produit test",
			Amount:    49.99,
			Status:    "en_attente",
			CreatedAt: time.Now().UTC(),
		},
	}, nil
}

func (f fakeCommandService) GetCommandeByID(_ context.Context, id string) (*models.Commande, error) {
	if id != mockID {
		return nil, errors.New("not found")
	}
	return &models.Commande{
		ID:        mockID,
		UserID:    "123e4567-e89b-12d3-a456-426614174000",
		Product:   "Produit test",
		Amount:    49.99,
		Status:    "en_attente",
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (f fakeCommandService) UpdateCommande(_ context.Context, id string, _ models.Commande) error {
	if id != mockID {
		return errors.New("not found")
	}
	return nil
}

func (f fakeCommandService) DeleteCommande(_ context.Context, id string) error {
	if id != mockID {
		return errors.New("not found")
	}
	return nil
}

// === TESTS ===

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

func TestGetAllCommandesHandler(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	req, _ := http.NewRequest(http.MethodGet, "/commandes", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetCommandeByIDHandler_Success(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	req, _ := http.NewRequest(http.MethodGet, "/commandes/"+mockID, nil)
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
		"status": "en_attente",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPut, "/commandes/"+mockID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteCommandeHandler_Success(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	req, _ := http.NewRequest(http.MethodDelete, "/commandes/"+mockID, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateCommandeHandler_MissingFields(t *testing.T) {
	router := setupRouterWith(fakeCommandService{})

	payload := map[string]interface{}{
		"user_id": uuid.New().String(),
		"amount":  10.00, // manque product
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/commandes", bytes.NewBuffer(body))
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

// setupRouterWith est une fonction utilitaire locale aux tests
func setupRouterWith(service business.CommandeService) *gin.Engine {
	router := gin.Default()
	handler := NewHandler(service)

	router.POST("/commandes", handler.CreateCommandeHandler)
	router.GET("/commandes", handler.GetAllCommandesHandler)
	router.GET("/commandes/:id", handler.GetCommandeByIDHandler)
	router.PUT("/commandes/:id", handler.UpdateCommandeHandler)
	router.DELETE("/commandes/:id", handler.DeleteCommandeHandler)

	return router
}
