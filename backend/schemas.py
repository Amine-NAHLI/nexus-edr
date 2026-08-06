from pydantic import BaseModel
from typing import Optional

# 1. Le contrat pour RECEVOIR l'attaque (Ancien TelemetryData)
class AlertCreate(BaseModel):
    agent_id: str
    cpu_percent: float
    ram_percent: float
    suspicious_process_found: bool
    ip_destination: Optional[str] = None

# 2. Le contrat pour ENVOYER l'attaque au Dashboard
# Il hérite de AlertCreate (il a donc toutes ses propriétés) + on lui rajoute l'ID !
class AlertResponse(AlertCreate):
    id: int 

    class Config:
        # Cette ligne magique permet à Pydantic de comprendre le format de SQLAlchemy
        from_attributes = True 
