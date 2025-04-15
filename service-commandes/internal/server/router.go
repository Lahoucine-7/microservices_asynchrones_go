package server

import (
	"github.com/gin-gonic/gin"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/api"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/business"
)

// SetupRouter configure les routes HTTP pour le service commandes.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", api.HealthHandler)

	handler := api.NewHandler(business.Service{}) // instance réelle ici

	// Routes REST
	router.POST("/commandes", handler.CreateCommandeHandler)
	router.GET("/commandes", handler.GetAllCommandesHandler)
	router.GET("/commandes/:id", handler.GetCommandeByIDHandler)
	router.PUT("/commandes/:id", handler.UpdateCommandeHandler)
	router.DELETE("/commandes/:id", handler.DeleteCommandeHandler)

	return router
}
