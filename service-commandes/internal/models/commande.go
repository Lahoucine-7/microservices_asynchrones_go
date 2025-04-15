package models

import "time"

type Commande struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Product   string    `json:"product"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CommandeInput struct {
	UserID  string  `json:"user_id" binding:"required"`
	Product string  `json:"product" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
}
