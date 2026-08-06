export default function Home() {
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
          <p>1 PC Connecté</p>
        </div>

        {/* Boîte 2 : Les Alertes */}
        <div className="card alert-card">
          <h2>🚨 Alertes Critiques</h2>
          <div className="alert-item">
            <strong>PC-TEST-GO-01</strong>
            <p>CPU à 99.9% - Processus suspect détecté !</p>
            <button className="btn-block">Bloquer l'IP</button>
          </div>
        </div>
      </div>
    </main>
  );
}
