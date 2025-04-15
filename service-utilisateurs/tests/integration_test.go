package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type userResp struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func TestCreateUserIntegration(t *testing.T) {
	_ = godotenv.Load("../.env")
	payload := map[string]string{
		"username": "integration_test",
		"email":    "integration@example.com",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://service-utilisateurs:8081/users", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	defer resp.Body.Close()

	var user userResp
	json.NewDecoder(resp.Body).Decode(&user)

	assert.Equal(t, "integration_test", user.Username)
	assert.Equal(t, "integration@example.com", user.Email)
	assert.NotEmpty(t, user.ID)

	t.Log("POSTGRES_CONN =", os.Getenv("POSTGRES_CONN"))

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONN"))
	assert.NoError(t, err)
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", user.ID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
