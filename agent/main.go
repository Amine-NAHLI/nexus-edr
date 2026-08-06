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

// Nouvelle fonction magique pour calculer les paliers de 5% (80, 85, 90...)
func getPalier(valeur float64) float64 {
	return float64(int(valeur/5) * 5)
}

func main() {
	fmt.Println("🚀 Agent EDR Actif. Surveillance Anti-Spam (Par Paliers) en cours...")
	
	agentName := "MON-PC-WINDOWS"
	
	// On mémorise le dernier palier alerté. On commence à 75.
	dernierPalierAlerte := 75.0

	for {
		// 1. Lire le processeur (CPU)
		cpuPercent, _ := cpu.Percent(0, false)
		
		// 2. Lire la mémoire (RAM)
		vMem, _ := mem.VirtualMemory()
		
		currentCpu := cpuPercent[0]
		currentRam := vMem.UsedPercent

		fmt.Printf("🔍 Analyse - CPU: %.2f%% | RAM: %.2f%%\n", currentCpu, currentRam)

		// On cherche la valeur la plus haute entre le CPU et la RAM
		maxUtilisation := currentCpu
		if currentRam > maxUtilisation {
			maxUtilisation = currentRam
		}

		// On calcule à quel "palier" de 5% on se trouve (ex: 84 -> 80, 87 -> 85)
		palierActuel := getPalier(maxUtilisation)

		// 3. Détecter une anomalie : on dépasse 80% ET on a franchi un NOUVEAU palier !
		if palierActuel >= 80.0 && palierActuel > dernierPalierAlerte {
			fmt.Printf("⚠️ ALERTE ROUGE : Nouveau palier critique franchi (%.0f%%) !\n", palierActuel)
			
			alert := TelemetryData{
				AgentID:                agentName,
				CpuPercent:             currentCpu,
				RamPercent:             currentRam,
				SuspiciousProcessFound: true,
				IpDestination:          "IP INCONNUE", 
			}
			sendAlert(alert)
			
			// On met à jour le dernier palier pour ne plus spammer la base de données
			dernierPalierAlerte = palierActuel
			
		} else if maxUtilisation < 75.0 {
			// Si le PC se calme (en dessous de 75%), on réinitialise le système d'alerte !
			dernierPalierAlerte = 75.0
		}

		time.Sleep(3 * time.Second)
	}
}
