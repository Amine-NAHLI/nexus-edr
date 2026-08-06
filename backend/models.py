from sqlalchemy import Column, Integer, String, Float, Boolean
from database import Base

class Alert(Base):
    __tablename__ = "alerts"

    id = Column(Integer, primary_key=True, index=True)
    agent_id = Column(String, index=True)
    cpu_percent = Column(Float)
    ram_percent = Column(Float)
    suspicious_process_found = Column(Boolean)
    ip_destination = Column(String, nullable=True)
