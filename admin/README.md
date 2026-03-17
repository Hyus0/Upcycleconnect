# UpcycleConnect Admin

Interface d'administration Vue.js pour UpcycleConnect avec mode local persistant.

## Installation

```bash
cd admin
npm install
npm run dev
```

Build :

```bash
npm run build
```

Tests :

```bash
npm run test
```

## Variables

- `VITE_BACKOFFICE_API_BASE` : `http://localhost:8080/api`
- `VITE_GO_API_BASE` : `http://localhost:8080/api-go`

## Base locale

- seed JSON : `admin/public/mock-db.json`
- schema locale : `admin/local-db/schema.sql`
- persistance navigateur : `localStorage`

L'application essaie d'abord les endpoints existants, puis bascule sur la base locale si l'API est indisponible.

## Ecrans

- Dashboard admin
- Utilisateurs CRUD local
- Prestations CRUD local + tentative d'appel API `/api/annonces`
- Categories CRUD local
- Evenements CRUD local

## Endpoints utilises

### Dashboard

- `GET /api/admin/metrics`
- `GET /api/admin/users`
- `GET /api/annonces`

### Utilisateurs

- `GET /api/admin/users`

### Prestations

- `GET /api/annonces`
- `POST /api/annonces`

## Endpoints toujours a prevoir cote backend

- `GET /api/admin/users/:id`
- `POST /api/admin/users`
- `PUT /api/admin/users/:id`
- `PATCH /api/admin/users/:id/status`
- `DELETE /api/admin/users/:id`
- `PUT /api/annonces/:id`
- `DELETE /api/annonces/:id`
- `GET /api/categories`
- `POST /api/categories`
- `PUT /api/categories/:id`
- `DELETE /api/categories/:id`
- `GET /api/events`
- `POST /api/events`
- `PUT /api/events/:id`
- `DELETE /api/events/:id`
