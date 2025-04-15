package business

import (
	"context"
	"testing"
	"time"

	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/models"
	"github.com/stretchr/testify/assert"
)
func TestCreateUser_MockRabbitMQ(t *testing.T) {
	// 🔁 Mock InsertUser
	originalInsertUser := insertUser
	defer func() { insertUser = originalInsertUser }()
	insertUser = func(user models.User) error {
		assert.Equal(t, "testuser@example.com", user.Email)
		return nil
	}

	// 🔁 Mock publishUserCreated
	originalPublisher := publishUserCreated
	defer func() { publishUserCreated = originalPublisher }()
	publishUserCreated = func(user models.User) error {
		assert.Equal(t, "lahoucine", user.Username)
		return nil
	}

	user := models.User{
		ID:        "test-id",
		Username:  "lahoucine",
		Email:     "testuser@example.com",
		CreatedAt: time.Now().UTC(),
	}

	service := Service{}
	err := service.CreateUser(context.Background(), user)
	assert.NoError(t, err)
}

