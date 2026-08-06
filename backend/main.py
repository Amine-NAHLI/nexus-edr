from fastapi import FastAPI
from pydantic import BaseModel
from typing import Optional

# 1. On définit le "Contrat" avec Pydantic
class TelemetryData(BaseModel):
    agent_id: str
    cpu_percent: float
    ram_percent: float
    suspicious_process_found: bool
    ip_destination: Optional[str] = None

# On initialise notre API
app = FastAPI(title="Nexus EDR API")

# On crée la porte d'entrée par défaut
@app.get("/")
def read_root():
    return {"status": "success", "message": "Bienvenue sur le Backend de Nexus EDR !"}

# 2. La nouvelle route qui écoute les envois de données (POST)
@app.post("/api/telemetry")
def receive_telemetry(data: TelemetryData):
    # Affichage dans le terminal pour débugger
    print(f"📡 Données reçues du PC {data.agent_id}: CPU={data.cpu_percent}%")
    
    # Si le PC nous dit qu'il a vu un processus bizarre
    if data.suspicious_process_found:
        print(f"🚨 ALERTE ROUGE sur le PC {data.agent_id} !!")
        return {"status": "alert_registered", "action": "block_process"}
        
    return {"status": "ok", "message": "Données enregistrées"}
