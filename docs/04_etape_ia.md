# Étape 4 : L'Intelligence Artificielle (Le Cerveau)

## 🎯 Objectif
Donner à notre EDR la capacité de détecter des attaques inconnues (Machine Learning) et d'apprendre des décisions de l'administrateur pour devenir autonome (Reinforcement Learning).

## 🛠️ Pourquoi Scikit-Learn et RLlib (PyTorch) ?
*   **Scikit-Learn** est parfait pour créer rapidement des modèles de "Détection d'anomalies" très légers et mathématiquement prouvés (comme l'Isolation Forest).
*   **PyTorch / RLlib** est le standard pour l'apprentissage par renforcement. Cela permettra à notre agent de prendre des décisions (Bloquer, Ignorer) basées sur un système de récompenses.

## 📝 Ce que nous allons faire en pratique :

1. **Créer le module IA (Worker) :**
   Ce module tournera en parallèle du Backend FastAPI.
   ```text
   ai_engine/
   ├── train_ml.py        (Entraîne le modèle sur l'historique ClickHouse)
   ├── detector.py        (Analyse les logs en temps réel via Kafka/Redis)
   └── rl_agent.py        (Prend les décisions et apprend des feedbacks)
   ```

2. **Phase 1 : Le Machine Learning (Apprentissage non supervisé)**
   - Extraire les logs des 7 derniers jours.
   - Transformer les données (ex: l'heure de la journée devient un chiffre, le % CPU reste un chiffre). C'est le "Feature Engineering".
   - Entraîner une *Isolation Forest* pour qu'elle apprenne ce qu'est un comportement "Normal".

3. **Phase 2 : Le Reinforcement Learning (Apprentissage par renforcement)**
   - Créer un environnement où l'État (State) = L'alerte générée par le ML.
   - Les Actions = [Alerter Admin, Tuer Processus, Isoler le PC du réseau].
   - La Récompense (Reward) = Si l'admin clique sur "Vrai Positif" sur le Dashboard, l'IA gagne +1. Si l'admin clique sur "Faux Positif", l'IA perd -1.

**👉 Finalité :** Une plateforme capable d'agir seule la nuit quand l'administrateur dort, avec un taux de faux positifs proche de zéro grâce à l'entraînement !
