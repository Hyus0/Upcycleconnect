# Prompt de reprise IA - UpcycleConnect

Tu reprends le projet UpcycleConnect en local sur la machine Windows.

## Contexte rapide

- Projet stocke ici : `C:\Users\shado\Desktop\Cours\PA\Upcycleconnect`
- Depot GitHub : `https://github.com/Hyus0/Upcycleconnect.git`
- Branche active : `main`
- Dernier commit recupere depuis GitHub : `1cc988d` - `Finalisation de la moderation des forum`
- Date de preparation locale : `2026-06-24`

Le dossier local a ete mis a jour depuis GitHub avec `git pull --ff-only` sur `main`. Une autre branche locale existe encore : `test/admin-backoffice-api-20260611-165012`, avec un commit specifique admin/back-office (`3765b37`) qui n'est pas dans `main`. Ne la supprime pas sans demande explicite.

## Stack

- Front public : Vue 3, Vue Router, Vite
- Back-office admin : Vue 3, Vue Router, Vite, Vitest
- API principale : Go, `database/sql`, driver MySQL
- Backoffice historique : PHP 8.2 Apache
- Base locale : MySQL 8.4 via Docker Compose
- Gateway : Nginx
- Outil DB : phpMyAdmin

## Etat local verifie

Les services Docker sont demarres depuis `deploy/docker-compose` :

- `gateway` : `http://localhost:8088/`
- `phpmyadmin` : `http://localhost:8089/`
- `api-go`
- `backoffice`
- `mysql`

Verifications faites :

- `http://localhost:8088/` repond en `200`
- `http://localhost:8088/health/api-go` repond `"en vie."`
- `http://localhost:8088/health/backoffice` repond `{"status":"ok","service":"backoffice-php"}`
- Base MySQL active : `upcycletest`
- Identifiants Docker DB : user `root`, password `password`
- La base `upcycletest` contient 39 tables.
- Nouvelle base importee le `2026-06-24` depuis `C:\Users\shado\Downloads\bdd(3).txt`, via la copie corrigee `temp/bdd3-20260624.sql`.
- Le fichier fourni ne contient aucun `INSERT` : la base importee est un schema vide en donnees (`UTILISATEUR = 0`, `ANNONCE = 0` apres import).

Note importante : un volume Docker MySQL existait deja (`docker-compose_mysql-data`). L'ancienne base a ete sauvegardee avant remplacement. Les dumps locaux non suivis par Git sont dans `temp/`, notamment :

- `temp/bdd2-fixed.sql`
- `temp/seed-bdd2-data.sql`
- `temp/upcycletest-backup-before-bdd2-20260611.sql`
- `temp/bdd3-20260624.sql`
- `temp/upcycletest-backup-before-bdd3-20260624-094718.sql`

Si une nouvelle piece jointe SQL doit remplacer la base, demander confirmation avant de supprimer le volume ou d'ecraser `upcycletest`.

## Commandes utiles

Depuis PowerShell :

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect
git status
```

Demarrer ou relancer la stack :

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect\deploy\docker-compose
docker compose up -d --build
docker compose ps
```

Arreter la stack sans supprimer la base :

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect\deploy\docker-compose
docker compose down
```

Reinitialiser completement la base Docker uniquement si l'utilisateur confirme la perte des donnees locales :

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect\deploy\docker-compose
docker compose down
docker volume rm docker-compose_mysql-data
docker compose up -d --build
```

Verifier la base :

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect\deploy\docker-compose
docker compose exec -T mysql mysql -uroot -ppassword -N -e "SHOW TABLES FROM upcycletest;"
```

Rebuilder les fronts sans Node local, via Docker :

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect\mission-1-app\web-app
docker run --rm -v "${PWD}:/app" -v /app/node_modules -w /app node:20-alpine sh -lc "npm ci && npm run build"

cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect\admin
docker run --rm -v "${PWD}:/app" -v /app/node_modules -w /app node:20-alpine sh -lc "npm ci && npm run build"
```

Note : `node` et `npm` ne sont pas dans le PATH Windows global sur cette machine. Les builds front ont ete faits avec Docker.

## Points d'entree fonctionnels

- Site public : `http://localhost:8088/`
- Back-office admin : `http://localhost:8088/admin/`
- phpMyAdmin : `http://localhost:8089/`
- API Go health : `http://localhost:8088/health/api-go`
- Backoffice PHP health : `http://localhost:8088/health/backoffice`

Routes API utiles :

- `GET /annonces`
- `GET /evenements`
- `GET /formations`
- `GET /categories`
- `GET /api-go/admin/users`
- `GET /api-go/admin/prestations`
- `GET /api-go/admin/categories`
- `GET /api-go/admin/events`
- `GET /api-go/admin/notifications`
- `GET /api-go/admin/finance/overview`

## Structure a connaitre

- `README.md` : documentation centrale du projet
- `deploy/docker-compose/docker-compose.yml` : stack locale
- `deploy/mysql/init/` : scripts SQL d'initialisation MySQL pour volume neuf
- `mission-1-app/web-app/` : front public Vue
- `mission-1-app/api-go/` : API Go
- `mission-1-app/backoffice/` : backoffice PHP historique
- `admin/` : back-office admin Vue
- `docs/` : documentation projet
- `temp/` : dumps SQL locaux non suivis par Git

## Mission pour l'IA qui reprend

1. Commencer par lire `README.md`, puis verifier `git status --short --branch`.
2. Ne pas supprimer `temp/` ni le volume Docker MySQL sans confirmation.
3. Si une tache touche la base, inspecter d'abord le schema reel dans `upcycletest`, car il est plus riche que les scripts `deploy/mysql/init`.
4. Si une tache touche le front public ou admin, rebuilder les `dist` avec Docker avant de tester via `http://localhost:8088/`.
5. Tester les changements avec les endpoints health et, si possible, avec les pages concernees dans le navigateur.
6. Garder les changements scopes : eviter les refontes non demandees.

## Utilisation avec Claude Cowork

Dans Claude Cowork, ouvrir ou declarer comme dossier de travail :

```text
C:\Users\shado\Desktop\Cours\PA\Upcycleconnect
```

Puis donner a Claude Cowork ce fichier comme prompt de demarrage, ou coller son contenu dans le premier message. Demander ensuite explicitement a Claude de :

- travailler dans le dossier ci-dessus ;
- commencer par `git status --short --branch` ;
- utiliser Docker Compose depuis `deploy/docker-compose` ;
- ne pas supprimer le volume `docker-compose_mysql-data` sans accord ;
- tester via `http://localhost:8088/`, `http://localhost:8088/admin/` et `http://localhost:8089/`.
