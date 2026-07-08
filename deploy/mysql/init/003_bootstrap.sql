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

COMMIT;

