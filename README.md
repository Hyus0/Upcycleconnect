# UpcycleConnect - Refonte SI (PA)

Socle de demarrage conforme au CDC fourni pour les 3 lots, avec implementation initiale de la Mission 1.

## Structure

- `mission-1-app/web-app`: Frontend JavaScript (espaces Particuliers, Pro/Artisans, Salaries, Admin).
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

- Front JavaScript: OK (MVP fonctionnel).
- Back PHP: OK (routes principales, stockage JSON, metriques admin).
- API Go: OK (score/validation).
- Stripe/OneSignal/PDF: points d'integration prepares, implementation complete a brancher dans l'iteration suivante.
