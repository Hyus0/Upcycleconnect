# Mapping CDC -> Mission 1 (MVP courant)

## Technologies imposees

- Front JavaScript: implemente (`mission-1-app/web-app`).
- Back PHP: implemente (`mission-1-app/backoffice`).
- API Go: implemente (`mission-1-app/api-go`).
- Stripe: integration preparee (etat `pending_stripe_integration`).
- OneSignal: integration preparee (etat `onesignal_pending_integration`).
- PDF: non implemente dans ce sprint (a ajouter cote PHP).
- Deploiement serveur reel: stack Docker prete pour deploiement cible.

## Espaces fonctionnels

### Particuliers

- Depot annonces don/vente: OK.
- Depot conteneur validation + code barre: OK.
- Conseils: OK.
- Catalogue formations/services (achat): OK (achat simule, Stripe a connecter).
- Upcycling Score: OK via API Go.
- Planning personnel: OK.
- Tutoriel premiere connexion: OK (modal + localStorage).

### Pro / Artisans

- Abonnements / facturation: MVP OK.
- Acces annonces + achat objets: MVP OK.
- Recuperation objets via conteneurs: flux present via depot/validation.
- Suivi et mise en avant projets: MVP OK.

### Salaries

- Creation/animation formations: MVP OK.
- Gestion planning: extension a finaliser (endpoint deja present pour planning global).
- Gestion conseils/news: MVP news OK.
- Moderation forum: MVP OK (journal des actions).

### Back-Office

- Gestion utilisateurs/acteurs: lecture users OK.
- Gestion financiere: metriques de base OK.
- Validation evenements: endpoint moderation/depot present.
- Gestion catalogue offres: base formations presente.
- Notifications: endpoint admin present (OneSignal en attente).
- Gestion conteneurs/box: depot conteneur present.
