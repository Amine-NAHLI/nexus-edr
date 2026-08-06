package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time" // Pour mettre le script en pause
)

type TelemetryData struct {
	AgentID                string  `json:"agent_id"`
	CpuPercent             float64 `json:"cpu_percent"`
	RamPercent             float64 `json:"ram_percent"`
	SuspiciousProcessFound bool    `json:"suspicious_process_found"`
	IpDestination          string  `json:"ip_destination"`
}

// Petite fonction maison pour tirer nos requêtes POST
func sendAlert(data TelemetryData) {
	jsonData, _ := json.Marshal(data)
	url := "http://127.0.0.1:8000/api/telemetry"
	reponse, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	
	if err != nil {
		fmt.Println("❌ Erreur de connexion pour", data.AgentID)
		return
	}
	defer reponse.Body.Close()
	fmt.Printf("✅ Attaque envoyée pour %s (Code %d)\n", data.AgentID, reponse.StatusCode)
}

func main() {
	fmt.Println("🚀 Agent EDR démarré. Lancement des simulations d'attaques...")

	// Attaque 1 : Cryptomining (CPU à fond)
	attaque1 := TelemetryData{
		AgentID:                "PC-COMPTA-01",
		CpuPercent:             99.9,
		RamPercent:             45.0,
		SuspiciousProcessFound: true,
		IpDestination:          "198.51.100.23",
	}

	// Attaque 2 : Fuite de données (RAM à fond vers une IP louche)
	attaque2 := TelemetryData{
		AgentID:                "SERVEUR-BASE-05",
		CpuPercent:             12.0,
		RamPercent:             98.5,
		SuspiciousProcessFound: true,
		IpDestination:          "45.33.22.11", 
	}

	// Attaque 3 : Ransomware total
	attaque3 := TelemetryData{
		AgentID:                "PC-PATRON-01",
		CpuPercent:             100.0,
		RamPercent:             100.0,
		SuspiciousProcessFound: true,
		IpDestination:          "9.9.9.9",
	}

	// On tire nos 3 missiles avec 1 seconde d'écart
	sendAlert(attaque1)
	time.Sleep(1 * time.Second)
	
	sendAlert(attaque2)
	time.Sleep(1 * time.Second)
	
	sendAlert(attaque3)

	fmt.Println("🏁 Toutes les attaques ont été envoyées au Backend !")
}
