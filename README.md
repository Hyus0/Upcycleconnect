# UpcycleConnect - Refonte SI (PA)

Socle de demarrage conforme au CDC fourni pour les 3 lots, avec implementation initiale de la Mission 1.

## Structure

- `admin`: Frontend Vue.js d'administration servi en statique apres build.
- `mission-1-app/backoffice`: Backend PHP (API metier + administration).
- `mission-1-app/api-go`: API Go (upcycling score, validation code-barres, preview notifications).
- `deploy/docker-compose`: Orchestration Docker de l'application.
- `docs/`: Livrables projet et documentation.

## Lancer la stack

Pre-requis: Docker Desktop.

```powershell
./deploy/up.ps1 -Build
```

Puis ouvrir: `http://localhost:8080`

Health checks:
- `http://localhost:8080/health/backoffice`
- `http://localhost:8080/health/api-go`

## Etat d'implementation Mission 1

- Front Admin Vue.js: OK (servi depuis `admin/dist` dans la stack Docker).
- Back PHP: OK (routes principales, stockage JSON, metriques admin).
- API Go: OK (score/validation).
- Stripe/OneSignal/PDF: points d'integration prepares, implementation complete a brancher dans l'iteration suivante.
