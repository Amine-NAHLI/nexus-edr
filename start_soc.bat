@echo off
title Nexus EDR - Lanceur
echo ===================================
echo     Demarrage de Nexus EDR...
echo ===================================

echo [1/3] Demarrage de la base de donnees PostgreSQL...
docker start nexus-db

echo [2/3] Demarrage du Backend Python...
start "Nexus EDR - Backend" cmd /k "cd backend && venv\Scripts\activate && uvicorn main:app --reload"

echo [3/3] Demarrage du Dashboard Next.js...
start "Nexus EDR - Frontend" cmd /k "cd frontend && npm run dev"

echo.
echo ✅ Tous les services sont en cours de demarrage !
echo ------------------------------------------------
echo - API Backend : http://127.0.0.1:8000
echo - API Docs    : http://127.0.0.1:8000/docs
echo - Dashboard   : http://localhost:3000
echo ------------------------------------------------
echo Vous pouvez reduire cette fenetre.
pause
