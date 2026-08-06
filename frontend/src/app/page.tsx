"use client"; // Indique à Next.js que c'est une page interactive (côté client)
import { useEffect, useState } from "react";

// On définit le "Contrat" (comme Pydantic en Python, mais pour TypeScript)
interface TelemetryData {
  agent_id: string;
  cpu_percent: number;
  ram_percent: number;
  suspicious_process_found: boolean;
  ip_destination: string | null;
}

export default function Home() {
  // La mémoire locale de notre page : une liste d'alertes, vide au début
  const [alerts, setAlerts] = useState<TelemetryData[]>([]);

  // La fonction qui va interroger l'API Python
  const fetchAlerts = async () => {
    try {
      const response = await fetch("http://127.0.0.1:8000/api/alerts");
      const data = await response.json();
      setAlerts(data); // On sauvegarde les données reçues
    } catch (error) {
      console.error("Erreur de connexion au Backend:", error);
    }
  };

  // useEffect : le moteur de la page.
  // Il lance la recherche au démarrage, puis recommence toutes les 2 secondes !
  useEffect(() => {
    fetchAlerts(); 
    const interval = setInterval(fetchAlerts, 2000); 
    return () => clearInterval(interval);
  }, []);

  return (
    <main className="dashboard-container">
      <div className="header">
        <h1>NEXUS <span>EDR</span></h1>
        <p>Security Operations Center</p>
      </div>

      <div className="grid">
        {/* Boîte 1 : Statut du parc */}
        <div className="card">
          <h2>Statut Global</h2>
          <div className="status online">Système Opérationnel</div>
          <p>{alerts.length > 0 ? "⚠️ Des attaques sont en cours !" : "✅ Tout est calme"}</p>
        </div>

        {/* Boîte 2 : Les Alertes Dynamiques */}
        <div className="card alert-card">
          <h2>🚨 Alertes Critiques ({alerts.length})</h2>
          
          {/* S'il n'y a pas d'alerte */}
          {alerts.length === 0 && <p style={{color: "gray"}}>Aucune menace détectée.</p>}
          
          {/* On crée une boîte rouge pour CHAQUE alerte reçue (La boucle "map") */}
          {alerts.map((alerte, index) => (
            <div key={index} className="alert-item">
              <strong>{alerte.agent_id}</strong>
              <p>CPU à {alerte.cpu_percent}% - RAM à {alerte.ram_percent}%</p>
              
              {/* Si une IP louche a été trouvée, on l'affiche */}
              {alerte.ip_destination && (
                <p style={{ color: "var(--neon-red)", fontWeight: "bold" }}>
                  🌐 IP Cible : {alerte.ip_destination}
                </p>
              )}
              
              <button className="btn-block">Bloquer l'IP</button>
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}
