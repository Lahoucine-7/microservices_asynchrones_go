package business

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/models"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/repository"
)

// CommandeCreatedEvent représente un message CommandeCreated publié dans RabbitMQ.
type CommandeCreatedEvent struct {
	EventType string          `json:"eventType"`
	Version   string          `json:"version"`
	Timestamp string          `json:"timestamp"`
	Payload   models.Commande `json:"payload"`
}

// Service est l’implémentation concrète de l’interface CommandeService.
type Service struct{}

var (
	insertCommande   = repository.InsertCommande
	getAllCommandes  = repository.GetAllCommandes
	getCommandeByID  = repository.GetCommandeByID
	updateCommande   = repository.UpdateCommande
	deleteCommande   = repository.DeleteCommande
)

// CreateCommande insère la commande en base et publie l'événement.
func (s Service) CreateCommande(ctx context.Context, commande models.Commande) error {
	if err := insertCommande(commande); err != nil {
		return err
	}
	return publishCommandeCreated(commande)
}

// GetAllCommandes retourne toutes les commandes.
func (s Service) GetAllCommandes(ctx context.Context) ([]models.Commande, error) {
	return getAllCommandes()
}

// GetCommandeByID retourne une commande par ID.
func (s Service) GetCommandeByID(ctx context.Context, id string) (*models.Commande, error) {
	c, err := getCommandeByID(id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// UpdateCommande met à jour une commande existante.
func (s Service) UpdateCommande(ctx context.Context, id string, updated models.Commande) error {
	updated.ID = id // assurer que l’ID reste le même
	return updateCommande(updated)
}

// DeleteCommande supprime une commande.
func (s Service) DeleteCommande(ctx context.Context, id string) error {
	return deleteCommande(id)
}

// publication RabbitMQ
var publishCommandeCreated = func(commande models.Commande) error {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Println("[RabbitMQ] Connexion échouée :", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("[RabbitMQ] Échec ouverture channel :", err)
		return err
	}
	defer ch.Close()

	event := CommandeCreatedEvent{
		EventType: "CommandeCreated",
		Version:   "1.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Payload:   commande,
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Println("[RabbitMQ] JSON invalide :", err)
		return err
	}

	err = ch.Publish(
		"events",           // exchange
		"commande.created", // routing key
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("[RabbitMQ] Erreur publication :", err)
	}

	return err
}
