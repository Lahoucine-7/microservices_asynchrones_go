package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/repository"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/server"
)

func main() {
	// Charger les variables d'environnement
	err := godotenv.Load()
	if err != nil {
		log.Println("Aucun fichier .env trouvé, on continue avec les variables système...")
	}

	if err := repository.InitDB(); err != nil {
		log.Fatalf("Échec d'initialisation de la base de données : %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // fallback
	}

	fmt.Println("Démarrage du service utilisateurs sur le port", port)

	router := server.SetupRouter()
	router.Run(":" + port)
}