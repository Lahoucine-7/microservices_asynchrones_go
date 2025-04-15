package business

import (
	"context"

	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/models"
)

// UserService définit les opérations offertes par la couche métier.
type UserService interface {
	CreateUser(ctx context.Context, user models.User) error
}
