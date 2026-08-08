# Feuille de Route d'Implémentation (Roadmap)

Ce document détaille pas à pas comment nous allons coder l'architecture complète de Nexus-EDR, de la première ligne de code jusqu'à la version finale XDR.

## Phase 1 : Les Fondations (Processus & Base de Données)
**Objectif :** Rendre l'Agent capable de voir ce qui tourne sur le PC et préparer le Backend à recevoir ces données complexes.
1. **Backend (Python) :**
   - Modifier `models.py` (SQLAlchemy) pour ajouter `process_name`, `process_path`, `process_hash`, `risk_score`.
   - Modifier `schemas.py` pour valider ces nouvelles données.
   - Supprimer l'ancienne base de données Docker pour qu'elle se recrée avec les nouvelles colonnes.
2. **Agent (Go) :**
   - Ajouter la librairie `crypto/sha256` et `process` de gopsutil.
   - Coder une fonction pour lister les processus en cours (ex: `notepad.exe`).
   - Coder une fonction qui lit le fichier `.exe` sur le disque dur et calcule son Hash.
   - Envoyer ces données au Backend.

## Phase 2 : Le Cerveau & Threat Intelligence
**Objectif :** Connecter le Backend aux bases de données mondiales pour juger les processus détectés à la Phase 1.
1. **Backend (Python) :**
   - Créer un module `threat_intel.py`.
   - Intégrer l'API **VirusTotal** (ou MalwareBazaar) : envoyer le Hash reçu par l'Agent et récupérer le score de dangerosité.
   - Créer la première version du **Risk Engine** : un algorithme mathématique simple qui calcule un score sur 100 en fonction de la réponse de VirusTotal.

## Phase 3 : Le Dashboard & Human-in-the-Loop
**Objectif :** Donner le pouvoir de décision à l'Administrateur pour agir sur les menaces de la Phase 2.
1. **Frontend (Next.js) :**
   - Mettre à jour l'interface pour afficher les détails du processus (Hash, Nom, Risk Score).
   - Créer des boutons d'action **[AUTORISER]** et **[BLOQUER]**.
2. **Backend & Agent :**
   - Créer une file d'attente (Endpoint API) où l'Agent Go vient vérifier s'il a reçu un ordre "BLOQUER".
   - Coder la fonction `process.Kill()` dans l'Agent Go pour qu'il détruise le malware si l'Admin a cliqué sur Bloquer.

## Phase 4 : Le Réseau & La Corrélation
**Objectif :** Analyser le trafic internet et le lier aux processus de la Phase 1.
1. **Agent (Go) :**
   - Extraire les connexions TCP/UDP et l'IP de destination.
   - Associer la connexion au PID du processus qui l'a créée.
2. **Backend (Python) :**
   - Intégrer l'API **AbuseIPDB** dans `threat_intel.py` pour vérifier la réputation de l'IP.
   - Mettre à jour le Risk Engine : *Score = Score Processus + Score IP*.

## Phase 5 : File Integrity Monitoring (Anti-Ransomware)
**Objectif :** Détecter la destruction de fichiers.
1. **Agent (Go) :**
   - Surveiller un dossier test.
   - Compter le nombre de modifications/suppressions par seconde effectuées par un processus.
2. **Backend & IA :**
   - Si les modifications dépassent la limite de l'Admin (Policy), déclencher une alerte CRITIQUE immédiate.

## Phase 6 : Logs Windows & Persistance
**Objectif :** Achever la transformation en XDR.
1. **Agent (Go) :**
   - Lire les événements de sécurité Windows (Connexions échouées).
   - Surveiller les clés de Registre (Démarrage automatique).
2. **Backend (Python) :**
   - Croiser toutes les alertes (Logs -> Processus -> Fichier -> Réseau) pour reconstruire la "Kill Chain" complète d'un attaquant.
