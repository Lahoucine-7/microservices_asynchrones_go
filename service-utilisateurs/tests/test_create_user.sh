#!/bin/bash

echo "🔁 Test création utilisateur (POST /users)"

curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com"
  }'

echo -e "\n✅ Requête envoyée"
