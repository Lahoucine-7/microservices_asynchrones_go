// File: common/events/events.go
package events

import "time"

// BaseEvent définit les champs communs à tous les événements.
type BaseEvent struct {
	EventType string    `json:"eventType"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

// --- Événements du Service Utilisateurs ---

// UserCreatedPayload représente le contenu d'un UserCreated.
type UserCreatedPayload struct {
	UserID    string    `json:"userID"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// UserCreatedEvent représente l'événement de création d'un utilisateur.
type UserCreatedEvent struct {
	BaseEvent
	Payload UserCreatedPayload `json:"payload"`
}

// UserUpdatedPayload représente le contenu d'un UserUpdated.
type UserUpdatedPayload struct {
	UserID    string    `json:"userID"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserUpdatedEvent représente l'événement de mise à jour d'un utilisateur.
type UserUpdatedEvent struct {
	BaseEvent
	Payload UserUpdatedPayload `json:"payload"`
}

// UserDeletedPayload représente le contenu d'un UserDeleted.
type UserDeletedPayload struct {
	UserID    string    `json:"userID"`
	DeletedAt time.Time `json:"deletedAt"`
}

// UserDeletedEvent représente l'événement de suppression d'un utilisateur.
type UserDeletedEvent struct {
	BaseEvent
	Payload UserDeletedPayload `json:"payload"`
}

// --- Événements du Service Commandes ---

// OrderItem représente un article dans une commande.
type OrderItem struct {
	ProductID string `json:"productID"`
	Quantity  int    `json:"quantity"`
}

// OrderCreatedPayload représente le contenu d'un OrderCreated.
type OrderCreatedPayload struct {
	OrderID     string      `json:"orderID"`
	UserID      string      `json:"userID"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
	OrderDate   time.Time   `json:"orderDate"`
}

// OrderCreatedEvent représente l'événement de création d'une commande.
type OrderCreatedEvent struct {
	BaseEvent
	Payload OrderCreatedPayload `json:"payload"`
}

// OrderUpdatedPayload représente le contenu d'un OrderUpdated.
type OrderUpdatedPayload struct {
	OrderID     string      `json:"orderID"`
	UserID      string      `json:"userID"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
	OrderDate   time.Time   `json:"orderDate"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// OrderUpdatedEvent représente l'événement de mise à jour d'une commande.
type OrderUpdatedEvent struct {
	BaseEvent
	Payload OrderUpdatedPayload `json:"payload"`
}

// OrderCanceledPayload représente le contenu d'un OrderCanceled.
type OrderCanceledPayload struct {
	OrderID    string    `json:"orderID"`
	UserID     string    `json:"userID"`
	CanceledAt time.Time `json:"canceledAt"`
	Reason     string    `json:"reason"`
}

// OrderCanceledEvent représente l'événement d'annulation d'une commande.
type OrderCanceledEvent struct {
	BaseEvent
	Payload OrderCanceledPayload `json:"payload"`
}

// --- Événement du Service Notifications ---

// NotificationTriggeredPayload représente le contenu d'un NotificationTriggered.
type NotificationTriggeredPayload struct {
	NotificationID string    `json:"notificationID"`
	UserID         string    `json:"userID"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"createdAt"`
}

// NotificationTriggeredEvent représente l'événement de déclenchement d'une notification.
type NotificationTriggeredEvent struct {
	BaseEvent
	Payload NotificationTriggeredPayload `json:"payload"`
}
