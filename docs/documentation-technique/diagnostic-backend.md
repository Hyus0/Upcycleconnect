# Diagnostic backend UpcycleConnect — état des lieux & plan d'action

Date : 2026-06-24
Auteur : prestation backend (analyse en lecture seule, aucun fichier modifié)

---

## 0. Méthode et limite importante

Cette phase est un **diagnostic en lecture seule**. Aucun fichier, aucune table, aucun volume Docker n'a été modifié.

Limite à connaître : l'environnement d'exécution de l'assistant est un sandbox Linux séparé de ta machine Windows. Il **ne peut pas** atteindre `http://localhost:8088`, ni lancer `docker compose`, ni interroger ton MySQL local. Tout ce qui suit vient de la **lecture du code et des scripts**. Les vérifications « runtime » (docker compose ps, requêtes SQL réelles) devront être lancées par toi ; des commandes prêtes à l'emploi sont fournies en fin de document.

---

## 1. Ce qui a été lu

- `Sujet du projet 2A.pdf` (Mission 1 : solutions applicatives — périmètre fonctionnel des 4 espaces).
- `charte_graphique_upcycleconnect.html` (identité visuelle — hors périmètre backend).
- `README.md`, `PROMPT_REPRISE_IA.md` (stack, lancement, état local).
- `docs/cahier-des-charges/mission-1-mapping.md`, `docs/documentation-technique/mission-1-tech.md`.
- API Go : `api.go` (routeur, CORS), `db/db.go` (connexion), tout `app/` (handlers) et `db/` (accès données), `passwordHashing/`.
- Infra : `deploy/docker-compose/docker-compose.yml`, `deploy/server-config/nginx/default.conf`, `deploy/mysql/init/001..003`.
- Front public Vue : `mission-1-app/web-app/src/**` (pages, services, composables).
- Back-office admin Vue : `admin/src/**` (pages, services `api.js`, `mockDb.js`).

---

## 2. Compréhension du projet

UpcycleConnect = plateforme d'upcycling avec **4 espaces** : Particulier, Prestataire/Artisan, Salarié, Admin (back-office). Fonctionnalités attendues par le sujet : annonces don/vente (validées admin), dépôt d'objets en conteneur avec code-barre puis récupération par un pro, conseils, catalogues (formations/événements/services payants via Stripe), Upcycling Score, planning, projets d'upcycling, forum + modération, messagerie, panier/factures (PDF), abonnements, notifications (OneSignal), multilingue piloté en base.

Architecture réelle du dépôt :

- **Front public** : Vue 3 / Vite → `mission-1-app/web-app`, servi en statique par Nginx (`/`).
- **Back-office admin** : Vue 3 / Vite → `admin`, servi sous `/admin/`.
- **API Go** : `mission-1-app/api-go`, `net/http` standard (routeur `http.ServeMux` Go 1.22+), `database/sql` + driver MySQL. Port 8081.
- **Backoffice PHP** historique : derrière `/api/`.
- **MySQL 8.4**, **Nginx** gateway (port 8088), **phpMyAdmin** (8089), le tout en Docker Compose.

Organisation du code Go (correcte dans le principe) :

- `api.go` : déclaration des routes + CORS + bootstrap.
- `app/` : handlers HTTP (`app.go` 3095 lignes, + `admin_db.go`, `admin_memory.go`).
- `db/` : accès SQL par domaine (annonce, user, panier, messaging, forum, facture…).
- `models/` : structs.
- `passwordHashing/` : bcrypt.

La séparation handlers / db / models existe déjà — la base est saine. Les problèmes sont concentrés sur la **configuration de connexion**, la **cohérence schéma**, la **sécurité/autorisation** et la **cohérence d'appel front↔gateway**.

---

## 3. Cartographie : front → gateway → API Go → tables

Règles de réécriture Nginx (`default.conf`) :

| Préfixe front | Cible Nginx | Effet |
|---|---|---|
| `/go/` | `http://api-go:8081/` | le `/go` est retiré → frappe la racine de l'API Go |
| `/annonces`, `/evenements`, `/formations`, `/categories` | `http://api-go:8081` | proxy direct (préfixe conservé) |
| `/api-go/` | `http://api-go:8081/api/` | `/api-go/admin/x` → `/api/admin/x` côté Go |
| `/api/` | `http://backoffice:80` | **backoffice PHP**, pas l'API Go |
| `/health/api-go`, `/health/backoffice` | health checks | OK |

Domaines fonctionnels et endpoints Go correspondants (extraits de `api.go`) :

| Domaine | Endpoints Go (exemples) | Tables MySQL principales |
|---|---|---|
| Auth / session | `POST /login`, `GET /check-session`, `PUT /users/{id}/password` | UTILISATEUR (colonne `token`) |
| Utilisateurs | `GET/POST /users`, `GET/PUT/DELETE /users/{id}`, `POST /users/{id}/images` | UTILISATEUR |
| Annonces | `GET/POST /annonces`, `GET/PUT/DELETE /annonces/{id}`, favoris, image | ANNONCE, FAVORIS, CATEGORIE |
| Achat / casier / dépôt | `POST /annonces/{id}/acheter`, `/reserver`, `/retirer`, `POST /depot` | ANNONCE, CASIER, CONTENEUR, OBJET, SITE |
| Événements | `GET/POST/PUT/DELETE /evenements`, join/quit/participants | EVENEMENT, EVENEMENT_INSCRIPTION |
| Formations | `GET/POST/PUT/DELETE /formation(s)`, join/quit/participants | FORMATION, FORMATION_INSCRIPTION |
| Catégories | `GET /categories`, CRUD `/category/{id}` | CATEGORIE |
| Projets | CRUD `/projets`, like, join, vues, upload image | PROJET_UPCYCLING/PROJET, PROJET_LIKE, PROJET_VUE, PROJET_INSCRIPTION, ETAPE |
| Conseils (tips) | `GET /tips`, `/tips/role/{role}`, CRUD | TIPS |
| Forum + modération | `GET /forums`, message/topic, signaler, ban, bannis | FORUM, FORUM_MESSAGE, FORUM_SALON, MESSAGE_SIGNALEMENT |
| Panier / factures | `/users/{id}/panier`, `/checkout`, `/factures`, download/send | PANIER_ITEM, COMMANDE, LIGNE_COMMANDE, TRANSACTION, FACTURE |
| Abonnements | `/users/{id}/abonnement` (souscrire/résilier), `/subscription` | ABONNEMENT, ABONNEMENT_UTILISATEUR, TYPE_ABONNEMENT |
| Messagerie (DM) | `/users/{id}/messages…`, offres, ventes, réception, avis | DM_CONVERSATION, DM_MESSAGE, DM_OFFER, DM_SALE, AVIS |
| Notifications | `/users/{id}/notifications`, `/notifications/{id}/read` | NOTIFICATION |
| Stats / éco | `/stats/platform`, `/users/{id}/eco-stats`, `/materiaux/stats` | UPCYCLING_SCORE, OBJET, … |
| Multilingue | `GET /langues`, `/traductions/{code}`, `PUT /users/{id}/langue` | LANGUE, TRADUCTION |
| Admin (back-office) | `/api/admin/users|prestations|categories|events|moderation|finance|notifications` | mappées en base **si** la connexion DB est active (sinon JSON local) |

---

## 4. Manques et incohérences (par criticité)

### 🔴 Bloquants

1. **Connexion DB non configurée pour Docker — c'est LE défaut central.**
   `db/db.go` utilise des constantes **en dur** : `host=localhost`, `password=""` (vide), `dbname=test2_upcycle`. La version qui lit les variables d'environnement (`DB_HOST`, `DB_PASSWORD`, `DB_NAME`…) est **commentée**. Or `docker-compose.yml` fournit `DB_HOST=mysql`, `DB_PASSWORD=password`, `DB_NAME=upcycletest`.
   Conséquence : dans le conteneur, l'API Go tente de joindre `localhost:3306` (= elle-même), mot de passe vide, base `test2_upcycle` inexistante. **Tous les endpoints qui touchent la base échouent (500)** ; seul `/health` (sans DB) répond. Le code marche probablement uniquement en `go run` sur l'hôte avec un MySQL local nommé `test2_upcycle` — ce qui explique le `localhost:8081` codé en dur dans le front (voir point 4).

2. **`sql.Open` ne valide jamais la connexion.** `NewDB()` ne fait pas de `Ping`, ne configure pas le pool (`SetMaxOpenConns`, `SetConnMaxLifetime`), et `panic` sur erreur de parsing seulement. `sql.Open` réussit même si le serveur est injoignable → `db.Conn` est **toujours non-nil**.
   Effet de bord majeur : `adminDBEnabled()` = `db.Conn != nil` renvoie **toujours true**. Le « fallback » mémoire/JSON de l'admin (`admin_memory.go` + `storage/admin_data.json`) est donc **du code mort** quand la DB est en panne → l'admin renvoie des 500 au lieu de basculer. Le toggle DB/mémoire doit reposer sur un vrai `Ping`.

3. **Base de données vide.** D'après `PROMPT_REPRISE_IA.md`, l'import `bdd3` ne contient **aucun INSERT** (`UTILISATEUR = 0`, `ANNONCE = 0`). Même connexion réparée, le front n'aura rien à afficher tant qu'un jeu de données de seed cohérent (et avec **mots de passe bcrypt**) n'est pas chargé.

4. **Front : base d'API incohérente et codée en dur.** **63 fichiers** front utilisent `http://localhost:8081` en dur (`const API_URL = "http://localhost:8081"`), tandis que quelques-uns utilisent `/go` (`SiteNavbar`, `annoncesApi.js`) ou `/api-go` (`userDashboardApi.js`).
   Conséquence : servi derrière le gateway (8088) ou sur le serveur de prod, le navigateur appelle `localhost:8081` qui n'existe pas pour le visiteur → l'app casse hors `go run` local. Le sujet est explicite : « **Un site en localhost ne sera pas corrigé** ». Il faut une **base d'URL unique et relative** (`/go`) passant par le gateway.

### 🟠 Sécurité

5. **Aucune autorisation sur les routes sensibles.** Les endpoints `/api/admin/*` (CRUD users, finance, modération, notifications) et `GET /users` (liste de tous les utilisateurs) n'ont **aucun contrôle d'authentification ni de rôle**. N'importe qui atteignant l'API peut lister/modifier/supprimer.

6. **Pas de middleware d'auth ; vérifs manuelles et partielles.** Le modèle de session = un `token` aléatoire stocké en clair dans `UTILISATEUR.token` (1 seule session par user, pas d'expiration), vérifié au cas par cas dans certains handlers via l'en-tête `Authorization`. Beaucoup de handlers ne vérifient rien. Pas de séparation rôle (Particulier/Prestataire/Salarié/Admin).

7. **CORS figé sur le dev.** `Access-Control-Allow-Origin: http://localhost:5174` en dur. Cassé derrière le gateway et en prod. À aligner avec la stratégie « tout relatif via gateway » (idéalement, plus de CORS cross-origin nécessaire).

8. **Seed avec mots de passe en clair.** `003_bootstrap.sql` insère `password='AdminTemp1!'` (texte brut), alors que `/login` compare en **bcrypt**. Ces comptes seed **ne peuvent pas se connecter**. Les seeds doivent contenir des hash bcrypt.

### 🟡 Cohérence schéma / qualité

9. **Dérive de schéma majeure.** `deploy/mysql/init/001_schema.sql` définit ~21 tables, mais le **code Go référence ~40 tables** absentes des scripts d'init : `CASIER`, `COMMANDE`, `LIGNE_COMMANDE`, `PANIER_ITEM`, `AVIS`, `FAVORIS`, `COMMENTAIRE`, `DM_CONVERSATION`, `DM_MESSAGE`, `DM_OFFER`, `DM_SALE`, `FORUM_SALON`, `MESSAGE_SIGNALEMENT`, `PROJET`/`PROJET_LIKE`/`PROJET_VUE`, `TIPS`, `ABONNEMENT_UTILISATEUR`, la colonne `UTILISATEUR.token`, etc. Le PROMPT confirme : la base réelle a **39 tables**, les scripts d'init sont incomplets.
   Risque direct : **réinitialiser le volume Docker depuis `init/` produit une base qui casse l'app**. Les scripts d'init ne sont **pas** la source de vérité. Il faut régénérer un schéma d'init complet à partir d'un `mysqldump --no-data` de la base réelle.

10. **Doublons/ambiguïtés de modèle** : présence simultanée de `PROJET_UPCYCLING` (init) et `PROJET` (code), `ABONNEMENT` et `ABONNEMENT_UTILISATEUR`. À clarifier (table canonique unique).

11. **Gestion d'erreurs et codes HTTP irréguliers.** Mélange `http.Error`, `fmt.Fprintf`, `json.Encode` ; certains handlers loggent l'erreur SQL et continuent quand même (ex. `GetAllUsers` n'arrête pas après l'erreur). Pas de format d'erreur JSON homogène. Détails techniques potentiellement exposés.

12. **Deux systèmes admin parallèles** (`admin_db.go` vs `admin_memory.go` + `mockDb.js` côté front) — dette à trancher : la base doit être source de vérité, le mock ne sert qu'au dev front isolé.

13. **Dockerfile API** ne copie pas `uploads/` (images annonces/profils servies par `GET /img/`) → images cassées dans le conteneur. À vérifier/monter en volume.

---

## 5. Plan d'action backend priorisé

Principe : **réparer le socle d'abord** (connexion + cohérence schéma), **puis** sécuriser, **puis** compléter/normaliser. Implémentation par lots vérifiables, sans casser les routes consommées par le front.

### Lot 0 — Sauvegarde & vérification (préalable, non destructif)
- Dump horodaté de `upcycletest` dans `temp/` avant toute action DB.
- `docker compose ps`, health checks, `SHOW TABLES`, comptages de lignes : établir l'état réel runtime.
- Geler la base réelle comme **source de vérité** schéma.

### Lot 1 — Réparer la connexion MySQL (🔴 priorité absolue)
- `db/db.go` : activer la lecture des variables d'env (`DB_HOST/PORT/USER/PASSWORD/NAME`) avec valeurs par défaut, **ajouter `Ping` au démarrage**, configurer le pool, échec explicite et journalisé.
- Faire reposer `adminDBEnabled()` sur un vrai état « DB joignable » (Ping), pas sur `Conn != nil`.
- Critère de succès : `/annonces`, `/api/admin/users` renvoient des données réelles dans la stack Docker.

### Lot 2 — Aligner le front sur le gateway (🔴, peu de risque, fort impact)
- Centraliser la base d'API dans **un seul module** (`API_URL = "/go"`, configurable via `VITE_*`), remplacer les 63 occurrences `http://localhost:8081`.
- Vérifier les règles Nginx (`/go/`, `/api-go/`) couvrent tous les préfixes appelés (`/users`, `/projets`, `/tips`, `/notifications`, `/sites`, `/depot`, `/forums`, `/langues`, `/traductions`…). Compléter `default.conf` si besoin.
- Rebuild des `dist` via Docker, test via `http://localhost:8088/`.

### Lot 3 — Cohérence du schéma (🟡→🔴 pour la reproductibilité)
- `mysqldump --no-data` de la base réelle → régénérer `deploy/mysql/init/001_schema.sql` complet (40 tables, colonne `token`, contraintes, index).
- Trancher les doublons (`PROJET` vs `PROJET_UPCYCLING`, `ABONNEMENT` vs `ABONNEMENT_UTILISATEUR`).
- `003_bootstrap.sql` : seeds avec **mots de passe bcrypt**, un compte par rôle, données minimales cohérentes (langue FR, catégories, 1 site + conteneurs).
- Critère : volume neuf reconstruit depuis `init/` → app pleinement fonctionnelle.

### Lot 4 — Sécurité & autorisation (🟠)
- Middleware d'authentification (vérif token) + middleware de rôle, appliqués aux routes sensibles (`/api/admin/*`, mutations, données d'un autre utilisateur).
- Durcir les sessions : expiration du token, ne jamais renvoyer le hash de mot de passe.
- CORS aligné sur la stratégie gateway (relatif) ; en prod, réponses d'erreur non verbeuses.

### Lot 5 — Normalisation handlers/erreurs (🟡)
- Helpers communs `writeJSON` / `writeError` (déjà présents côté admin) généralisés ; codes HTTP cohérents (400/401/403/404/409/422/500).
- Validation stricte des entrées (déjà amorcée côté admin : `validateUser`, etc.).
- Confirmer 100 % de requêtes paramétrées (échantillon analysé : OK, à généraliser).
- Transactions SQL pour les opérations multi-tables : checkout (panier→commande→transaction→facture), achat annonce, dépôt/réservation casier, souscription abonnement.

### Lot 6 — Intégrations & finitions
- Brancher Stripe (paiements) et OneSignal (push) là où c'est `pending`.
- Génération PDF des factures (côté PHP ou Go) + archivage accessible.
- `uploads/` correctement servi/persisté (volume Docker), Dockerfile à corriger.

### Lot 7 — Tests & doc
- Tests Go (`net/http/httptest`) sur les handlers critiques + script de smoke-test des endpoints (curl) via le gateway.
- Doc courte « lancer & tester le backend ».

---

## 6. Risques principaux

- **Schéma réel ≠ scripts d'init** : toute réinit de volume sans régénérer `init/` casse l'app. Ne pas `docker volume rm` sans dump préalable et accord explicite.
- **Base vide** : tests fonctionnels impossibles sans seed cohérent (et bcrypt).
- **Refactor front massif** (63 fichiers) : mécanique mais à faire proprement (module central) pour ne pas casser des routes.
- **Sécurité actuellement ouverte** sur l'admin : à corriger avant tout déploiement réseau réel.
- **Limite d'outillage** : l'assistant ne peut pas exécuter Docker/MySQL de ton poste ; les étapes runtime nécessitent que tu lances les commandes (fournies ci-dessous) et me renvoies les sorties.

---

## 7. Commandes de vérification à lancer (par toi)

```powershell
cd C:\Users\shado\Desktop\Cours\PA\Upcycleconnect
git status --short --branch

cd deploy\docker-compose
docker compose ps

# Schéma réel (source de vérité)
docker compose exec -T mysql mysql -uroot -ppassword -N -e "SHOW TABLES FROM upcycletest;"
docker compose exec -T mysql mysql -uroot -ppassword -N -e "SELECT COUNT(*) FROM upcycletest.UTILISATEUR; SELECT COUNT(*) FROM upcycletest.ANNONCE;"

# Sauvegarde horodatée AVANT toute modif DB
docker compose exec -T mysql mysqldump -uroot -ppassword upcycletest > ..\..\temp\upcycletest-backup-$(Get-Date -Format yyyyMMdd-HHmmss).sql

# Logs API Go (voir l'échec de connexion DB attendu)
docker compose logs --tail=120 api-go

# Endpoints
curl http://localhost:8088/health/api-go
curl http://localhost:8088/go/annonces
curl http://localhost:8088/api-go/admin/users
```

---

## 8. Prochaine étape proposée

Démarrer par le **Lot 0 (sauvegarde + état runtime)** puis le **Lot 1 (connexion DB)**, qui débloque tout le reste. Chaque lot sera implémenté, testé, et tracé. Aucune opération destructive (suppression de volume, écrasement de base, suppression de `temp/`) ne sera faite sans ton accord explicite et une sauvegarde préalable.
