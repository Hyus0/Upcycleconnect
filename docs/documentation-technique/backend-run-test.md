# Backend UpcycleConnect — lancer & tester

Date : 2026-06-24

Guide court pour démarrer, tester et maintenir le backend après les travaux de professionnalisation.

## 1. Prérequis

- Docker Desktop démarré.
- Node/npm ne sont **pas** requis sur l'hôte : les builds front se font via Docker.

## 2. Lancer la stack

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect\deploy\docker-compose
docker compose up -d --build
docker compose ps
```

Points d'entrée :

- Site public : http://localhost:8088/
- Back-office admin : http://localhost:8088/admin/
- phpMyAdmin : http://localhost:8089/
- Santé API Go : http://localhost:8088/health/api-go
- Santé backoffice PHP : http://localhost:8088/health/backoffice

## 3. Comptes de démonstration (seed)

Mots de passe stockés en **bcrypt** ; valeurs en clair pour les tests :

| Rôle | Email | Mot de passe |
|---|---|---|
| Admin | admin@upcycleconnect.local | Admin123! |
| Particulier | particulier@upcycleconnect.local | Test123! |
| Prestataire | prestataire@upcycleconnect.local | Test123! |
| Salarié | salarie@upcycleconnect.local | Test123! |

Le back-office n'accepte que les comptes **Admin**.

## 4. Architecture des appels (gateway Nginx)

Le navigateur ne parle qu'au gateway (port 8088). Règles de routage :

- `/go/*` → API Go (le préfixe `/go` est retiré). Ex. `/go/annonces` → `/annonces`.
- `/api-go/*` → API Go sous `/api/*` (back-office). Ex. `/api-go/admin/users` → `/api/admin/users`.
- `/api/*` → backoffice PHP historique.

Le front public utilise une base d'URL unique `"/go"` (plus aucun `localhost:8081` en dur).

## 5. Authentification & autorisations

- Connexion : `POST /go/login` `{email, password}` → renvoie un `token` (stocké en session côté front, colonne `UTILISATEUR.token`).
- Les requêtes protégées envoient l'en-tête `Authorization: <token>`.
- Middlewares Go (`app/auth.go`) :
  - `RequireRole("Admin")` : verrouille toutes les routes `/api/admin/*` (401 sans token, 403 si non-Admin).
  - `RequireSelf("id")` : routes per-utilisateur (panier, checkout, factures, abonnement, messagerie, matériaux-recherchés) — un utilisateur n'accède qu'à ses propres données ; les Admin passent.

## 6. Tests rapides (PowerShell)

```powershell
# Public
curl.exe http://localhost:8088/go/annonces

# Login (recupere un token)
'{"email":"admin@upcycleconnect.local","password":"Admin123!"}' | Set-Content t.json -Encoding ascii -NoNewline
curl.exe -s -X POST http://localhost:8088/go/login -H "Content-Type: application/json" --data-binary "@t.json"

# Admin protege : 401 sans token, 200 avec token Admin
curl.exe -s -o NUL -w "%{http_code}`n" http://localhost:8088/api-go/admin/users
curl.exe -s -H "Authorization: <TOKEN_ADMIN>" http://localhost:8088/api-go/admin/users
```

## 7. Rebuild après modification

```powershell
# API Go
cd deploy\docker-compose
docker compose up -d --build api-go

# Front public
cd ..\..\mission-1-app\web-app
docker run --rm -v "${PWD}:/app" -v /app/node_modules -w /app node:20-alpine sh -lc "npm ci && npm run build"

# Back-office admin
cd ..\..\admin
docker run --rm -v "${PWD}:/app" -v /app/node_modules -w /app node:20-alpine sh -lc "npm ci && npm run build"
```

Les `dist/` sont servis directement par Nginx (pas besoin de redémarrer le gateway).

## 8. Base de données

- Source de vérité du schéma : `deploy/mysql/init/` (régénéré depuis la base réelle, 39 tables).
  - `001_schema.sql` : schéma complet.
  - `002_admin_extensions.sql` : neutralisé (colonnes intégrées dans 001).
  - `003_bootstrap.sql` : seed idempotent (comptes bcrypt + catalogue).
- Connexion de l'API via variables d'env Docker (`DB_HOST=mysql`, `DB_NAME=upcycletest`, …) avec `Ping` au démarrage.

Réinitialiser complètement la base (⚠️ perte des données locales, sauvegarder avant) :

```powershell
cd deploy\docker-compose
docker compose exec -T mysql mysqldump -uroot -ppassword upcycletest > ..\..\temp\backup-avant-reset.sql
docker compose down
docker volume rm docker-compose_mysql-data
docker compose up -d --build
```

## 9. Limites connues / pistes restantes

- **Module finance** : agrège désormais les `COMMANDE` (montant_total / statut). Les montants restent à 0 tant qu'aucune commande réelle n'est passée.
- **Durcissement public** : l'ownership est appliqué aux routes les plus sensibles ; les autres routes per-utilisateur (lecture de stats, etc.) restent à durcir au cas par cas.
- **Stripe / OneSignal / PDF** : intégrations préparées, non branchées.
