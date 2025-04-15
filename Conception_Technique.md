# Document de Conception Technique

## 1. Contexte et Objectifs du Projet

**Titre du Projet :**  
Construction d’un écosystème microservices avec communication asynchrone

**Contexte :**  
Le projet vise à transformer une application monolithique en un ensemble de microservices autonomes, capables de communiquer de manière asynchrone via RabbitMQ. Cela permet d’améliorer la scalabilité, la résilience et la facilité d’évolution, répondant ainsi aux besoins d’un client e-commerce soumis à des pics de charge.

**Objectifs Principaux :**
- Découper l’application en microservices (Utilisateurs, Commandes, Notifications).
- Mettre en place une communication asynchrone fiable (publication/consommation d’événements).
- Assurer l’autonomie technique de chaque service avec leur propre base de données (PostgreSQL).
- Respecter les bonnes pratiques de développement en Go (avec Gin pour les API REST) et en gestion de projet.

---

## 2. Architecture Globale et Découpage Fonctionnel

### 2.1 Diagramme d’Architecture
- **Vue d'ensemble :**  
  Le diagramme représente les microservices (Utilisateurs, Commandes, Notifications) connectés via un broker RabbitMQ.
- **Points Clés :**  
  - Un API Gateway pour centraliser et router les requêtes externes.
  - Chaque microservice possède sa propre base de données PostgreSQL.
  - Les flux asynchrones de publication et de consommation d’événements sont clairement identifiés.

### 2.2 Découpage des Microservices
- **Service Utilisateurs :**  
  Gère la création de profils, inscriptions, mises à jour.
- **Service Commandes :**  
  Responsable de la gestion des commandes (création, mise à jour, annulation).
- **Service Notifications :**  
  Consomme les événements pour déclencher l’envoi de notifications aux utilisateurs.

---

## 3. Communication Asynchrone

### 3.1 Choix Techniques
- **Broker de Messages :**  
  RabbitMQ, choisi pour sa robustesse et sa capacité à gérer les communications asynchrones en mode publication/consommation.
- **Types d’Échange :**  
  Le choix peut se porter sur un exchange direct, fanout ou topic. Dans notre cas, un exchange de type topic permettrait un routage flexible en fonction des clés des événements.

### 3.2 Événements et Contrats de Communication
- **Liste des Événements :**
  - *Service Utilisateurs* : UserCreated, UserUpdated, UserDeleted.
  - *Service Commandes* : OrderCreated, OrderUpdated, OrderCanceled.
  - *Service Notifications* : NotificationTriggered.
- **Structure Commune :**  
  Chaque message comporte les champs suivants :
  - **eventType** : indique le type d’événement.
  - **version** : permet de versionner le contrat pour une éventuelle évolution.
  - **timestamp** : date et heure d’émission.
  - **payload** : contient les données spécifiques à l’événement.
- **Schémas JSON :**  
  Des fichiers JSON Schema (par exemple, user_created.schema.json, order_created.schema.json, etc.) ont été rédigés pour formaliser et valider la structure des messages.

---

## 4. Structure du Projet et Organisation des Dossiers

Le projet est organisé dans le dossier racine **microservices_asynchones_go** et se présente ainsi :

    microservices_asynchones_go/
    ├── docker-compose.yaml         — Orchestration des conteneurs (Microservices, RabbitMQ, PostgreSQL)
    ├── common/                     — Code et ressources communes
    │   ├── events/                 — Schémas JSON et définitions des événements
    │   └── utils/                  — Fonctions utilitaires partagées
    ├── service-utilisateurs/       — Microservice dédié aux utilisateurs
    │   ├── cmd/                    — Point d'entrée (exemple : main.go)
    │   ├── internal/               
    │   │   ├── api/                — Gestion des endpoints REST avec Gin
    │   │   ├── business/           — Logique métier
    │   │   └── repository/         — Accès aux données (PostgreSQL)
    │   ├── configs/                
    │   ├── tests/                  
    │   ├── Dockerfile              
    │   └── README.md               
    ├── service-commandes/          — Microservice dédié aux commandes (structure similaire)
    ├── service-notifications/      — Microservice dédié aux notifications (structure similaire)
    └── README.md                   — Documentation globale du projet

**Explications :**
- Chaque microservice est autonome et présente une structure identique pour faciliter la maintenance.
- Le dossier "common" centralise les contrats de communication et les utilitaires partagés.
- Le fichier docker-compose.yaml permet de lancer l’ensemble de l’infrastructure avec une seule commande.

---

## 5. Contrats de Communication et Modèles de Données

### 5.1 Schémas JSON
Chaque événement dispose d’un fichier JSON Schema décrivant sa structure. Par exemple, pour "UserCreated" :
- Champs obligatoires : eventType (doit être "UserCreated"), version, timestamp, et payload (contenant userID, username, email, createdAt).

### 5.2 Fichier Go des Structs
Les définitions en Go se trouvent dans **common/events/events.go** et comprennent :
- **BaseEvent :** Structure commune (eventType, version, timestamp).
- **Structures Spécifiques :**  
  Par exemple, la structure **UserCreatedEvent** intègre **BaseEvent** et possède un champ **Payload** de type **UserCreatedPayload** qui regroupe les données spécifiques (userID, username, email, createdAt).

---

## 6. Conclusion et Prochaines Étapes

Cette étape a permis de :
- Comprendre le cahier des charges et définir précisément les objectifs du projet.
- Concevoir une architecture globale claire avec un découpage fonctionnel et une communication asynchrone via RabbitMQ.
- Organiser la structure du projet en adoptant des pratiques de développement modulaires et évolutives.
- Élaborer des contrats de communication (schémas JSON et structs en Go) pour standardiser les échanges entre microservices.

Ce document de conception technique servira de référence tout au long du développement du projet.

---

*Fin du Document de Conception Technique*
