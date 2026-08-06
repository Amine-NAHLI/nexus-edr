# Étape 2 : L'Agent (Golang)

## 🎯 Objectif
Développer le logiciel "Client" (l'Agent) qui sera installé sur le PC de l'utilisateur pour le surveiller et envoyer ses données au Backend.

## 🛠️ Pourquoi Golang (Go) ?
Go a été inventé par Google pour être ultra-rapide et consommer très peu de mémoire. Surtout, Go compile le code en **un seul fichier exécutable** (.exe sur Windows, ou binaire sur Linux). Il n'y a donc pas besoin d'installer Python ou Java sur les PC cibles, juste de lancer l'exécutable !

## 📝 Ce que nous allons faire en pratique :

1. **Initialiser le module Go :**
   C'est l'équivalent de créer un projet.
   *Commande :* `go mod init nexus-agent`

2. **Installer la bibliothèque gopsutil :**
   C'est la librairie magique en Go qui permet d'extraire la consommation CPU, la RAM, et les processus de n'importe quel système d'exploitation.
   *Commande :* `go get github.com/shirou/gopsutil/v3`

3. **Créer l'architecture :**
   ```text
   agent/
   ├── main.go          (Boucle principale du programme)
   ├── collector.go     (S'occupe de lire le CPU/RAM/Logs)
   └── sender.go        (S'occupe d'envoyer la requête HTTP/gRPC au Backend)
   ```

4. **Coder le flux principal :**
   - Le programme démarre.
   - Il lit un "Secret Token" dans un fichier ou via ligne de commande.
   - Toutes les 10 secondes (Boucle infinie), il appelle `collector.go` pour avoir le % de CPU.
   - Il formate ça en JSON.
   - Il appelle `sender.go` pour faire une requête `POST http://localhost:8000/api/telemetry`.

**👉 Prochaine étape logique une fois terminé :** Créer l'interface web pour visualiser ces données.
