package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type commandeResp struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Product   string  `json:"product"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

func TestCreateCommandeIntegration(t *testing.T) {
	_ = godotenv.Load("../.env")

	payload := map[string]interface{}{
		"user_id": "123e4567-e89b-12d3-a456-426614174000", // UUID valide
		"product": "Test integration",
		"amount":  99.99,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://service-commandes:8082/commandes", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	defer resp.Body.Close()

	var commande commandeResp
	err = json.NewDecoder(resp.Body).Decode(&commande)
	assert.NoError(t, err)

	assert.Equal(t, "Test integration", commande.Product)
	assert.Equal(t, 99.99, commande.Amount)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", commande.UserID)
	assert.Equal(t, "en_attente", commande.Status)
	assert.NotEmpty(t, commande.ID)

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONN"))
	assert.NoError(t, err)
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM commandes WHERE id = $1", commande.ID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
