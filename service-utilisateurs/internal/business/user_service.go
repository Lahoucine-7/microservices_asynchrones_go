package business

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/models"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/repository"
)

// UserCreatedEvent représente un message UserCreated publié dans RabbitMQ.
type UserCreatedEvent struct {
	EventType string        `json:"eventType"`
	Version   string        `json:"version"`
	Timestamp string        `json:"timestamp"`
	Payload   models.User   `json:"payload"`
}

// Service est l’implémentation concrète de l’interface UserService.
type Service struct{}

var insertUser = repository.InsertUser

// CreateUser insère l'utilisateur en base et publie l'événement.
func (s Service) CreateUser(ctx context.Context, user models.User) error {
	if err := insertUser(user); err != nil {
		return err
	}
	return publishUserCreated(user)
}


var publishUserCreated = func(user models.User) error {
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

	event := UserCreatedEvent{
		EventType: "UserCreated",
		Version:   "1.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Payload:   user,
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Println("[RabbitMQ] JSON invalide :", err)
		return err
	}

	err = ch.Publish(
		"events",         // exchange
		"user.created",   // routing key
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
