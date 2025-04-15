package business

import (
	"context"

	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/models"
)

// CommandeService définit les opérations offertes par la couche métier.
type CommandeService interface {
	CreateCommande(ctx context.Context, commande models.Commande) error
	GetAllCommandes(ctx context.Context) ([]models.Commande, error)
	GetCommandeByID(ctx context.Context, id string) (*models.Commande, error)
	UpdateCommande(ctx context.Context, id string, update models.Commande) error
	DeleteCommande(ctx context.Context, id string) error
}
