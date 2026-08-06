from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, declarative_base

# 1. L'URL de connexion au serveur PostgreSQL
# Format : postgresql://utilisateur:mot_de_passe@serveur:port/nom_de_la_base
# (psycopg2 est le pilote que Python utilise pour faire transiter les infos)
SQLALCHEMY_DATABASE_URL = "postgresql+psycopg2://postgres:admin@localhost:5432/postgres"

# 2. Le Moteur (Le traducteur)
engine = create_engine(SQLALCHEMY_DATABASE_URL)

# 3. Le créateur de Sessions
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# 4. La Base pour nos modèles
Base = declarative_base()
