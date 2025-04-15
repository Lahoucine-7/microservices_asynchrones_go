package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/business"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/models"
)

// CreateCommandeInput représente les données pour créer une commande
type CreateCommandeInput struct {
	UserID  string  `json:"user_id" binding:"required,uuid"`
	Product string  `json:"product" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
}

// UpdateCommandeInput représente les données pour mettre à jour une commande
type UpdateCommandeInput struct {
	Product string  `json:"product" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
	Status  string  `json:"status" binding:"required"`
}

// Handler structure injectée avec un service
type Handler struct {
	CommandeService business.CommandeService
}

// NewHandler crée un handler avec dépendance injectée.
func NewHandler(service business.CommandeService) *Handler {
	return &Handler{CommandeService: service}
}

// HealthHandler traite GET /health
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateCommandeHandler traite POST /commandes
func (h *Handler) CreateCommandeHandler(c *gin.Context) {
	var input CreateCommandeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commande := models.Commande{
		ID:        uuid.New().String(),
		UserID:    input.UserID,
		Product:   input.Product,
		Amount:    input.Amount,
		Status:    "en_attente",
		CreatedAt: time.Now().UTC(),
	}

	if err := h.CommandeService.CreateCommande(c, commande); err != nil {
		log.Println("[Handler] Erreur création commande :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create commande"})
		return
	}

	c.JSON(http.StatusCreated, commande)
}

// GetAllCommandesHandler traite GET /commandes
func (h *Handler) GetAllCommandesHandler(c *gin.Context) {
	commandes, err := h.CommandeService.GetAllCommandes(c)
	if err != nil {
		log.Println("[Handler] Erreur récupération commandes :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch commandes"})
		return
	}

	c.JSON(http.StatusOK, commandes)
}

// GetCommandeByIDHandler traite GET /commandes/:id
func (h *Handler) GetCommandeByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	commande, err := h.CommandeService.GetCommandeByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "commande not found"})
		return
	}

	c.JSON(http.StatusOK, commande)
}

// UpdateCommandeHandler traite PUT /commandes/:id
func (h *Handler) UpdateCommandeHandler(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var input UpdateCommandeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := models.Commande{
		ID:        id,
		Product:   input.Product,
		Amount:    input.Amount,
		Status:    input.Status,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.CommandeService.UpdateCommande(c, id, updated); err != nil {
		log.Println("[Handler] Erreur mise à jour commande :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update commande"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteCommandeHandler traite DELETE /commandes/:id
func (h *Handler) DeleteCommandeHandler(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := h.CommandeService.DeleteCommande(c, id); err != nil {
		log.Println("[Handler] Erreur suppression commande :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete commande"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "commande deleted"})
}
