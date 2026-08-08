# Mécanismes de Détection EDR (Agent Go)

## 1. Détection intelligente des processus par analyse comportementale IA et enrichissement Threat Intelligence

### 🎯 Le Concept (Double Moteur de Risque)
Cette fonctionnalité est le cœur d'un EDR de nouvelle génération (XDR). Elle ne se contente pas de bloquer un fichier connu comme malveillant, mais utilise un double moteur de décision :
1. **Threat Intelligence (API Externe) :** Interrogation de bases de données mondiales (ex: VirusTotal) via le hash du fichier pour savoir si le monde entier le connaît.
2. **Intelligence Artificielle (Comportemental) :** Analyse du comportement local du processus (chemin inhabituel, parent suspect, ligne de commande étrange, etc.) pour bloquer même les virus "Zero-Day" inconnus.
3. **Moteur de Risque (Risk Engine) :** Croise les deux informations pour générer un "Risk Score" final sur 100.

### 🛠️ Implémentation Technique

#### Partie 1 : L'Agent EDR (Golang)
L'Agent Go devient un véritable collecteur d'informations forensiques (Telemetry). En plus du CPU/RAM, il devra surveiller les processus et extraire pour chacun :
- **Identité** : Nom (`p.Name()`), PID et Parent PID (`p.Ppid()`).
- **Comportement** : Chemin d'exécution (`p.Exe()`) et Ligne de commande exacte (`p.Cmdline()`).
- **Signature (Crucial)** : L'Agent Go devra ouvrir le fichier binaire sur le disque dur de Windows et calculer lui-même son **Hash SHA-256** avec la librairie native `crypto/sha256`.
- **Réseau** : Liste des connexions liées à ce processus spécifique (`net.ConnectionsPid()`).

*L'Agent envoie ce "Dossier Complet" au Backend Python.*

#### Partie 2 : Le Cerveau Central (Python / FastAPI)
Le Backend reçoit le dossier et orchestre l'investigation (idéalement en parallèle pour ne pas ralentir le système) :
1. **Module Threat Intel :** Requête HTTP automatique vers l'API de VirusTotal avec le hash SHA-256 reçu de Go.
2. **Module IA :** Envoi du contexte (ex: "Processus inconnu qui spawn cmd.exe depuis un dossier temporaire") à une IA (soit un modèle ML local, soit une API LLM) pour avoir un *Anomaly Score*.
3. **Risk Engine :** Un algorithme combine les deux résultats. Si le `Risk Score > 75/100`, il génère une alerte Critique.


---s


## 2. Flux d'Approbation des Actions de Remédiation (Human-in-the-Loop)

### 🎯 Le Concept
Un EDR ne doit pas être un "cowboy" qui détruit des processus au hasard et risque de casser le PC de l'utilisateur. Toute action critique de remédiation (comme tuer un processus, isoler la machine du réseau, ou supprimer un fichier) doit passer par un système d'autorisation granulaire. 
C'est le concept du **"Human-in-the-Loop"**. L'Administrateur du SOC garde le contrôle total via le Dashboard et définit les permissions.

### 🛠️ Implémentation Technique

#### Partie 1 : Le Backend (Gestionnaire des Permissions)
- **Nouvelle Table SQL `remediation_requests`** : Contient l'action demandée par le Risk Engine (ex: `KILL_PROCESS`, PID: 1234), l'ID de l'alerte liée, et le statut actuel (`PENDING`, `APPROVED`, `REJECTED`).
- **Table `admin_settings`** : Stocke les niveaux de permissions. L'admin peut choisir quelles actions nécessitent une validation manuelle et lesquelles peuvent être automatiques (ex: "Auto-kill si Score > 95", mais "Demande d'autorisation obligatoire si Score entre 75 et 95").

#### Partie 2 : Le Dashboard (L'Interface de Validation)
- L'Administrateur reçoit une notification interactive : *"Le Risk Engine demande l'autorisation de tuer le processus unknown.exe (Score de Risque: 86). [AUTORISER] [REFUSER]"*.
- Une page "Settings" permet de configurer finement les autorisations selon le type d'attaque, la sévérité, ou les groupes de PC.

#### Partie 3 : L'Agent Go (Le Soldat)
- Au lieu de tuer le processus de sa propre initiative, l'Agent Go interroge le Backend régulièrement (ou utilise des WebSockets pour du vrai temps réel) pour vérifier s'il y a des "Ordres Approuvés" dans sa file d'attente.
- S'il reçoit l'ordre `KILL` officiellement approuvé pour le PID 1234, il exécute la commande système (`process.Kill()`) et renvoie un accusé de réception de succès au Backend.


---


## 3. Détection intelligente des connexions réseau suspectes (Corrélation)

### 🎯 Le Concept
Au lieu d'analyser le réseau de façon isolée (comme un pare-feu classique), l'EDR effectue une **Corrélation d'Événements**. Il lie chaque connexion réseau au processus exact qui l'a générée.
L'analyse se fait selon un double prisme :
1. **Threat Intelligence :** "L'IP de destination ou le domaine est-il connu comme malveillant par la communauté ?"
2. **IA (Analyse Comportementale) :** "Est-ce normal que *ce* processus spécifique (ex: `notepad.exe`) communique vers internet ? Est-ce normal qu'il scanne 50 IPs d'un coup ?"

### 🛠️ Implémentation Technique

#### Partie 1 : L'Agent Go (Le Collecteur Réseau)
- L'Agent utilise `net.Connections("tcp")` et `net.Connections("udp")` de gopsutil.
- Il extrait la "Télémétrie Réseau" complète : 
  - `Processus Associé` (Nom et PID)
  - `IP et Port de Destination`
  - `Statut` (ex: ESTABLISHED)
- L'Agent envoie ce paquet de connexion au Backend en l'associant au processus parent.

#### Partie 2 : Le Backend & Threat Intel
- **API Threat Intel :** Le Backend interroge des API de confiance.
  - *AbuseIPDB* ou *AlienVault OTX* : Vérifie la réputation de l'IP.
- L'API retourne un niveau de dangerosité, et les tags éventuels (ex: "Botnet C2", "Tor Exit Node").

#### Partie 3 : Corrélation et Machine Learning
- Le **Risk Engine** combine les données de la Détection N°1 (Score du Processus) avec la Détection N°3 (Score de l'IP).
- *La Puissance de la Corrélation :* Si `update.exe` a un score légèrement suspect (ex: 40/100) mais qu'il se connecte à une IP blacklistée, le score combiné explose instantanément à 95/100.
- *Anomalie de Fréquence (IA) :* Le modèle IA analyse le rythme des requêtes. Un "Port Scan" interne ou un logiciel cherchant des bases de données déclenchera une anomalie comportementale, même si les IPs sont techniquement "saines".


---


## 4. File Activity Monitoring & Ransomware Behavior Detection

### 🎯 Objectif
Surveiller en temps réel les actions effectuées sur les fichiers du PC afin de détecter des comportements anormaux ou potentiellement malveillants, notamment ceux qui ressemblent à une activité de ransomware (chiffrement, renommage ou suppression massive).

### 🛠️ Implémentation Technique

#### 1. Données collectées par l’agent (Golang)
L'Agent Go agit comme un FIM (File Integrity Monitor). Pour chaque activité, il récupère :
* Type d’action : Création, Modification, Renommage, Suppression.
* Chemin du fichier et extension.
* Processus responsable et son PID.
* Timestamp et **nombre d’actions effectuées par seconde**.

#### 2. Politique configurable par l’administrateur (Policy Engine)
Depuis le dashboard Next.js, l’administrateur définit lui-même les seuils. Exemple :
- Maximum fichiers modifiés : 30 / seconde.
- Fenêtre d’analyse : 10 secondes.
Ces valeurs sont synchronisées dynamiquement avec l'Agent Go par le Backend, sans avoir besoin de recompiler l'agent.

#### 3. Deuxième couche : Intelligence Artificielle
En parallèle des règles strictes (Policy), le modèle IA analyse si le comportement est inhabituel par rapport à l'historique de *ce* processus spécifique.
Par exemple, si `backup.exe` modifie 10 fichiers/s d'habitude et passe soudainement à 73/s, l'IA détecte une *AI Behavioral Anomaly* (même si le seuil de l'Admin n'est pas encore dépassé).

#### 4. Fusion des deux résultats (Risk Engine)
Le **Risk Engine** combine les deux moteurs : `Policy Detection` + `AI Anomaly Detection`.
Si les deux remontent une alerte (Policy=HIGH + IA Anomaly=0.96), le Risk Score explose à 97/100 et la sévérité passe en CRITICAL.

#### 5. Actions de Réponse Automatisée
L'administrateur choisit les actions selon la gravité, et pour la criticité maximale liée au chiffrement de fichiers, il peut activer l'option "Terminer le processus" automatiquement (qui passera outre la validation humaine pour sauver les données en une fraction de seconde).


---


## 5. Persistence Monitoring & Admin Validation

### 🎯 Objectif
Détecter lorsqu’un programme essaie de s'installer de manière persistante sur la machine (lancement automatique après un redémarrage ou une reconnexion).
Plutôt que de bloquer automatiquement toute nouvelle persistance (ce qui casserait les mises à jour légitimes de l'OS), le système utilise une approche intelligente : **Detect → Analyze → Ask Admin → Respond → Learn**.

### 🛠️ Implémentation Technique

#### 1. L'Agent Go (Le Veilleur de Démarrage)
L'Agent surveille les points clés de l'OS :
- Nouveaux programmes au démarrage (Dossier Startup, Clés de Registre).
- Nouvelles tâches planifiées.
- Nouveaux services Windows.

Pour chaque événement, l'Agent récolte un **Profil de Persistance complet** :
- Processus responsable et Commande exacte exécutée.
- Chemin du fichier et Hash SHA-256.
- **Signature numérique et Éditeur** (Crucial : ex: "Microsoft Corporation" vs "Non Signé").

#### 2. Analyse Côté Backend (IA + Threat Intel)
Le Backend reçoit l'événement et l'enrichit :
- **Threat Intel** : Le hash est-il connu comme malveillant ?
- **Contexte** : Quel est le "Process Risk" actuel du programme qui a créé la persistance ?
- **IA** : Ce type de persistance est-il habituel pour cet éditeur ?
Le Backend génère un *Persistence Risk Score*.

#### 3. Décision Humaine (Human-in-the-Loop)
L'administrateur reçoit une alerte détaillée avec 3 boutons d'action :
- **ALLOW (Autoriser)** : Persistance légitime. Le système l'ajoute à une politique de confiance.
- **INVESTIGATE (Investiguer)** : Conserver l'alerte pour analyse approfondie.
- **BLOCK (Bloquer)** : Ordre envoyé à l'Agent Go de désactiver ou supprimer la persistance immédiatement.

#### 4. Le Moteur d'Apprentissage (Trust Policy)
Point de génie : Si l'Admin clique sur "ALLOW", l'IA n'apprend pas simplement que "updater.exe = gentil". 
Le système crée une règle de confiance forte basée sur une **empreinte composite** : 
`[Hash du Binaire] + [Signature Valide de l'Éditeur] + [Chemin d'accès] + [Type de persistance]`.
Ainsi, si un malware se renomme "updater.exe", son Hash ou sa Signature sera différent, et l'EDR le bloquera.


---


## 6. System & Security Event Monitoring (Windows Logs)

### 🎯 Objectif
Surveiller en temps réel les **Logs Système et de Sécurité** de la machine (les fameux *Windows Event Logs*) pour y détecter des séquences d'événements malveillants (ex: Brute Force, Usurpation d'identité, Escalade de Privilèges).
Nexus-EDR ne se contente pas d'être un simple lecteur de logs, il les transforme en détections actives et les corrèle.

### 🛠️ Implémentation Technique

#### 1. L'Agent Go (Le Parseur d'Événements)
L'Agent utilise l'API native de l'OS (ex: WMI ou Event Viewer sous Windows) pour collecter et normaliser les événements critiques :
- Connexions (Logons) réussies et échouées.
- Verrouillages de compte.
- Création/Suppression d'utilisateurs locaux.
- Ajout d'un utilisateur à un groupe sensible (ex: Administrateurs).
- Modifications critiques des paramètres de sécurité.

#### 2. Politiques Configurables (Policy Engine)
L'administrateur définit des seuils sur le Dashboard (ex: `Échecs de connexion : 10/5min = WARNING, 50/5min = CRITICAL`). 
Ces politiques sont transmises à l'Agent.

#### 3. Analyse Intelligente (IA Comportementale)
L'IA apprend la "Baseline" (la routine) de chaque utilisateur.
- *Exemple :* `user01` se connecte habituellement en journée depuis `PC-01`. 
- Si `user01` se connecte à 3h20 du matin sur un nouveau `PC-08` après 18 échecs, l'IA génère un **Anomaly Score très élevé** (ex: 0.93), même si le seuil strict de "50 échecs" de l'Admin n'est pas atteint.

#### 4. La Corrélation (Reconstruction de la "Kill Chain")
C'est ici que ce module brille. Le Risk Engine peut désormais corréler une attaque complète de bout en bout :
1. *Module 6 (Logs)* : 18 échecs de connexion, suivis d'une connexion réussie (Attaque Brute Force probable).
2. *Module 6 (Logs)* : Privilèges Administrateur obtenus.
3. *Module 1 (Processus)* : Un processus inconnu (`unknown.exe`) est lancé par ce compte fraîchement compromis.
4. *Module 3 (Réseau)* : Ce processus se connecte à une IP de commande (C2) malveillante.
5. *Module 4 (Fichiers)* : Début d'une altération massive de fichiers (Comportement Ransomware).
👉 Grâce à la corrélation de ces 5 événements, le **Risk Score devient instantanément CRITICAL (100/100)** et la machine est isolée.


---


## L'Architecture d'Intelligence Artificielle (Le Cerveau Central)

### 🎯 Objectif : L'Immunité Collective (Crowd Immunity)
Plutôt que d'avoir une petite IA isolée sur chaque PC qui n'apprend que de l'utilisateur local, Nexus-EDR utilisera une **IA Centralisée sur le Backend**. 
Dès que l'Agent EDR est déployé chez 10, 100 ou 10 000 utilisateurs, toute leur télémétrie (comportements normaux et attaques subies) remonte vers l'API centrale pour enrichir le même modèle.

### 🛠️ Le Modèle "Fait Maison" (Souveraineté des Données)
Nous allons coder notre propre moteur de Machine Learning en Python (ex: `scikit-learn`). Ce modèle sera hébergé sur le serveur central. 

1. **Apprentissage Continu (Continuous Learning) :** Le modèle se nourrit en permanence de la base de données globale PostgreSQL. Si un utilisateur subit une attaque d'un malware inconnu ou d'une nouvelle technique, le modèle enregistre ce modèle comportemental.
2. **L'Immunité Collective :** Dès que la nouvelle méthode d'attaque est comprise par le Cerveau Central, l'IA protège instantanément *tous les autres utilisateurs* du réseau, car les agents interrogent tous la même base centrale.


---


## L'Écosystème d'Enrichissement (Threat Intelligence APIs)

Pour que notre "Cerveau Central" (Backend) puisse prendre des décisions intelligentes instantanées en s'appuyant sur l'expertise mondiale, il interrogera des bases de données de **Threat Intelligence** gratuites et open-source. Le Backend fera ces requêtes HTTP (asynchrones) dès qu'il recevra de la télémétrie de l'Agent.

### 1. APIs de Réputation de Fichiers (Hashes SHA-256)
- **VirusTotal API :** (Jusqu'à 500 req/jour gratuits). Vérifie le hash contre 70+ moteurs antivirus mondiaux. La référence absolue.
- **MalwareBazaar (abuse.ch) :** (Gratuit, illimité). Excellente base de données pour identifier la *famille* exacte d'un malware et ne pas juste avoir un score binaire.
- **Hybrid Analysis API :** (Gratuit). Fournit des rapports de "Sandboxing" comportementaux très détaillés pour un Hash.

### 2. APIs de Réputation Réseau (IPs & Domaines)
- **AbuseIPDB :** (Jusqu'à 1000 req/jour). L'API de référence pour savoir si une adresse IP participe à un réseau de Botnets ou fait des attaques Brute Force. Retourne un `Abuse Score` sur 100.
- **URLhaus (abuse.ch) :** (Gratuit, illimité). Spécialisée dans la détection d'URLs distribuant des malwares (souvent appelées par des scripts malveillants).
- **AlienVault OTX (Open Threat Exchange) :** (Gratuit, illimité). Le plus grand réseau communautaire. Permet de vérifier n'importe quel "IoC" (Indicateur de Compromission) pour voir si la communauté mondiale l'a déjà signalé.

### 3. APIs de Contexte (Analyse de la cible)
- **Shodan API :** (Plan gratuit limité). Le "moteur de recherche des hackers". Permet de savoir si l'IP avec laquelle notre processus communique héberge un service étrange (ex: un serveur de Commande et Contrôle C2 ouvert).
- **GreyNoise (Community API) :** (Gratuit). Filtre le "bruit de fond" d'Internet. Indispensable pour éviter que notre EDR s'affole à chaque fois qu'un scanner de recherche bénin (comme ceux de Google) "touche" notre machine.
- **IPinfo / ip-api.com :** (Gratuit). API de Géolocalisation. Savoir qu'un processus système caché se connecte soudainement à un serveur en Russie ou en Corée du Nord est un modificateur de score massif pour le Risk Engine.

### 4. Détections "Gratuites" via Corrélation d'API
Grâce au simple croisement de la télémétrie locale (Hash + IP) avec ces API, notre Backend est capable de détecter des menaces complexes sans aucun code supplémentaire dans l'Agent Go :
- **Cryptojacking (Minage caché) :** L'API AbuseIPDB/Shodan tagge l'IP de destination comme "Mining Pool" (Bitcoin/Monero).
- **Trafic Dark Web (Tor) :** L'API AlienVault indique que l'IP de destination est un "Tor Exit Node", un comportement typique des ransomwares pour cacher leur C2.
- **Serveurs Cobalt Strike / C2 :** L'API Shodan identifie un serveur offensif ouvert sur l'IP distante.
- **Gestion des Vulnérabilités (CVE) :** L'Agent envoie le Hash d'un fichier légitime. VirusTotal répond qu'il n'y a pas de virus, MAIS que ce fichier possède une vulnérabilité critique connue (CVE). L'EDR sert alors de scanner de failles.
- **Proxies/VPN Malveillants :** L'API IPinfo détecte que l'IP distante appartient à un fournisseur de VPN anonyme louche (utilisé pour l'exfiltration de données).
