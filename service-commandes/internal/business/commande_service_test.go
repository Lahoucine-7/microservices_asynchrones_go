package business

import (
	"context"
	"testing"
	"time"

	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommande_MockDependencies(t *testing.T) {
	// 🔁 Mock InsertCommande
	originalInsert := insertCommande
	defer func() { insertCommande = originalInsert }()

	insertCommande = func(cmd models.Commande) error {
		assert.Equal(t, "Souris ergonomique", cmd.Product)
		assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", cmd.UserID)
		assert.Equal(t, 39.99, cmd.Amount)
		return nil
	}

	// 🔁 Mock publishCommandeCreated
	originalPublisher := publishCommandeCreated
	defer func() { publishCommandeCreated = originalPublisher }()

	publishCommandeCreated = func(cmd models.Commande) error {
		assert.Equal(t, "en_attente", cmd.Status)
		assert.NotEmpty(t, cmd.ID)
		return nil
	}

	cmd := models.Commande{
		ID:        "test-id-commande",
		UserID:    "123e4567-e89b-12d3-a456-426614174000",
		Product:   "Souris ergonomique",
		Amount:    39.99,
		Status:    "en_attente",
		CreatedAt: time.Now().UTC(),
	}

	service := Service{}
	err := service.CreateCommande(context.Background(), cmd)
	assert.NoError(t, err)
}
