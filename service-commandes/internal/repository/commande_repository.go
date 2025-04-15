package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/Lahoucine-7/microservices_asynchrones_go/service-commandes/internal/models"
)

var db *sql.DB

// InitDB initialise la connexion à la base de données PostgreSQL.
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

// InsertCommande insère une commande dans la base.
func InsertCommande(c models.Commande) error {
	_, err := db.Exec(`
		INSERT INTO commandes (id, user_id, product, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, c.ID, c.UserID, c.Product, c.Amount, c.Status, c.CreatedAt)
	return err
}

// GetAllCommandes retourne toutes les commandes.
func GetAllCommandes() ([]models.Commande, error) {
	rows, err := db.Query(`SELECT id, user_id, product, amount, status, created_at FROM commandes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commandes []models.Commande
	for rows.Next() {
		var c models.Commande
		if err := rows.Scan(&c.ID, &c.UserID, &c.Product, &c.Amount, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		commandes = append(commandes, c)
	}
	return commandes, nil
}

// GetCommandeByID retourne une commande par ID.
func GetCommandeByID(id string) (models.Commande, error) {
	var c models.Commande
	err := db.QueryRow(`
		SELECT id, user_id, product, amount, status, created_at
		FROM commandes WHERE id = $1
	`, id).Scan(&c.ID, &c.UserID, &c.Product, &c.Amount, &c.Status, &c.CreatedAt)
	return c, err
}

// UpdateCommande met à jour une commande.
func UpdateCommande(c models.Commande) error {
	_, err := db.Exec(`
		UPDATE commandes
		SET user_id = $1, product = $2, amount = $3, status = $4
		WHERE id = $5
	`, c.UserID, c.Product, c.Amount, c.Status, c.ID)
	return err
}

// DeleteCommande supprime une commande.
func DeleteCommande(id string) error {
	_, err := db.Exec(`DELETE FROM commandes WHERE id = $1`, id)
	return err
}
