# UpcycleConnect Admin

Interface d'administration Vue.js dediee au projet UpcycleConnect.

## Installation

```bash
cd admin
npm install
npm run dev
```

Build production :

```bash
npm run build
```

Tests :

```bash
npm run test
```

## Variables d'environnement

Copier `.env.example` puis ajuster si besoin :

- `VITE_BACKOFFICE_API_BASE` : base du backoffice PHP, par defaut `http://localhost:8080/api`
- `VITE_GO_API_BASE` : base de l'API Go, par defaut `http://localhost:8080/api-go`

## Structure

- `src/pages/` : ecrans admin
- `src/components/` : composants reutilisables
- `src/services/` : couche API centralisee
- `src/router/` : routes admin
- `src/store/` : store local des toasts
- `src/utils/` : formatage

## Endpoints utilises

### Dashboard

- `GET /api/admin/metrics`
- `GET /api/admin/users`
- `GET /api/annonces`

### Utilisateurs

- Utilise actuellement : `GET /api/admin/users`
- Manquants pour CRUD complet :
  - `GET /api/admin/users/:id`
  - `POST /api/admin/users`
  - `PUT /api/admin/users/:id`
  - `PATCH /api/admin/users/:id/status`
  - `DELETE /api/admin/users/:id`

### Prestations

- Utilises actuellement :
  - `GET /api/annonces`
  - `POST /api/annonces`
- Manquants :
  - `PUT /api/annonces/:id`
  - `DELETE /api/annonces/:id`

### Categories de prestations

- Aucun endpoint disponible dans le repo courant
- Requis :
  - `GET /api/categories`
  - `GET /api/categories/:id`
  - `POST /api/categories`
  - `PUT /api/categories/:id`
  - `DELETE /api/categories/:id`

### Evenements

- Aucun endpoint disponible dans le repo courant
- Requis :
  - `GET /api/events`
  - `GET /api/events/:id`
  - `POST /api/events`
  - `PUT /api/events/:id`
  - `DELETE /api/events/:id`

## Fichiers crees

- `admin/package.json`
- `admin/vite.config.js`
- `admin/index.html`
- `admin/.env.example`
- `admin/src/main.js`
- `admin/src/App.vue`
- `admin/src/router/index.js`
- `admin/src/services/http.js`
- `admin/src/services/api.js`
- `admin/src/services/api.spec.js`
- `admin/src/store/toastStore.js`
- `admin/src/utils/format.js`
- `admin/src/components/*`
- `admin/src/pages/*`

## Notes d'integration

- Les couleurs, polices et tons d'action reprennent la charte `UpcycleConnect`.
- Les actions principales sont en vert, les alertes en ambre, les suppressions en corail.
- Les vues Categories et Evenements assument explicitement l'absence d'API plutot que de simuler des donnees backend.
