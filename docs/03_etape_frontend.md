# Étape 3 : Le Dashboard (Next.js)

## 🎯 Objectif
Créer l'interface d'administration (SOC - Security Operations Center) belle, fluide et en temps réel pour gérer les menaces et la flotte d'Agents.

## 🛠️ Pourquoi Next.js & React ?
C'est le standard de l'industrie pour les applications Web (WebApps). Next.js gère le routage facilement. Combiné avec TailwindCSS (ou du Vanilla CSS bien architecturé), il permet de créer des designs "Premium" (Dark mode, effets de verre/Glassmorphism) indispensables pour un produit de Cybersécurité moderne.

## 📝 Ce que nous allons faire en pratique :

1. **Générer le projet Next.js :**
   Nous utiliserons le générateur officiel qui installe tout automatiquement (React, TypeScript, outils de build).
   *Commande :* `npx create-next-app@latest frontend --ts --app --src-dir --eslint`

2. **Créer le Design System (CSS) :**
   Nous mettrons en place les couleurs (Fonds très sombres, textes cyan et violet) et les animations fluides dans un fichier CSS global. L'objectif est l'effet "WOW".

3. **Créer les Pages (Routage App Router) :**
   ```text
   frontend/src/app/
   ├── layout.tsx         (Le menu latéral présent sur toutes les pages)
   ├── page.tsx           (La page d'accueil / Vue globale)
   ├── agents/
   │   └── page.tsx       (La liste des PC connectés)
   └── alerts/
       └── page.tsx       (Le centre d'alertes IA avec les boutons de feedback +1/-1)
   ```

4. **Interfacer avec le Backend :**
   Dans React, nous utiliserons la fonction `fetch()` ou la librairie `axios` pour interroger notre API Python (ex: `GET http://localhost:8000/api/agents`) et afficher les résultats dans de jolis tableaux et graphiques.

**👉 Et ensuite ?** Une fois ces 3 briques (Backend, Agent, Frontend) connectées, nous attaquerons la **Phase 2 : L'intégration de l'Intelligence Artificielle (Machine Learning).**
