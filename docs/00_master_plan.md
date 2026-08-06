# Master Plan : Plateforme EDR & Intelligence Artificielle (Nexus-EDR)

## 🎯 Objectif du Projet
Créer une plateforme de cybersécurité (EDR - Endpoint Detection and Response) complète, capable de surveiller un parc informatique en temps réel et d'utiliser l'Intelligence Artificielle pour bloquer les menaces.

## 🏗️ Architecture Globale
Le projet suit une architecture moderne en "Monorepo" (tous les sous-projets sont dans ce même dossier Git). Il est divisé en 4 grands piliers techniques :
1. **L'Agent (Golang) :** Le logiciel espion bienveillant installé sur les PC cibles.
2. **Le Backend (Python/FastAPI) :** Le cerveau central qui gère l'API et la base de données.
3. **Le Dashboard (Next.js) :** L'interface web pour l'administrateur.
4. **Le Moteur IA (Python) :** Les modèles de Machine Learning et de Reinforcement Learning.

## 🗺️ Comment lire cette documentation ?
Afin de construire ce projet étape par étape de manière pédagogique et structurée, nous avons divisé le travail en phases distinctes. 

Veuillez suivre les documents dans cet ordre :
- 👉 **[Étape 1]** Lisez `01_etape_backend.md` pour comprendre comment nous allons créer l'API centrale.
- 👉 **[Étape 2]** Lisez `02_etape_agent.md` pour le développement du collecteur de données en Go.
- 👉 **[Étape 3]** Lisez `03_etape_frontend.md` pour la conception de l'interface administrateur.
- 👉 **[Étape 4]** Lisez `04_etape_ia.md` pour comprendre comment nous allons implémenter le Machine Learning et le Reinforcement Learning.
