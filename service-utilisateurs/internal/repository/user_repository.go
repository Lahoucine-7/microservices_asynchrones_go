package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-utilisateurs/internal/models"
)

var db *sql.DB

// InitDB initialise la connexion à la base de données.
func InitDB() error {
	var err error
	db, err = sql.Open("postgres", os.Getenv("POSTGRES_CONN"))
	if err != nil {
		log.Println("[DB] Connexion échouée :", err)
		return err
	}

	if err := db.Ping(); err != nil {
		log.Println("[DB] Ping échoué :", err)
		return err
	}

	log.Println("[DB] Connexion réussie")
	return nil
}

// InsertUser insère un utilisateur dans la base.
func InsertUser(u models.User) error {
	_, err := db.Exec(`
		INSERT INTO users (id, username, email, created_at)
		VALUES ($1, $2, $3, $4)
	`, u.ID, u.Username, u.Email, u.CreatedAt)
	return err
}
