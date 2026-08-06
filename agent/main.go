package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type TelemetryData struct {
	AgentID                string  `json:"agent_id"`
	CpuPercent             float64 `json:"cpu_percent"`
	RamPercent             float64 `json:"ram_percent"`
	SuspiciousProcessFound bool    `json:"suspicious_process_found"`
	IpDestination          string  `json:"ip_destination"`
}

func sendAlert(data TelemetryData) {
	jsonData, _ := json.Marshal(data)
	url := "http://127.0.0.1:8000/api/telemetry"
	reponse, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	
	if err != nil {
		fmt.Println("❌ Erreur de connexion au serveur")
		return
	}
	defer reponse.Body.Close()
	fmt.Printf("✅ Alerte envoyée (Code %d)\n", reponse.StatusCode)
}

func main() {
	fmt.Println("🚀 Agent EDR Actif. Surveillance du système en cours...")
	
	agentName := "MON-PC-WINDOWS"

	// Boucle infinie : on vérifie l'ordinateur toutes les 3 secondes
	for {
		// 1. Lire le VRAI processeur (CPU)
		cpuPercent, _ := cpu.Percent(0, false)
		
		// 2. Lire la VRAIE mémoire (RAM)
		vMem, _ := mem.VirtualMemory()
		
		currentCpu := cpuPercent[0]
		currentRam := vMem.UsedPercent

		fmt.Printf("🔍 Analyse - CPU: %.2f%% | RAM: %.2f%%\n", currentCpu, currentRam)

		// 3. Détecter une anomalie
		// Si ton processeur OU ta RAM dépasse 80%, c'est suspect !
		if currentCpu > 80.0 || currentRam > 80.0 {
			fmt.Println("⚠️ ANOMALIE DÉTECTÉE ! Surcharge du système !")
			
			alert := TelemetryData{
				AgentID:                agentName,
				CpuPercent:             currentCpu,
				RamPercent:             currentRam,
				SuspiciousProcessFound: true,
				IpDestination:          "IP INCONNUE", 
			}
			sendAlert(alert)
		}

		time.Sleep(3 * time.Second)
	}
}
