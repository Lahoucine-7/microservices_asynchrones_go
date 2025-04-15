package server

import (
	"github.com/gin-gonic/gin"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/api"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/business"
)

// SetupRouter configure les routes HTTP.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", api.HealthHandler)

	handler := api.NewHandler(business.Service{}) // ← instance réelle ici
	router.POST("/users", handler.CreateUserHandler)

	return router
}
