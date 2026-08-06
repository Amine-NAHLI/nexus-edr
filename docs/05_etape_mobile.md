# Étape 5 (Future) : L'Application Mobile

## 🎯 Objectif
Permettre à l'administrateur de recevoir des alertes en temps réel sur son téléphone (iOS et Android) et de prendre des décisions critiques sans avoir besoin d'ouvrir son ordinateur.

## 🛠️ Pourquoi React Native ?
React Native permet de coder une seule fois (en JavaScript/TypeScript) pour déployer l'application à la fois sur l'App Store (Apple) et le Google Play Store. De plus, comme nous utilisons déjà Next.js (React) pour le Dashboard web, nous pourrons réutiliser une grande partie de la logique et du code métier !

## 📝 Ce que fera l'application :

1. **Notifications Push (Firebase / APNs) :**
   Si le Moteur IA détecte une menace grave (ex: tentative de suppression de fichiers système), le téléphone vibre instantanément. L'alerte s'affiche sur l'écran de verrouillage.

2. **Actions Rapides (Quick Actions) :**
   Depuis la notification, l'utilisateur pourra :
   - ✅ Ignorer (Faux Positif)
   - 🛑 Bloquer l'IP cible
   - 💀 Isoler le PC du réseau
   (Ces choix viendront nourrir l'algorithme de Reinforcement Learning).

3. **Vue d'Ensemble Simplifiée :**
   Un mini-dashboard affichant le nombre de PC en ligne et le score de santé global du parc informatique.

**👉 Statut actuel :** Cette étape sera développée après la validation de la plateforme Core (Agent, Backend, Dashboard, IA).
