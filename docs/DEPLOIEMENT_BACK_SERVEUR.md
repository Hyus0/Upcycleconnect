# Deploiement du back sur serveur

Ce guide deploie uniquement le back d'UpcycleConnect. Le front public peut rester dans son emplacement actuel sur le serveur.

## Infrastructure actuelle

Le serveur physique expose Proxmox. Dans Proxmox :

- VM `100` : EVE-NG, ne pas utiliser pour UpcycleConnect.
- VM `101` : web, sert deja le front sur `http://95.216.9.170:8088/`.

Le deploiement ci-dessous vise donc la VM `101`, ou une nouvelle VM/CT dediee au back si la RAM de la VM `101` est trop limitee.

La VM `101` visible sur Proxmox a seulement 2 Go de RAM et une utilisation memoire elevee. Pour une demo rapide, le back peut tourner dessus. Pour une pratique plus propre, creer une VM/CT dediee au back avec 4 Go de RAM minimum, puis faire pointer le Nginx de la VM `101` vers cette VM.

## Emplacement conseille

```bash
/opt/upcycleconnect
```

## Recuperer le projet

```bash
cd /opt
git clone https://github.com/Hyus0/Upcycleconnect.git upcycleconnect
cd /opt/upcycleconnect
git fetch origin
git checkout test/back
git pull origin test/back
```

Si le dossier existe deja :

```bash
cd /opt/upcycleconnect
git fetch origin
git checkout test/back
git pull origin test/back
```

## Variables serveur

Creer le fichier `/opt/upcycleconnect/deploy/docker-compose/.env` :

```bash
MYSQL_ROOT_PASSWORD=remplacer_par_un_mot_de_passe_fort
MYSQL_DATABASE=upcycletest
```

## Lancer les services back

```bash
cd /opt/upcycleconnect/deploy/docker-compose
docker compose -f docker-compose.server-back.yml up -d --build
```

Services exposes localement sur le serveur :

- API Go : `127.0.0.1:8081`
- Backoffice PHP : `127.0.0.1:8090`
- phpMyAdmin : `127.0.0.1:8089`
- MySQL : accessible uniquement par le reseau Docker

## Brancher le Nginx existant

Ajouter ces blocs dans le `server { ... }` qui sert deja le front public.

```nginx
location /go/ {
  proxy_pass http://127.0.0.1:8081/;
  proxy_http_version 1.1;
  proxy_set_header Host $host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
}

location /api-go/ {
  proxy_pass http://127.0.0.1:8081/api/;
  proxy_http_version 1.1;
  proxy_set_header Host $host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
}

location /api/ {
  proxy_pass http://127.0.0.1:8090;
  proxy_http_version 1.1;
  proxy_set_header Host $host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
}

location /health/api-go {
  proxy_pass http://127.0.0.1:8081/health;
}

location /health/backoffice {
  proxy_pass http://127.0.0.1:8090/health;
}
```

Si l'admin Vue doit aussi etre disponible sur le serveur :

```bash
cd /opt/upcycleconnect/admin
npm ci
npm run build
```

Puis ajouter au meme `server { ... }` Nginx :

```nginx
location /admin/ {
  alias /opt/upcycleconnect/admin/dist/;
  try_files $uri $uri/ /admin/index.html;
}
```

## Verifier

```bash
curl http://127.0.0.1:8081/health
curl http://127.0.0.1:8090/health
docker compose -f /opt/upcycleconnect/deploy/docker-compose/docker-compose.server-back.yml ps
```

Depuis le navigateur :

- `http://IP_DU_SERVEUR/health/api-go`
- `http://IP_DU_SERVEUR/health/backoffice`
- `http://IP_DU_SERVEUR/admin/` si l'admin est active

## Mise a jour

```bash
cd /opt/upcycleconnect
git pull origin test/back
cd deploy/docker-compose
docker compose -f docker-compose.server-back.yml up -d --build
```
