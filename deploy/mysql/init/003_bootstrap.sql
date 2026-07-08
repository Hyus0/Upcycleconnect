-- 003_bootstrap.sql
-- Donnees de demarrage (idempotent). Mots de passe en bcrypt.
-- Comptes: admin@/Admin123!, particulier@/Test123!, prestataire@/Test123!, salarie@/Test123!

-- seed-complete.sql - Jeu de donnees de demonstration UpcycleConnect
-- Idempotent (re-executable) et transactionnel. INSERT uniquement, aucune suppression.
-- Mots de passe stockes en bcrypt (cout 14). Comptes de test :
--   admin@upcycleconnect.local        / Admin123!
--   particulier@upcycleconnect.local  / Test123!
--   prestataire@upcycleconnect.local  / Test123!
--   salarie@upcycleconnect.local      / Test123!

SET NAMES utf8mb4;
START TRANSACTION;

-- ---------- LANGUES ----------
INSERT INTO LANGUE (code, nom_langue)
SELECT 'fr', 'Francais' WHERE NOT EXISTS (SELECT 1 FROM LANGUE WHERE code = 'fr');
INSERT INTO LANGUE (code, nom_langue)
SELECT 'en', 'English' WHERE NOT EXISTS (SELECT 1 FROM LANGUE WHERE code = 'en');
SET @lang_fr = (SELECT id FROM LANGUE WHERE code = 'fr' LIMIT 1);

-- ---------- UTILISATEURS (1 par role) ----------
INSERT INTO UTILISATEUR (prenom, nom, password, mail, mail_valide, ville, code_postal, is_admin, role, statut, id_langue)
SELECT 'Nina', 'Roux', '$2b$14$GjnYMjTrtwQ84RIEUOapOOVgOsVW4mulSI29RsuLwG4JtK0/rXYp.', 'admin@upcycleconnect.local', 1, 'Paris', '75011', 1, 'Admin', 'Actif', @lang_fr
WHERE NOT EXISTS (SELECT 1 FROM UTILISATEUR WHERE mail = 'admin@upcycleconnect.local');

INSERT INTO UTILISATEUR (prenom, nom, password, mail, mail_valide, ville, code_postal, is_admin, role, statut, id_langue, materiaux_recherches)
SELECT 'Marie', 'Lambert', '$2b$14$.VOPrsLkikXAvW4is99CtOMtWYx7XGZurKBeAD/gW/QOBrH2ZqMH6', 'particulier@upcycleconnect.local', 1, 'Paris', '75011', 0, 'Particulier', 'Actif', @lang_fr, 'bois,palette,textile'
WHERE NOT EXISTS (SELECT 1 FROM UTILISATEUR WHERE mail = 'particulier@upcycleconnect.local');

INSERT INTO UTILISATEUR (prenom, nom, password, mail, mail_valide, siret, siret_valide, ville, code_postal, is_admin, role, statut, id_langue)
SELECT 'Karim', 'Benali', '$2b$14$HoblX5En65JUvK8z4NFMQOOg/KGXCLk6Ae0JAwE2lqtGuyJ2MXXxy', 'prestataire@upcycleconnect.local', 1, '12345678901234', 1, 'Lyon', '69002', 0, 'Prestataire', 'Actif', @lang_fr
WHERE NOT EXISTS (SELECT 1 FROM UTILISATEUR WHERE mail = 'prestataire@upcycleconnect.local');

INSERT INTO UTILISATEUR (prenom, nom, password, mail, mail_valide, ville, code_postal, is_admin, role, statut, id_langue)
SELECT 'Lucas', 'Moreau', '$2b$14$Q/vm6v31MlA5xgmC95J/Teq3krRp71cgRjIWwxMi3iJYnuCdLAQJO', 'salarie@upcycleconnect.local', 1, 'Lille', '59000', 0, 'Salarie', 'Actif', @lang_fr
WHERE NOT EXISTS (SELECT 1 FROM UTILISATEUR WHERE mail = 'salarie@upcycleconnect.local');

SET @u_admin  = (SELECT id FROM UTILISATEUR WHERE mail = 'admin@upcycleconnect.local' LIMIT 1);
SET @u_part   = (SELECT id FROM UTILISATEUR WHERE mail = 'particulier@upcycleconnect.local' LIMIT 1);
SET @u_presta = (SELECT id FROM UTILISATEUR WHERE mail = 'prestataire@upcycleconnect.local' LIMIT 1);
SET @u_sal    = (SELECT id FROM UTILISATEUR WHERE mail = 'salarie@upcycleconnect.local' LIMIT 1);

-- ---------- CATEGORIES ----------
INSERT INTO CATEGORIE (nom, description, statut)
SELECT 'Mobilier', 'Meubles et accessoires', 'active' WHERE NOT EXISTS (SELECT 1 FROM CATEGORIE WHERE nom = 'Mobilier');
INSERT INTO CATEGORIE (nom, description, statut)
SELECT 'Textile', 'Tissus et vetements a revaloriser', 'active' WHERE NOT EXISTS (SELECT 1 FROM CATEGORIE WHERE nom = 'Textile');
INSERT INTO CATEGORIE (nom, description, statut)
SELECT 'Bois', 'Chutes et palettes', 'active' WHERE NOT EXISTS (SELECT 1 FROM CATEGORIE WHERE nom = 'Bois');
INSERT INTO CATEGORIE (nom, description, statut)
SELECT 'Electronique', 'Composants et appareils', 'active' WHERE NOT EXISTS (SELECT 1 FROM CATEGORIE WHERE nom = 'Electronique');
INSERT INTO CATEGORIE (nom, description, statut)
SELECT 'Metal', 'Pieces metalliques', 'active' WHERE NOT EXISTS (SELECT 1 FROM CATEGORIE WHERE nom = 'Metal');

SET @cat_mobilier = (SELECT id FROM CATEGORIE WHERE nom = 'Mobilier' LIMIT 1);
SET @cat_textile  = (SELECT id FROM CATEGORIE WHERE nom = 'Textile' LIMIT 1);
SET @cat_bois     = (SELECT id FROM CATEGORIE WHERE nom = 'Bois' LIMIT 1);

-- ---------- SITE / CONTENEUR / CASIER ----------
INSERT INTO SITE (nom, ville, code_postal, adresse, telephone, type, actif)
SELECT 'Point Collecte Paris 11', 'Paris', '75011', '174 rue La Fayette', '0145000000', 'Point de collecte', 1
WHERE NOT EXISTS (SELECT 1 FROM SITE WHERE nom = 'Point Collecte Paris 11');
SET @site1 = (SELECT id FROM SITE WHERE nom = 'Point Collecte Paris 11' LIMIT 1);

-- type_dechet : enum('Verre','Plastique','Metal','Papier','Electronique')
INSERT INTO CONTENEUR (id_site, type_dechet, statut, capacite_max_kg, niveau_remplissage)
SELECT @site1, 'Papier', 'Operationnel', 500.00, 120.00
WHERE NOT EXISTS (SELECT 1 FROM CONTENEUR WHERE id_site = @site1 AND type_dechet = 'Papier');
INSERT INTO CONTENEUR (id_site, type_dechet, statut, capacite_max_kg, niveau_remplissage)
SELECT @site1, 'Metal', 'Operationnel', 800.00, 300.00
WHERE NOT EXISTS (SELECT 1 FROM CONTENEUR WHERE id_site = @site1 AND type_dechet = 'Metal');
SET @cont1 = (SELECT id FROM CONTENEUR WHERE id_site = @site1 AND type_dechet = 'Papier' LIMIT 1);

INSERT INTO CASIER (id_conteneur, numero_casier, taille, statut)
SELECT @cont1, 'A01', 'Moyen', 'Libre' WHERE NOT EXISTS (SELECT 1 FROM CASIER WHERE id_conteneur = @cont1 AND numero_casier = 'A01');
INSERT INTO CASIER (id_conteneur, numero_casier, taille, statut)
SELECT @cont1, 'A02', 'Grand', 'Libre' WHERE NOT EXISTS (SELECT 1 FROM CASIER WHERE id_conteneur = @cont1 AND numero_casier = 'A02');

-- ---------- ANNONCES (validees, disponibles) ----------
INSERT INTO ANNONCE (id_vendeur, id_categorie, id_site, titre, description, type_materiau, poids_estime_kg, prix, etat_objet, statut, est_valide, type, ville, code_postal, adresse)
SELECT @u_part, @cat_bois, @site1, 'Lot de palettes bois', 'Cinq palettes en bon etat pour mobilier DIY.', 'Bois', 40.00, 0.00, 'Bon etat', 'Disponible', 'Valide', 'Don', 'Paris', '75011', '174 rue La Fayette'
WHERE NOT EXISTS (SELECT 1 FROM ANNONCE WHERE titre = 'Lot de palettes bois' AND id_vendeur = @u_part);

INSERT INTO ANNONCE (id_vendeur, id_categorie, id_site, titre, description, type_materiau, poids_estime_kg, prix, etat_objet, statut, est_valide, type, ville, code_postal, adresse)
SELECT @u_part, @cat_mobilier, @site1, 'Chaise vintage a restaurer', 'Chaise bois annees 60, structure solide.', 'Bois', 6.00, 25.00, 'Usage', 'Disponible', 'Valide', 'Vente', 'Paris', '75011', '174 rue La Fayette'
WHERE NOT EXISTS (SELECT 1 FROM ANNONCE WHERE titre = 'Chaise vintage a restaurer' AND id_vendeur = @u_part);

INSERT INTO ANNONCE (id_vendeur, id_categorie, id_site, titre, description, type_materiau, poids_estime_kg, prix, etat_objet, statut, est_valide, type, ville, code_postal, adresse)
SELECT @u_part, @cat_textile, @site1, 'Coupons de tissu coton', 'Chutes de coton multicolores, ideal couture.', 'Textile', 3.00, 10.00, 'Neuf', 'Disponible', 'Valide', 'Vente', 'Paris', '75011', '174 rue La Fayette'
WHERE NOT EXISTS (SELECT 1 FROM ANNONCE WHERE titre = 'Coupons de tissu coton' AND id_vendeur = @u_part);

-- ---------- FORMATION (par le salarie, validee) ----------
INSERT INTO FORMATION (id_formateur, type, titre, description, capacite_max, est_valide, date_debut, date_fin, statut, prix_unitaire, adresse, ville, code_postal)
SELECT @u_sal, 'Atelier', 'Initiation a la restauration de meubles', 'Apprenez les bases de la remise en etat du mobilier bois.', 12, 'Valide', '2026-07-15 10:00:00', '2026-07-15 17:00:00', 'Ouvert', 49.00, '174 rue La Fayette', 'Paris', '75011'
WHERE NOT EXISTS (SELECT 1 FROM FORMATION WHERE titre = 'Initiation a la restauration de meubles');

-- ---------- EVENEMENT (par le salarie) ----------
INSERT INTO EVENEMENT (id_createur, titre, description, adresse, ville, code_postal, date_evenement, type)
SELECT @u_sal, 'Collecte textile de quartier', 'Journee de collecte et tri de textiles a revaloriser.', '174 rue La Fayette', 'Paris', '75011', '2026-07-20 09:00:00', 'Collecte'
WHERE NOT EXISTS (SELECT 1 FROM EVENEMENT WHERE titre = 'Collecte textile de quartier');

-- ---------- TIPS / CONSEILS ----------
INSERT INTO TIPS (id_createur, titre, description, role_cible, actif)
SELECT @u_sal, 'Bien preparer son bois avant poncage', 'Nettoyez, degraissez puis poncez du grain le plus gros au plus fin.', 'Particulier', 1
WHERE NOT EXISTS (SELECT 1 FROM TIPS WHERE titre = 'Bien preparer son bois avant poncage');
INSERT INTO TIPS (id_createur, titre, description, role_cible, actif)
SELECT @u_sal, 'Optimiser la recuperation en conteneur', 'Verifiez la categorie et le poids avant depot pour accelerer la validation.', 'Prestataire', 1
WHERE NOT EXISTS (SELECT 1 FROM TIPS WHERE titre = 'Optimiser la recuperation en conteneur');

-- ---------- TYPES D ABONNEMENT ----------
INSERT INTO TYPE_ABONNEMENT (nom, description, prix_ht, duree_mois)
SELECT 'Freemium', 'Acces de base gratuit', 0.00, 1 WHERE NOT EXISTS (SELECT 1 FROM TYPE_ABONNEMENT WHERE nom = 'Freemium');
INSERT INTO TYPE_ABONNEMENT (nom, description, prix_ht, duree_mois)
SELECT 'Premium', 'Mise en avant des projets et acces complet', 19.90, 1 WHERE NOT EXISTS (SELECT 1 FROM TYPE_ABONNEMENT WHERE nom = 'Premium');

-- ---------- PROJET UPCYCLING (par le prestataire) ----------
INSERT INTO PROJET_UPCYCLING (id_createur, image_url, titre, description_courte, score_impact, nb_vues, nb_likes, co2_evite_kg, visible_public)
SELECT @u_presta, '/img/projets/placeholder.png', 'Table basse en palettes', 'Transformation de palettes en table basse design.', 85.00, 12, 3, 18.50, 1
WHERE NOT EXISTS (SELECT 1 FROM PROJET_UPCYCLING WHERE titre = 'Table basse en palettes');

-- ---------- UPCYCLING SCORE (particulier) ----------
INSERT INTO UPCYCLING_SCORE (id_utilisateur, ressources_economisees, co2_total_evite_kg, nb_objets_recycles, total_points, niveau)
SELECT @u_part, 45.00, 30.00, 4, 120, 'Bronze'
WHERE NOT EXISTS (SELECT 1 FROM UPCYCLING_SCORE WHERE id_utilisateur = @u_part);


-- ---------- FORUM (salon + sujet) ----------
INSERT INTO FORUM_SALON (nom, description)
SELECT 'General', 'Discussions generales de la communaute' WHERE NOT EXISTS (SELECT 1 FROM FORUM_SALON WHERE nom = 'General');
SET @salon_gen = (SELECT id FROM FORUM_SALON WHERE nom = 'General' LIMIT 1);

INSERT INTO FORUM (id_utilisateur, id_salon, titre, sujet, ouvert)
SELECT @u_part, @salon_gen, 'Idees pour recycler des palettes', 'Partagez vos meilleures idees de projets a base de palettes.', 1
WHERE NOT EXISTS (SELECT 1 FROM FORUM WHERE titre = 'Idees pour recycler des palettes');

INSERT INTO LANGUE (code, nom_langue) VALUES 
('fr', 'Français'),
('en', 'English'),
('es', 'Espagnol');


-- ===========================
-- SITES
-- ===========================

INSERT INTO SITE (nom, ville, code_postal, adresse, telephone, type) VALUES
('EcoCentre Lille', 'Lille', '59000', '15 Rue des Recycleurs', '0320123456', 'Decheterie'),
('Point Vert Roubaix', 'Roubaix', '59100', '8 Avenue de l''Industrie', '0320654321', 'Point de collecte'),
('Association RecyNord', 'Tourcoing', '59200', '42 Rue des Freres Lumiere', '0320789456', 'Association'),
('EcoPark Villeneuve', 'Villeneuve-d''Ascq', '59650', '120 Boulevard de l''Innovation', '0320456789', 'Decheterie');

-- ===========================
-- CONTENEURS
-- ===========================

INSERT INTO CONTENEUR (id_site, type_dechet, statut, capacite_max_kg, niveau_remplissage) VALUES

-- Site 1
(1, 'Electronique', 'Operationnel', 800, 220),
(1, 'Metal', 'Operationnel', 1000, 540),

-- Site 2
(2, 'Plastique', 'Operationnel', 700, 180),
(2, 'Papier', 'Plein', 900, 900),

-- Site 3
(3, 'Verre', 'Operationnel', 1200, 450),
(3, 'Electronique', 'Maintenance', 600, 0),

-- Site 4
(4, 'Metal', 'Operationnel', 1000, 610),
(4, 'Plastique', 'Operationnel', 750, 340);

-- ===========================
-- CASIERS
-- 4 casiers par conteneur
-- ===========================

INSERT INTO CASIER (id_conteneur, numero_casier, taille, statut) VALUES

-- Conteneur 1
(1,'A1','Petit','Libre'),
(1,'A2','Petit','Occupe'),
(1,'A3','Moyen','Libre'),
(1,'A4','Grand','Reserve'),

-- Conteneur 2
(2,'B1','Petit','Libre'),
(2,'B2','Moyen','Occupe'),
(2,'B3','Grand','Libre'),
(2,'B4','Grand','Maintenance'),

-- Conteneur 3
(3,'C1','Petit','Libre'),
(3,'C2','Petit','Libre'),
(3,'C3','Moyen','Reserve'),
(3,'C4','Grand','Occupe'),

-- Conteneur 4
(4,'D1','Petit','Occupe'),
(4,'D2','Petit','Occupe'),
(4,'D3','Moyen','Maintenance'),
(4,'D4','Grand','Occupe'),

-- Conteneur 5
(5,'E1','Petit','Libre'),
(5,'E2','Petit','Libre'),
(5,'E3','Moyen','Occupe'),
(5,'E4','Grand','Reserve'),

-- Conteneur 6
(6,'F1','Petit','Maintenance'),
(6,'F2','Petit','Maintenance'),
(6,'F3','Moyen','Maintenance'),
(6,'F4','Grand','Maintenance'),

-- Conteneur 7
(7,'G1','Petit','Libre'),
(7,'G2','Petit','Occupe'),
(7,'G3','Moyen','Libre'),
(7,'G4','Grand','Reserve'),

-- Conteneur 8
(8,'H1','Petit','Libre'),
(8,'H2','Petit','Occupe'),
(8,'H3','Moyen','Libre'),
(8,'H4','Grand','Libre');

INSERT IGNORE INTO TYPE_ABONNEMENT (id, nom, description, prix_ht, duree_mois, color) VALUES 
(1, 'DM Plus', 'Messagerie illimitee entre membres UpcycleConnect', 2.99, 1, '2d7a4f'),
(2, 'Artisan Débutant', 'Plan gratuit', 0.00, 1, 'ffffff'),
(3, 'Pro & Entreprises', 'Plan complet', 25.00, 1, '2d7a4f');

INSERT INTO ABONNEMENT_AVANTAGE (id_type_abonnement, nom, disponible) VALUES
(2, 'Vitrine limitée à 5 projets', TRUE),
(2, 'Publication & recherche d''annonces', TRUE),
(2, 'Réservation standard de matériaux', TRUE),
(2, 'Statistiques & Impact écologique', FALSE),
(2, 'Alertes de collecte priorisées', FALSE),
(2, 'Campagnes de Sponsoring', FALSE);

INSERT INTO ABONNEMENT_AVANTAGE (id_type_abonnement, nom, disponible) VALUES
(3, 'Vitrine étendue à 20 projets', TRUE),
(3, 'Tout le plan gratuit', TRUE),
(3, 'Tableaux de bord avancés', TRUE),
(3, 'Analyse d''impact écologique détaillée', TRUE),
(3, 'Alertes priorisées pour la collecte', TRUE),
(3, 'Création de campagnes publicitaires', TRUE);

INSERT INTO FORUM_SALON (nom, description) VALUES
('Général', 'Discussions générales de la communauté'),
('Upcycling', 'Partagez vos créations et techniques'),
('Événements', 'Organisez et rejoignez des événements locaux'),
('Entraide', 'Posez vos questions, obtenez de l aide');

INSERT INTO TYPE_ABONNEMENT (nom, description, prix_ht, duree_mois)
SELECT 'DM Plus', 'Messagerie illimitee entre membres UpcycleConnect', 2.99, 1
WHERE NOT EXISTS (SELECT 1 FROM TYPE_ABONNEMENT WHERE nom = 'DM Plus');

INSERT INTO TYPE_ABONNEMENT (nom, description, prix_ht, duree_mois) 
VALUES (
    'Premium Pro', 
    'Accès aux tableaux de bord avancés, statistiques de matériaux et alertes priorisées.', 
    25.00, -- Prix entre 15 et 30 euros
    1      -- Renouvellement mensuel
);

INSERT INTO COMMENTAIRE (id_utilisateur, description) 
VALUES 
(2, "Grâce à UpcycleConnect, j'ai trouvé les matériaux parfaits pour mes créations. La plateforme a transformé mon activité d'artisan."),

(2, "Une interface simple et efficace. Les bornes sont très bien pensées pour récupérer des objets lourds sans encombrer les autres."),

(3, "Enfin un outil sérieux pour l'économie circulaire. Les échanges avec les autres membres sont toujours très professionnels."),

(3, "Le système de réservation de casier est brillant. Cela libère une énorme quantité d'espace chez soi tout en aidant des artisans locaux."),

(4, "Service au top ! J'ai pu donner une seconde vie à mes vieux meubles en toute simplicité. Je recommande vivement à tous ceux qui veulent réduire leurs déchets.");

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Dashboard', 'Tableau de bord'),
(2, 'Dashboard', 'Dashboard'),
(3, 'Dashboard', 'Panel');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Hello', 'Bonjour'),
(2, 'Hello', 'Hello'),
(3, 'Hello', 'Buen día');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'SummaryActivity', 'Voici un résumé de votre activité sur UpcycleConnect'),
(2, 'SummaryActivity', 'Here is a summary of your activity on UpcycleConnect
'),
(3, 'SummaryActivity', 'Aquí tienes un resumen de tu actividad en UpcycleConnect.');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Co2', 'CO2 évité'),
(2, 'Co2', 'CO2 avoided'),
(3, 'Co2', 'CO2 evitado');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Objets', 'Objets'),
(2, 'Objets', 'Objects'),
(3, 'Objets', 'Objetos');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Economise', 'Économisé'),
(2, 'Economise', 'Saved'),
(3, 'Economise', 'Guardado');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Co2', 'CO2 évité'),
(2, 'Co2', 'CO2 avoided'),
(3, 'Co2', 'CO2 evitado');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Co2', 'CO2 évité'),
(2, 'Co2', 'CO2 avoided'),
(3, 'Co2', 'CO2 evitado');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'activeAds', 'Annonces actives'),
(2, 'activeAds', 'Active announcements'),
(3, 'activeAds', 'Anuncios activos');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Month', 'ce mois'),
(2, 'Month', 'this month'),
(3, 'Month', 'este mes');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'NoActiveAds', 'Aucune active'),
(2, 'NoActiveAds', 'None active'),
(3, 'NoActiveAds', 'Ninguno activo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'MyDeposits', 'Mes dépôts'),
(2, 'MyDeposits', 'My deposits'),
(3, 'MyDeposits', 'Mis depósitos');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'BeDeposited', 'À déposer en consigne'),
(2, 'BeDeposited', 'To be deposited'),
(3, 'BeDeposited', 'Para ser depositado');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'NothingDeposits', 'Rien à déposer'),
(2, 'NothingDeposits', 'Nothing to deposit'),
(3, 'NothingDeposits', 'nada que depositar');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'LastAds', 'Mes dernières annonces'),
(2, 'LastAds', 'My latest announcements'),
(3, 'LastAds', 'Mis últimos anuncios');

--

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Home', 'ACCUEIL'),
(2, 'Home', 'HOME'),
(3, 'Home', 'HOGAR');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Objet', 'OBJET'),
(2, 'Objet', 'OBJECT'),
(3, 'Objet', 'OBJETO');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Type', 'TYPE'),
(2, 'Type', 'TYPE'),
(3, 'Type', 'AMABLE');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'StatutPublication', 'STATUT PUBLICATION'),
(2, 'StatutPublication', 'PUBLICATION STATUS'),
(3, 'StatutPublication', 'ESTADO DE PUBLICACIÓN');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'LogisticStatut', 'STATUT LOGISTIQUE'),
(2, 'LogisticStatut', 'LOGISTICS STATUS'),
(3, 'LogisticStatut', 'ESTADO LOGÍSTICO');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Date', 'DATE'),
(2, 'Date', 'DATE'),
(3, 'Date', 'FECHA');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Actions', 'ACTIONS'),
(2, 'Actions', 'ACTIONS'),
(3, 'Actions', 'ACCIONES ');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Voir', 'Voir'),
(2, 'Voir', 'See'),
(3, 'Voir', 'Ver');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Modifier', 'Modifier'),
(2, 'Modifier', 'To modify'),
(3, 'Modifier', 'para modificar');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Supprimer', 'Supprimer'),
(2, 'Supprimer', 'Delete'),
(3, 'Supprimer', 'Borrar');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'NewAnnonce', '+ Nouvelle annonce'),
(2, 'NewAnnonce', '+ New announcement'),
(3, 'NewAnnonce', '+ Nuevo anuncio');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'DeposeAnnonce', '+ Déposer une annonce'),
(2, 'DeposeAnnonce', '+ Post an announcement'),
(3, 'DeposeAnnonce', '+ Publicar un anuncio');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'MySchedule', 'Mon planning'),
(2, 'MySchedule', 'My schedule'),
(3, 'MySchedule', 'Mi horario');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Schedule', 'Cliquez pour voir le calendrier complet.'),
(2, 'Schedule', 'Click to see the full schedule'),
(3, 'Schedule', 'Haz clic para ver el programa completo.');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Lundi', 'LUN'),
(2, 'Lundi', 'MON'),
(3, 'Lundi', 'LUN');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Mardi', 'MAR'),
(2, 'Mardi', 'TUE'),
(3, 'Mardi', 'MAR');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Mercredi', 'MER'),
(2, 'Mercredi', 'WED'),
(3, 'Mercredi', 'MIE');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Jeudi', 'JEU'),
(2, 'Jeudi', 'THU'),
(3, 'Jeudi', 'JUE');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Vendredi', 'Vendredi'),
(2, 'Vendredi', 'FRI'),
(3, 'Vendredi', 'VIE');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Samedi', 'SAM'),
(2, 'Samedi', 'SAT'),
(3, 'Samedi', 'SAB');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Dimanche', 'DIM'),
(2, 'Dimanche', 'SUN'),
(3, 'Dimanche', 'DOM');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'VueMensuelle', 'Vue mensuelle'),
(2, 'VueMensuelle', 'Monthly view'),
(3, 'VueMensuelle', 'Vista mensual');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'PlanningInfo', 'Le planning se remplira dès que vous rejoindrez une formation ou un événement.'),
(2, 'PlanningInfo', 'The schedule will fill up as soon as you join a training course or event.'),
(3, 'PlanningInfo', 'El calendario se llenará en cuanto te inscribas en un curso de formación o evento.');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Tips', 'Conseil du jour'),
(2, 'Tips', 'Tip of the day'),
(3, 'Tips', 'Consejo del día');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Notification', 'Notification • Nouvelle'),
(2, 'Notification', 'Notification • New'),
(3, 'Notification', 'Notificación • Nuevo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'ReadNext', 'Lire la suite'),
(2, 'ReadNext', 'Read more'),
(3, 'ReadNext', 'Leer más');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'SeeAction', "Voir l'action"),
(2, 'SeeAction', 'See the action'),
(3, 'SeeAction', 'Ver la acción');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Nouvelle', 'Nouvelle'),
(2, 'Nouvelle', 'New'),
(3, 'Nouvelle', 'Nuevo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'NothingNotification', 'Aucune notification'),
(2, 'NothingNotification', 'No notification'),
(3, 'NothingNotification', 'Sin notificacion');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'NotificationText', "Vous n'avez aucun message ou alerte non lus dans votre espace personnel pour le moment."),
(2, 'NotificationText', 'You currently have no unread messages or alerts in your personal space.'),
(3, 'NotificationText', 'Actualmente no tienes mensajes ni alertas sin leer en tu espacio personal.');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'WaitingTips', "En attente d'astuces..."),
(2, 'WaitingTips', 'Looking for tips...'),
(3, 'WaitingTips', 'Buscando consejos...');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'TipsText', "Restez à l'affût, de nouveaux conseils pour votre profil arrivent bientôt !"),
(2, 'TipsText', 'Stay tuned, new tips for your profile are coming soon!'),
(3, 'TipsText', '¡Estad atentos, pronto habrá nuevos consejos para vuestro perfil!');

--

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Niveau', 'Niveau'),
(2, 'Niveau', 'Level'),
(3, 'Niveau', 'Nivel');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Calendrier', "Calendrier d'inscriptions"),
(2, 'Calendrier', 'Registration Calendar'),
(3, 'Calendrier', 'Calendario de Inscripciones');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Precedent', 'Précédent'),
(2, 'Precedent', 'Previous'),
(3, 'Precedent', 'Anterior');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Suivant', 'Suivant'),
(2, 'Suivant', 'Next'),
(3, 'Suivant', 'próximo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Today', "Aujourd'hui"),
(2, 'Today', 'Today'),
(3, 'Today', 'Hoy');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Formation', 'Formation'),
(2, 'Formation', 'Formation'),
(3, 'Formation', 'Formación');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Evenement', 'Événement'),
(2, 'Evenement', 'Event'),
(3, 'Evenement', 'Evento');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'BackToDashboard', 'Retour au tableau de bord'),
(2, 'BackToDashboard', 'Return to the dashboard'),
(3, 'BackToDashboard', 'Regresar al tablero');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Janvier', 'Janvier'),
(2, 'Janvier', 'January'),
(3, 'Janvier', 'Enero');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Fevrier', 'Février'),
(2, 'Fevrier', 'February'),
(3, 'Fevrier', 'Febrero');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Mars', 'Mars'),
(2, 'Mars', 'March'),
(3, 'Mars', 'Marzo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Avril', 'Avril'),
(2, 'Avril', 'April'),
(3, 'Avril', 'Abril');

-- MAI
INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Mai', 'Mai'),
(2, 'Mai', 'May'),
(3, 'Mai', 'Mayo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Juin', 'Juin'),
(2, 'Juin', 'June'),
(3, 'Juin', 'Junio');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Juillet', 'Juillet'),
(2, 'Juillet', 'July'),
(3, 'Juillet', 'Julio');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Aout', 'Août'),
(2, 'Aout', 'August'),
(3, 'Aout', 'Agosto');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Septembre', 'Septembre'),
(2, 'Septembre', 'September'),
(3, 'Septembre', 'Septiembre');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Octobre', 'Octobre'),
(2, 'Octobre', 'October'),
(3, 'Octobre', 'Octubre');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Novembre', 'Novembre'),
(2, 'Novembre', 'November'),
(3, 'Novembre', 'Noviembre');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Decembre', 'Décembre'),
(2, 'Decembre', 'December'),
(3, 'Decembre', 'Diciembre');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Suivant', 'Suivant'),
(2, 'Suivant', 'Next'),
(3, 'Suivant', 'próximo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Suivant', 'Suivant'),
(2, 'Suivant', 'Next'),
(3, 'Suivant', 'próximo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'Suivant', 'Suivant'),
(2, 'Suivant', 'Next'),
(3, 'Suivant', 'próximo');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'ActivePlatform', 'Plateforme active'),
(1, 'Members', 'membres'),
(1, 'HeroSubtitle', 'UpcycleConnect réunit particuliers, artisans et entreprises autour de l''upcycling. Déposez, échangez, transformez — et mesurez votre impact environnemental.'),
(1, 'StartFree', 'Commencer gratuitement'),
(1, 'ImaPro', 'Je suis professionnel'),
(1, 'CO2Avoided', 'CO₂ évité au total'),
(1, 'UpcycledObjects', 'Objets upcyclés'),
(1, 'PartnerArtisans', 'Artisans partenaires'),
(1, 'SitesIDF', 'Sites à Paris & IDF'),
(1, 'FeaturesSub', 'FONCTIONNALITÉS'),
(1, 'AllYouNeed', 'Tout ce dont vous avez besoin'),
(1, 'FeaturesDesc', 'Une plateforme complète pour donner une seconde vie à vos objets, du dépôt jusqu''à la création finale.'),
(1, 'SmartAdsTitle', 'Annonces intelligentes'),
(1, 'SmartAdsDesc', 'Publiez vos dons ou ventes avec photos, catégorie et localisation. Les artisans reçoivent des alertes personnalisées selon leurs besoins.'),
(1, 'SecureContainersTitle', 'Conteneurs sécurisés'),
(1, 'SecureContainersDesc', 'Demandez un code-barres pour déposer votre objet dans l''un de nos 30+ conteneurs à Paris et en IDF. Simple, rapide, sécurisé.'),
(1, 'UpcyclingScoreTitle', 'Upcycling Score'),
(1, 'UpcyclingScoreDesc', 'Chaque action compte. Visualisez en temps réel le CO₂ évité, les ressources économisées et votre contribution à l''économie circulaire.'),
(1, 'TrainingWorkshopsTitle', 'Formations & Ateliers'),
(1, 'TrainingWorkshopsDesc', 'Apprenez les techniques d''upcycling avec nos formateurs experts. Cours en ligne et ateliers en présentiel dans nos 7 sites.'),
(1, 'PersonalScheduleTitle', 'Planning personnel'),
(1, 'PersonalScheduleDesc', 'Gérez vos inscriptions, formations, dépôts et événements depuis un calendrier centralisé et synchronisable.'),
(1, 'LivingCommunityTitle', 'Communauté vivante'),
(1, 'LivingCommunityDesc', 'Forums, événements, conseils et projets partagés. Rejoignez une communauté de passionnés qui agissent pour l''environnement.'),
(1, 'ProcessSub', 'PROCESSUS'),
(1, 'HowItWorks', 'Comment ça marche ?'),
(1, 'Step1Title', 'Déposez une annonce'),
(1, 'Step1Desc', 'Photographiez et décrivez l''objet à donner ou vendre sur la plateforme.'),
(1, 'Step2Title', 'Validation'),
(1, 'Step2Desc', 'Notre équipe valide votre objet et vous envoie un code-barres d''accès.'),
(1, 'Step3Title', 'Dépôt en conteneur'),
(1, 'Step3Desc', 'Déposez l''objet dans le conteneur le plus proche grâce à votre code-barres.'),
(1, 'Step4Title', 'Seconde vie !'),
(1, 'Step4Desc', 'Un artisan récupère l''objet et lui donne une nouvelle vie.'),
(1, 'ReadyAct', 'Prêt à agir pour'),
(1, 'ThePlanet', 'la planète ?'),
(1, 'ReadyDesc1', 'Rejoignez les'),
(1, 'ReadyDesc2', 'membres qui donnent une seconde vie à leurs objets chaque jour.'),
(1, 'CreateFreeAccount', 'Créer mon compte gratuit'),
(1, 'ViewAds', 'Voir les annonces');


INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(2, 'ActivePlatform', 'Active platform'),
(2, 'Members', 'members'),
(2, 'HeroSubtitle', 'UpcycleConnect brings together individuals, artisans, and companies around upcycling. Drop off, exchange, transform — and measure your environmental impact.'),
(2, 'StartFree', 'Start for free'),
(2, 'ImaPro', 'I am a professional'),
(2, 'CO2Avoided', 'Total CO₂ avoided'),
(2, 'UpcycledObjects', 'Upcycled objects'),
(2, 'PartnerArtisans', 'Partner artisans'),
(2, 'SitesIDF', 'Locations in Paris & IDF'),
(2, 'FeaturesSub', 'FEATURES'),
(2, 'AllYouNeed', 'Everything you need'),
(2, 'FeaturesDesc', 'A complete platform to give a second life to your items, from drop-off to the final creation.'),
(2, 'SmartAdsTitle', 'Smart ads'),
(2, 'SmartAdsDesc', 'Publish your donations or sales with photos, category, and location. Artisans receive personalized alerts based on their needs.'),
(2, 'SecureContainersTitle', 'Secure containers'),
(2, 'SecureContainersDesc', 'Request a barcode to drop off your item in one of our 30+ containers in Paris and IDF. Simple, fast, secure.'),
(2, 'UpcyclingScoreTitle', 'Upcycling Score'),
(2, 'UpcyclingScoreDesc', 'Every action counts. View in real-time the CO₂ avoided, resources saved, and your contribution to the circular economy.'),
(2, 'TrainingWorkshopsTitle', 'Training & Workshops'),
(2, 'TrainingWorkshopsDesc', 'Learn upcycling techniques with our expert trainers. Online courses and in-person workshops at our 7 locations.'),
(2, 'PersonalScheduleTitle', 'Personal schedule'),
(2, 'PersonalScheduleDesc', 'Manage your registrations, training, drop-offs, and events from a centralized, synchronizable calendar.'),
(2, 'LivingCommunityTitle', 'Vibrant community'),
(2, 'LivingCommunityDesc', 'Forums, events, tips, and shared projects. Join a community of enthusiasts acting for the environment.'),
(2, 'ProcessSub', 'PROCESS'),
(2, 'HowItWorks', 'How does it work?'),
(2, 'Step1Title', 'Post an ad'),
(2, 'Step1Desc', 'Take a picture and describe the item to give away or sell on the platform.'),
(2, 'Step2Title', 'Validation'),
(2, 'Step2Desc', 'Our team validates your item and sends you an access barcode.'),
(2, 'Step3Title', 'Container drop-off'),
(2, 'Step3Desc', 'Drop off the item in the nearest container using your barcode.'),
(2, 'Step4Title', 'Second life!'),
(2, 'Step4Desc', 'An artisan retrieves the item and gives it a new life.'),
(2, 'ReadyAct', 'Ready to act for'),
(2, 'ThePlanet', 'the planet?'),
(2, 'ReadyDesc1', 'Join the'),
(2, 'ReadyDesc2', 'members who give a second life to their items every day.'),
(2, 'CreateFreeAccount', 'Create my free account'),
(2, 'ViewAds', 'View ads');


INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(3, 'ActivePlatform', 'Plataforma activa'),
(3, 'Members', 'miembros'),
(3, 'HeroSubtitle', 'UpcycleConnect reúne a particulares, artesanos y empresas en torno al upcycling. Deposita, intercambia, transforma y mide tu impacto ambiental.'),
(3, 'StartFree', 'Comenzar gratis'),
(3, 'ImaPro', 'Soy profesional'),
(3, 'CO2Avoided', 'CO₂ total evitado'),
(3, 'UpcycledObjects', 'Objetos reciclados'),
(3, 'PartnerArtisans', 'Artesanos asociados'),
(3, 'SitesIDF', 'Sitios en París e IDF'),
(3, 'FeaturesSub', 'CARACTERÍSTICAS'),
(3, 'AllYouNeed', 'Todo lo que necesitas'),
(3, 'FeaturesDesc', 'Una plataforma completa para dar una segunda vida a tus objetos, desde el depósito hasta la creación final.'),
(3, 'SmartAdsTitle', 'Anuncios inteligentes'),
(3, 'SmartAdsDesc', 'Publica tus donaciones o ventas con fotos, categoría y ubicación. Los artesanos reciben alertas personalizadas según sus necesidades.'),
(3, 'SecureContainersTitle', 'Contenedores seguros'),
(3, 'SecureContainersDesc', 'Solicita un código de barras para depositar tu objeto en uno de nuestros más de 30 contenedores en París e IDF. Simple, rápido, seguro.'),
(3, 'UpcyclingScoreTitle', 'Upcycling Score'),
(3, 'UpcyclingScoreDesc', 'Cada acción cuenta. Visualiza en tiempo real el CO₂ evitado, los recursos ahorrados y tu contribución a la economía circular.'),
(3, 'TrainingWorkshopsTitle', 'Formación y Talleres'),
(3, 'TrainingWorkshopsDesc', 'Aprende técnicas de upcycling con nuestros formadores expertos. Cursos en línea y talleres presenciales en nuestras 7 ubicaciones.'),
(3, 'PersonalScheduleTitle', 'Agenda personal'),
(3, 'PersonalScheduleDesc', 'Gestiona tus inscripciones, formaciones, depósitos y eventos desde un calendario centralizado y sincronizable.'),
(3, 'LivingCommunityTitle', 'Comunidad vibrante'),
(3, 'LivingCommunityDesc', 'Foros, eventos, consejos y proyectos compartidos. Únete a una comunidad de entusiastas que actúan por el medio ambiente.'),
(3, 'ProcessSub', 'PROCESO'),
(3, 'HowItWorks', '¿Cómo funciona?'),
(3, 'Step1Title', 'Publicar un anuncio'),
(3, 'Step1Desc', 'Fotografía y describe el objeto para regalar o vender en la plataforma.'),
(3, 'Step2Title', 'Validación'),
(3, 'Step2Desc', 'Nuestro equipo valida tu objeto y te envía un código de barras de acceso.'),
(3, 'Step3Title', 'Depósito en contenedor'),
(3, 'Step3Desc', 'Deposita el objeto en el contenedor más cercano utilizando tu código de barras.'),
(3, 'Step4Title', '¡Segunda vida!'),
(3, 'Step4Desc', 'Un artesano recupera el objeto y le da una nueva vida.'),
(3, 'ReadyAct', '¿Listo para actuar por'),
(3, 'ThePlanet', 'el planeta?'),
(3, 'ReadyDesc1', 'Únete a los'),
(3, 'ReadyDesc2', 'miembros que dan una segunda vida a sus objetos cada día.'),
(3, 'CreateFreeAccount', 'Crear mi cuenta gratis'),
(3, 'ViewAds', 'Ver los anuncios');

-- Footer

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'FooterDesc', 'La plateforme de référence pour l''upcycling en France et en Europe. Rejoignez une communauté engagée pour l''économie circulaire.'),
(1, 'FooterCol1', 'PLATEFORME'),
(1, 'FooterHowItWorks', 'Comment ça marche'),
(1, 'FooterAds', 'Annonces'),
(1, 'FooterDashboard', 'Tableau de bord'),
(1, 'FooterMyAds', 'Mes annonces'),
(1, 'FooterCol2', 'PROFESSIONNELS'),
(1, 'FooterSubOffers', 'Offres d''abonnement'),
(1, 'FooterContainers', 'Conteneurs'),
(1, 'FooterProjects', 'Projets upcycling'),
(1, 'FooterAdvertising', 'Publicité'),
(1, 'FooterCol3', 'ENTREPRISE'),
(1, 'FooterAbout', 'À propos'),
(1, 'FooterLocations', 'Nos sites'),
(1, 'FooterBlog', 'Blog'),
(1, 'FooterContact', 'Contact'),
(1, 'FooterLegal', 'Mentions légales'),
(1, 'FooterTOS', 'CGU'),
(1, 'FooterGDPR', 'RGPD');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(2, 'FooterDesc', 'The reference platform for upcycling in France and Europe. Join a community committed to the circular economy.'),
(2, 'FooterCol1', 'PLATFORM'),
(2, 'FooterHowItWorks', 'How it works'),
(2, 'FooterAds', 'Ads'),
(2, 'FooterDashboard', 'Dashboard'),
(2, 'FooterMyAds', 'My ads'),
(2, 'FooterCol2', 'PROFESSIONALS'),
(2, 'FooterSubOffers', 'Subscription offers'),
(2, 'FooterContainers', 'Containers'),
(2, 'FooterProjects', 'Upcycling projects'),
(2, 'FooterAdvertising', 'Advertising'),
(2, 'FooterCol3', 'COMPANY'),
(2, 'FooterAbout', 'About us'),
(2, 'FooterLocations', 'Our locations'),
(2, 'FooterBlog', 'Blog'),
(2, 'FooterContact', 'Contact'),
(2, 'FooterLegal', 'Legal Notice'),
(2, 'FooterTOS', 'Terms of Service'),
(2, 'FooterGDPR', 'GDPR');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(3, 'FooterDesc', 'La plataforma de referencia para el upcycling en Francia y Europa. Únete a una comunidad comprometida con la economía circular.'),
(3, 'FooterCol1', 'PLATAFORMA'),
(3, 'FooterHowItWorks', 'Cómo funciona'),
(3, 'FooterAds', 'Anuncios'),
(3, 'FooterDashboard', 'Tablero'),
(3, 'FooterMyAds', 'Mis anuncios'),
(3, 'FooterCol2', 'PROFESIONALES'),
(3, 'FooterSubOffers', 'Ofertas de suscripción'),
(3, 'FooterContainers', 'Contenedores'),
(3, 'FooterProjects', 'Proyectos de upcycling'),
(3, 'FooterAdvertising', 'Publicidad'),
(3, 'FooterCol3', 'EMPRESA'),
(3, 'FooterAbout', 'Acerca de'),
(3, 'FooterLocations', 'Nuestros sitios'),
(3, 'FooterBlog', 'Blog'),
(3, 'FooterContact', 'Contacto'),
(3, 'FooterLegal', 'Aviso legal'),
(3, 'FooterTOS', 'Términos de uso'),
(3, 'FooterGDPR', 'RGPD');

-- LoginPage
INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'WelcomeBackLeft', 'Heureux de vous '),
(1, 'WelcomeBackAccent', 'revoir'),
(1, 'LoginDesc', 'Connectez-vous pour retrouver vos annonces, vos projets en cours et continuer à faire grandir la communauté de l''upcycling.'),
(1, 'StatCO2Month', 'CO₂ évité / mois'),
(1, 'StatUpcycled', 'Objets upcyclés'),
(1, 'StatArtisans', 'Artisans actifs'),
(1, 'MemberUpcycle', 'Membre UpcycleConnect'),
(1, 'LoadingReviews', 'Chargement des avis...'),
(1, 'WelcomeBackRight', 'Bon retour parmi nous'),
(1, 'NotMemberYet', 'Pas encore membre ?'),
(1, 'CreateAccountFree', 'Créer un compte gratuitement'),
(1, 'EmailLabel', 'Adresse e-mail'),
(1, 'PasswordLabel', 'Mot de passe'),
(1, 'ForgotPwd', 'Oublié ?'),
(1, 'RememberMe', 'Se souvenir de moi'),
(1, 'LoginBtn', 'Se connecter'),
(1, 'OrContinueWith', 'ou continuer avec');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(2, 'WelcomeBackLeft', 'Happy to see you '),
(2, 'WelcomeBackAccent', 'again'),
(2, 'LoginDesc', 'Log in to find your ads, your current projects, and continue to grow the upcycling community.'),
(2, 'StatCO2Month', 'CO₂ avoided / month'),
(2, 'StatUpcycled', 'Upcycled objects'),
(2, 'StatArtisans', 'Active artisans'),
(2, 'MemberUpcycle', 'UpcycleConnect Member'),
(2, 'LoadingReviews', 'Loading reviews...'),
(2, 'WelcomeBackRight', 'Welcome back'),
(2, 'NotMemberYet', 'Not a member yet?'),
(2, 'CreateAccountFree', 'Create an account for free'),
(2, 'EmailLabel', 'Email address'),
(2, 'PasswordLabel', 'Password'),
(2, 'ForgotPwd', 'Forgot it?'),
(2, 'RememberMe', 'Remember me'),
(2, 'LoginBtn', 'Log in'),
(2, 'OrContinueWith', 'or continue with');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(3, 'WelcomeBackLeft', 'Feliz de verte '),
(3, 'WelcomeBackAccent', 'de nuevo'),
(3, 'LoginDesc', 'Inicia sesión para encontrar tus anuncios, tus proyectos actuales y continuar haciendo crecer la comunidad del upcycling.'),
(3, 'StatCO2Month', 'CO₂ evitado / mes'),
(3, 'StatUpcycled', 'Objetos reciclados'),
(3, 'StatArtisans', 'Artesanos activos'),
(3, 'MemberUpcycle', 'Miembro de UpcycleConnect'),
(3, 'LoadingReviews', 'Cargando reseñas...'),
(3, 'WelcomeBackRight', 'Bienvenido de nuevo'),
(3, 'NotMemberYet', '¿Aún no eres miembro?'),
(3, 'CreateAccountFree', 'Crea una cuenta gratis'),
(3, 'EmailLabel', 'Correo electrónico'),
(3, 'PasswordLabel', 'Contraseña'),
(3, 'ForgotPwd', '¿Olvidada?'),
(3, 'RememberMe', 'Recuérdame'),
(3, 'LoginBtn', 'Iniciar sesión'),
(3, 'OrContinueWith', 'o continuar con');

-- RegisterPage
INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(1, 'JoinThe', 'Rejoignez la '),
(1, 'Community', 'communauté'),
(1, 'ThatActs', 'qui agit.'),
(1, 'RegisterDesc', 'Plus de 12 400 membres, 340 artisans et des <br />centaines de projets d''upcycling chaque mois. <br />Commencez gratuitement dès aujourd''hui.'),
(1, 'TestimonialText', '"Grâce à UpcycleConnect, j''ai trouvé les matériaux parfaits pour mes créations. La plateforme a transformé mon activité d''artisan."'),
(1, 'TestimonialAuthorRole', 'Artisane ébéniste — Paris 11e'),
(1, 'CreateAccountTitle', 'Créer un compte'),
(1, 'AlreadyMember', 'Déjà membre ?'),
(1, 'LoginLink', 'Se connecter'),
(1, 'IAm', 'Je suis...'),
(1, 'Individual', 'Particulier'),
(1, 'ProArtisan', 'Pro / Artisan'),
(1, 'FirstName', 'Prénom'),
(1, 'LastName', 'Nom'),
(1, 'ZipCode', 'Code postal'),
(1, 'IAcceptThe', 'J''accepte les'),
(1, 'TOS', 'CGU'),
(1, 'AndThe', 'et la'),
(1, 'PrivacyPolicy', 'politique de confidentialité'),
(1, 'CreateMyAccountFreeBtn', 'Créer mon compte gratuitement'),
(1, 'Or', 'ou'),
(1, 'ContinueWithGoogle', 'Continuer avec Google');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(2, 'JoinThe', 'Join the '),
(2, 'Community', 'community'),
(2, 'ThatActs', 'that takes action.'),
(2, 'RegisterDesc', 'Over 12,400 members, 340 artisans, and <br />hundreds of upcycling projects every month. <br />Start for free today.'),
(2, 'TestimonialText', '"Thanks to UpcycleConnect, I found the perfect materials for my creations. The platform transformed my artisan business."'),
(2, 'TestimonialAuthorRole', 'Cabinetmaker — Paris 11th'),
(2, 'CreateAccountTitle', 'Create an account'),
(2, 'AlreadyMember', 'Already a member?'),
(2, 'LoginLink', 'Log in'),
(2, 'IAm', 'I am...'),
(2, 'Individual', 'Individual'),
(2, 'ProArtisan', 'Pro / Artisan'),
(2, 'FirstName', 'First name'),
(2, 'LastName', 'Last name'),
(2, 'ZipCode', 'Zip code'),
(2, 'IAcceptThe', 'I accept the'),
(2, 'TOS', 'TOS'),
(2, 'AndThe', 'and the'),
(2, 'PrivacyPolicy', 'privacy policy'),
(2, 'CreateMyAccountFreeBtn', 'Create my free account'),
(2, 'Or', 'or'),
(2, 'ContinueWithGoogle', 'Continue with Google');

INSERT INTO TRADUCTION (id_langue, cle_traduction, text_traduit) VALUES
(3, 'JoinThe', 'Únete a la '),
(3, 'Community', 'comunidad'),
(3, 'ThatActs', 'que actúa.'),
(3, 'RegisterDesc', 'Más de 12.400 miembros, 340 artesanos y <br />cientos de proyectos de upcycling cada mes. <br />Empieza gratis hoy mismo.'),
(3, 'TestimonialText', '"Gracias a UpcycleConnect, encontré los materiales perfectos para mis creaciones. La plataforma transformó mi negocio artesanal."'),
(3, 'TestimonialAuthorRole', 'Ebanista — París 11'),
(3, 'CreateAccountTitle', 'Crear una cuenta'),
(3, 'AlreadyMember', '¿Ya eres miembro?'),
(3, 'LoginLink', 'Iniciar sesión'),
(3, 'IAm', 'Soy...'),
(3, 'Individual', 'Particular'),
(3, 'ProArtisan', 'Pro / Artesano'),
(3, 'FirstName', 'Nombre'),
(3, 'LastName', 'Apellido'),
(3, 'ZipCode', 'Código postal'),
(3, 'IAcceptThe', 'Acepto los'),
(3, 'TOS', 'Términos'),
(3, 'AndThe', 'y la'),
(3, 'PrivacyPolicy', 'política de privacidad'),
(3, 'CreateMyAccountFreeBtn', 'Crear mi cuenta gratis'),
(3, 'Or', 'o'),
(3, 'ContinueWithGoogle', 'Continuar con Google');

COMMIT;

