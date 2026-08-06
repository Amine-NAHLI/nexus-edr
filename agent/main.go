package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Le contrat Pydantic, mais version Go ! (On appelle ça une "Struct")
type TelemetryData struct {
	AgentID                string  `json:"agent_id"`
	CpuPercent             float64 `json:"cpu_percent"`
	RamPercent             float64 `json:"ram_percent"`
	SuspiciousProcessFound bool    `json:"suspicious_process_found"`
	IpDestination          string  `json:"ip_destination"`
}

func main() {
	fmt.Println("🚀 Agent EDR démarré. Préparation de l'envoi des données...")

	// On simule un PC attaqué avec un CPU à 99.9%
	fakeData := TelemetryData{
		AgentID:                "PC-TEST-GO-01",
		CpuPercent:             99.9,
		RamPercent:             45.0,
		SuspiciousProcessFound: true,
		IpDestination:          "198.51.100.23",
	}

	// On transforme nos données Go en format JSON universel
	jsonData, _ := json.Marshal(fakeData)

	fmt.Println("📡 Envoi de l'alerte vers le serveur Python...")
	
	// On lance l'attaque (la requête POST) vers l'API
	url := "http://127.0.0.1:8000/api/telemetry"
	reponse, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println("❌ Erreur : Serveur introuvable.")
		return
	}
	defer reponse.Body.Close()

	fmt.Println("✅ Le serveur Python a répondu avec le code HTTP :", reponse.StatusCode)
}
