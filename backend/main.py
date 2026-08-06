from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional, List

# 1. Le Contrat
class TelemetryData(BaseModel):
    agent_id: str
    cpu_percent: float
    ram_percent: float
    suspicious_process_found: bool
    ip_destination: Optional[str] = None

app = FastAPI(title="Nexus EDR API")

# --- NOUVEAUTÉ 1 : CORS ---
# On autorise expressément ton Dashboard (port 3000) à lire les données
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"], 
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# --- NOUVEAUTÉ 2 : Base de données temporaire ---
# Une liste vide qui va garder nos alertes en mémoire
alertes_memoire = []

@app.get("/")
def read_root():
    return {"status": "success", "message": "Bienvenue sur le Backend de Nexus EDR !"}

# Route POST : L'Agent Go envoie les données ici
@app.post("/api/telemetry")
def receive_telemetry(data: TelemetryData):
    print(f"📡 Données reçues du PC {data.agent_id}: CPU={data.cpu_percent}%")
    
    if data.suspicious_process_found:
        print(f"🚨 ALERTE ROUGE sur le PC {data.agent_id} !!")
        # On SAUVEGARDE l'alerte dans notre liste
        alertes_memoire.append(data)
        return {"status": "alert_registered", "action": "block_process"}
        
    return {"status": "ok", "message": "Données enregistrées"}

# --- NOUVEAUTÉ 3 : Route GET (La route de lecture) ---
# Le Dashboard Next.js va appeler cette route pour récupérer la liste
@app.get("/api/alerts")
def get_alerts():
    return alertes_memoire
