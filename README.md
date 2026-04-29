# UpcycleConnect

UpcycleConnect est le projet de PA de refonte SI autour de l'upcycling. Le repo regroupe le site public, le back-office admin, les APIs metier, la base MySQL et l'infrastructure Docker locale.

Ce README est la documentation centrale du projet. Les informations qui etaient separees dans les README secondaires ont ete rassemblees ici.

## Stack technique

- Front public : Vue.js 3, Vue Router, Vite
- Back-office admin : Vue.js 3, Vue Router, Vite, Vitest
- API principale : Go, `database/sql`, driver MySQL
- Backoffice historique : PHP 8.2 Apache
- Base de donnees : MySQL 8.4
- Gateway : Nginx
- Outils locaux : Docker Compose, phpMyAdmin

## Structure du repo

```text
admin/
  Interface d'administration Vue.js.

deploy/
  Configuration Docker, Nginx et scripts SQL d'initialisation MySQL.

mission-1-app/
  web-app/
    Site public et espace utilisateur Vue.js.
  api-go/
    API Go connectee a MySQL.
  backoffice/
    Backoffice PHP historique.
```
hpMyAdmin : `http://localhost:8089/`

## Lancement rapide

Prerequis :

- Docker Desktop lance
- Le repo clone localement

Depuis la racine du repo :

```powershell
cd deploy/docker-compose
docker compose up -d --build
```

Puis ouvrir :

```text
http://localhost:8088/
```

Pour verifier les conteneurs :

```powershell
docker compose ps
```

Pour consulter les logs de l'API Go :

```powershell
docker compose logs --tail=120 api-go
```

## Base de donnees

La base MySQL est initialisee automatiquement au premier demarrage du volume Docker.

Scripts SQL :

- `deploy/mysql/init/001_schema.sql` : schema principal fourni pour le projet
- `deploy/mysql/init/002_admin_extensions.sql` : extensions necessaires au back-office admin
- `deploy/mysql/init/003_bootstrap.sql` : donnees minimales de demarrage

Le conteneur MySQL utilise un volume Docker nomme `docker-compose_mysql-data`. Si les scripts SQL changent et qu'il faut reinitialiser completement la base locale :

```powershell
cd deploy/docker-compose
docker compose down
docker volume rm docker-compose_mysql-data
docker compose up -d --build
```

Attention : cette operation supprime les donnees locales de la base Docker.

## API Go

Chemin :

```text
mission-1-app/api-go/
```

L'API Go est connectee a MySQL avec les variables fournies par Docker Compose :

- `DB_HOST=mysql`
- `DB_PORT=3306`
- `DB_USER=root`
- `DB_PASSWORD=password`
- `DB_NAME=upcycletest`

Domaines principaux exposes :

- utilisateurs
- annonces / prestations
- categories
- evenements
- formations
- moderation
- finances
- notifications

Routes utiles via le gateway :

```text
GET /health/api-go
GET /annonces
GET /evenements
GET /formations
GET /categories
GET /api-go/admin/users
GET /api-go/admin/prestations
GET /api-go/admin/categories
GET /api-go/admin/events
GET /api-go/admin/notifications
GET /api-go/admin/finance/overview
```

## Site public

Chemin :

```text
mission-1-app/web-app/
```

Commandes utiles :

```powershell
cd mission-1-app/web-app
npm install
npm run dev
npm run build
```

Le site public est servi par Nginx depuis :

```text
mission-1-app/web-app/dist
```

Pages principales :

- accueil
- annonces
- connexion
- inscription
- profil utilisateur
- tableau de bord utilisateur
- informations profil
- annonces utilisateur
- depots utilisateur

## Back-office admin

Chemin :

```text
admin/
```

Commandes utiles :

```powershell
cd admin
npm install
npm run dev
npm run build
npm run test
```

Le back-office admin est servi par Nginx sous :

```text
http://localhost:8088/admin/
```

Pages principales :

- dashboard
- utilisateurs
- prestations
- categories
- evenements
- moderation
- finances
- notifications

Le back-office admin consomme l'API Go via :

```text
/api-go/admin/*
```

La logique sensible doit rester cote API et base de donnees. Vue sert uniquement d'interface de gestion.

## Backoffice PHP

Chemin :

```text
mission-1-app/backoffice/
```

Le backoffice PHP historique reste disponible derriere le gateway pour certaines routes `/api/*`.

Health check :

```text
http://localhost:8088/health/backoffice
```
