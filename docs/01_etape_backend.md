# Étape 1 : Le Backend (FastAPI)

## 🎯 Objectif
Créer le cœur du système. C'est l'API qui va recevoir les données des agents et les enregistrer dans la base de données.

## 🛠️ Pourquoi FastAPI ?
Python est incontournable pour l'Intelligence Artificielle. FastAPI est actuellement le framework Python le plus rapide et moderne pour créer des API. Il gère l'asynchrone (pour encaisser des milliers de logs par seconde) et génère automatiquement la documentation (Swagger).

## 📝 Ce que nous allons faire en pratique :

1. **Créer l'environnement virtuel Python (`venv`) :**
   Pour isoler nos bibliothèques et ne pas polluer l'ordinateur.
   *Commande :* `python -m venv venv`

2. **Installer les dépendances :**
   *   `fastapi` : Le framework.
   *   `uvicorn` : Le serveur web qui fait tourner FastAPI.
   *   `sqlalchemy` : L'ORM (l'outil pour parler à la base de données PostgreSQL sans écrire de SQL pur).
   *   `psycopg2` : Le connecteur PostgreSQL.

3. **Créer l'architecture des dossiers Backend :**
   ```text
   backend/
   ├── app/
   │   ├── main.py        (Point d'entrée de l'API)
   │   ├── models.py      (Schéma de la base de données)
   │   ├── schemas.py     (Le "Contrat JSON" Pydantic)
   │   └── routes.py      (Les chemins comme /agents, /logs)
   └── requirements.txt
   ```

4. **Coder la route "Hello World" :**
   Nous écrirons 10 lignes de code dans `main.py` pour lancer le serveur et vérifier qu'il répond bien sur `http://localhost:8000`.

**👉 Prochaine étape logique une fois terminé :** Créer un Agent en Go pour envoyer des données à cette API.
