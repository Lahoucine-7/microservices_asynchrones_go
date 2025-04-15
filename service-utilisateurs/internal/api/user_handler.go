package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/business"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/models"
)

// CreateUserInput représente les données envoyées dans le POST /users
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Handler structure injectée avec un service
type Handler struct {
	UserService business.UserService
}

// NewHandler crée un handler avec dépendance injectée.
func NewHandler(service business.UserService) *Handler {
	return &Handler{UserService: service}
}

// CreateUserHandler traite POST /users
func (h *Handler) CreateUserHandler(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:        uuid.New().String(),
		Username:  input.Username,
		Email:     input.Email,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.UserService.CreateUser(c, user); err != nil {
		log.Println("[Handler] Erreur création utilisateur :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// HealthHandler traite GET /health
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
