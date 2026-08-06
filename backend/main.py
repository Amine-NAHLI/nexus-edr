from fastapi import FastAPI, Depends, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy.orm import Session
from typing import List

# On importe nos nouveaux fichiers architecturaux !
import models
import schemas
from database import engine, SessionLocal

# 1. CRÉATION AUTOMATIQUE DES TABLES
# Si la table "alerts" n'existe pas dans PostgreSQL, SQLAlchemy la crée tout seul
models.Base.metadata.create_all(bind=engine)

app = FastAPI(title="Nexus EDR API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"], 
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 2. LE DISTRIBUTEUR DE SESSIONS
# À chaque fois qu'un utilisateur appelle l'API, on lui ouvre une connexion (un panier)
# et on la ferme proprement à la fin.
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

@app.get("/")
def read_root():
    return {"status": "success", "message": "Backend Connecté à PostgreSQL !"}

# --- ROUTE 1 : RECEVOIR UNE ALERTE ---
# On utilise notre schéma "AlertCreate"
@app.post("/api/telemetry")
def receive_telemetry(data: schemas.AlertCreate, db: Session = Depends(get_db)):
    if data.suspicious_process_found:
        print(f"🚨 ALERTE ROUGE sur {data.agent_id} !! Sauvegarde en base...")
        
        # On transforme la donnée reçue en Modèle pour la base de données
        nouvelle_alerte = models.Alert(
            agent_id=data.agent_id,
            cpu_percent=data.cpu_percent,
            ram_percent=data.ram_percent,
            suspicious_process_found=data.suspicious_process_found,
            ip_destination=data.ip_destination
        )
        
        # On sauvegarde officiellement (INSERT INTO alerts ...)
        db.add(nouvelle_alerte)
        db.commit()
        
        return {"status": "alert_saved"}
        
    return {"status": "ok"}


# --- ROUTE 2 : LIRE LES ALERTES ---
# On renvoie une liste de "AlertResponse" (qui contient l'ID)
@app.get("/api/alerts", response_model=List[schemas.AlertResponse])
def get_alerts(db: Session = Depends(get_db)):
    # L'équivalent SQL : SELECT * FROM alerts;
    alertes = db.query(models.Alert).all()
    return alertes


# --- ROUTE 3 : SUPPRIMER UNE ALERTE (LA NOUVEAUTÉ !) ---
# Le Dashboard va nous envoyer un paramètre (ex: /api/alerts/5)
@app.delete("/api/alerts/{alert_id}")
def delete_alert(alert_id: int, db: Session = Depends(get_db)):
    # On cherche l'alerte qui a cet ID précis
    alerte = db.query(models.Alert).filter(models.Alert.id == alert_id).first()
    
    if alerte is None:
        raise HTTPException(status_code=404, detail="Alerte introuvable")
        
    # On supprime (DELETE FROM alerts WHERE id = 5)
    db.delete(alerte)
    db.commit()
    
    return {"status": "success", "message": "Alerte supprimée !"}
