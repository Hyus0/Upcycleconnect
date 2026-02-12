# Documentation technique - Stack applicative

## Architecture runtime

- Nginx `gateway` expose `:8080`.
- `gateway` sert le front statique et reverse proxy:
  - `/api/*` -> backend PHP.
  - `/api-go/*` -> API Go.
- Backend PHP stocke les donnees dans `mission-1-app/backoffice/storage/data.json`.

## Endpoints principaux

### Backoffice PHP

- `GET /health`
- `GET /api/conseils`
- `GET /api/formations`
- `GET|POST /api/annonces`
- `GET|POST /api/planning`
- `POST /api/container-deposits`
- `POST /api/pro/subscriptions`
- `POST /api/pro/projects`
- `POST /api/employee/trainings`
- `POST /api/news`
- `POST /api/employee/moderation`
- `GET /api/admin/metrics`
- `GET /api/admin/users`
- `POST /api/admin/notifications`

### API Go

- `GET /health`
- `POST /api/upcycling-score`
- `POST /api/barcode/validate`
- `POST /api/notifications/preview`

## Points d'integration a realiser ensuite

- Stripe checkout + webhook paiement.
- OneSignal push transactionnelles et admin.
- Generation PDF (factures, attestations formations, bordereaux depot).
- Base SQL (migration depuis JSON) + authentification JWT + RBAC.
