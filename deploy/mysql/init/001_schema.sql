CREATE TABLE LANGUE (
	id INT AUTO_INCREMENT PRIMARY KEY,
	code char(2),
	nom_langue varchar(100)
);

CREATE TABLE TRADUCTION (
	id INT AUTO_INCREMENT PRIMARY KEY,
	id_langue int,
	cle_traduction text,
	text_traduit text,
	FOREIGN KEY (id_langue) REFERENCES LANGUE(id)
);

-- 1. UTILISATEUR
CREATE TABLE UTILISATEUR (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	prenom VARCHAR(50),
    	nom VARCHAR(50),
    	password VARCHAR(255),
    	mail VARCHAR(100) UNIQUE,
    	adresse VARCHAR(100),
    	ville VARCHAR(50),
    	code_postal VARCHAR(10),
    	date_naissance DATE,
        date_inscription DATETIME DEFAULT CURRENT_TIMESTAMP,
    	role ENUM('Particulier', 'Prestataire', 'Admin', 'Salarie'),
	statut ENUM('Actif', 'Inactif', 'Banni', 'En attente de validation') DEFAULT 'Actif',
    	id_langue INT,
	FOREIGN KEY (id_langue) REFERENCES LANGUE(id)
);

-- 2. FORMATION
CREATE TABLE FORMATION (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_formateur INT NOT NULL,
    	type ENUM('Atelier', 'Cours', 'Webinaire'),
    	titre VARCHAR(255),
    	description TEXT,
    	capacite_max INT,
	est_valide ENUM('En attente', 'Valide', 'Refuse') DEFAULT 'En attente',
    	date_debut DATETIME,
    	date_fin DATETIME,
    	statut ENUM('Ouvert', 'Complet', 'Termine', 'Annule'),
    	prix_unitaire DECIMAL(10,2),
    	adresse VARCHAR(100),
    	ville VARCHAR(50),
    	code_postal VARCHAR(5),
    	FOREIGN KEY (id_formateur) REFERENCES UTILISATEUR(id)
);

-- 3. FORMATION_INSCRIPTION
CREATE TABLE FORMATION_INSCRIPTION (
    	id_utilisateur INT NOT NULL,
    	id_formation INT NOT NULL,
    	date_inscription DATETIME DEFAULT CURRENT_TIMESTAMP,
    	PRIMARY KEY (id_utilisateur, id_formation),
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id) ON DELETE CASCADE,
    	FOREIGN KEY (id_formation) REFERENCES FORMATION(id) ON DELETE CASCADE
);

-- 4. FORUM
CREATE TABLE FORUM (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_utilisateur INT NOT NULL,
    	titre VARCHAR(200),
    	sujet TEXT,
    	date_creation DATETIME DEFAULT CURRENT_TIMESTAMP,
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id)
);

-- 5. FORUM_MESSAGE
CREATE TABLE FORUM_MESSAGE (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_utilisateur INT NOT NULL,
    	id_forum INT NOT NULL,
    	contenu TEXT,
    	date_envoi DATETIME DEFAULT CURRENT_TIMESTAMP,
    	est_signale TINYINT(1) DEFAULT 0,
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id),
    	FOREIGN KEY (id_forum) REFERENCES FORUM(id)
);

-- 6. NOTIFICATION
CREATE TABLE NOTIFICATION (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_utilisateur INT NOT NULL,
    	type ENUM('Alerte', 'Message', 'Rappel'),
    	titre VARCHAR(150),
    	message TEXT,
    	lu TINYINT(1) DEFAULT 0,
    	date_envoi DATETIME DEFAULT CURRENT_TIMESTAMP,
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id)
);

-- 7. ANNONCE
CREATE TABLE ANNONCE (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_vendeur INT NOT NULL,
    	id_acheteur INT,
    	titre VARCHAR(100),
    	description TEXT,
    	statut ENUM('Disponible', 'Reserve', 'Vendu', 'Annule'),
    	est_valide ENUM('En attente', 'Valide', 'Refuse') DEFAULT 'En attente',
    	prix DECIMAL(10,2),
    	etat_objet ENUM('Neuf', 'Bon etat', 'Usage'),
    	adresse VARCHAR(100),
    	ville VARCHAR(50),
    	code_postal VARCHAR(5),
    	date_creation DATETIME DEFAULT CURRENT_TIMESTAMP,
    	type ENUM('Don', 'Vente'),
    	FOREIGN KEY (id_vendeur) REFERENCES UTILISATEUR(id),
    	FOREIGN KEY (id_acheteur) REFERENCES UTILISATEUR(id)
);

-- 8. SITE
CREATE TABLE SITE (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	nom VARCHAR(50),
    	ville VARCHAR(50),
    	code_postal VARCHAR(5),
    	adresse VARCHAR(100),
    	telephone VARCHAR(10),
    	type ENUM('Decheterie', 'Point de collecte', 'Association'),
    	actif TINYINT(1) DEFAULT 1
);

-- 9. CONTENEUR
CREATE TABLE CONTENEUR (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_site INT NOT NULL,
    	type_dechet ENUM('Verre', 'Plastique', 'Metal', 'Papier', 'Electronique'),
    	statut ENUM('Operationnel', 'Plein', 'Maintenance'),
    	capacite_max_kg DECIMAL(10,2),
    	niveau_remplissage DECIMAL(10,2) DEFAULT 0,
    	FOREIGN KEY (id_site) REFERENCES SITE(id)
);

-- 10. CODE_ACCES
CREATE TABLE CODE_ACCES (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_conteneur INT NOT NULL,
    	code_barre VARCHAR(255),
    	type_acces VARCHAR(30),
    	est_utilise TINYINT(1) DEFAULT 0,
    	FOREIGN KEY (id_conteneur) REFERENCES CONTENEUR(id)
);

-- 10,5. CATEGORIE
CREATE TABLE CATEGORIE (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    description TEXT,
    id_parent INT DEFAULT NULL,
    statut ENUM('active', 'inactive') DEFAULT 'active',
    date_creation DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_parent) REFERENCES CATEGORIE(id) ON DELETE SET NULL
);

-- 11. OBJET
CREATE TABLE OBJET (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_conteneur INT,
    	id_annonce INT,
    	id_utilisateur INT NOT NULL,
	id_categorie INT NOT NULL,
    	nom VARCHAR(50),
    	description TEXT,
    	type_materiau VARCHAR(50),
	prix_neuf_estime DECIMAL(10,2),
    	poids_estime_kg DECIMAL(10,2),
    	prix DECIMAL(10,2),
    	date_depot DATETIME DEFAULT CURRENT_TIMESTAMP,
    	statut ENUM('En possession', 'En ligne', 'Depose', 'Recycle'),
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id),
    	FOREIGN KEY (id_annonce) REFERENCES ANNONCE(id),
    	FOREIGN KEY (id_conteneur) REFERENCES CONTENEUR(id),
    	FOREIGN KEY (id_categorie) REFERENCES CATEGORIE(id)
);

-- 12. DEPOT_CONTENEUR
CREATE TABLE DEPOT_CONTENEUR (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_utilisateur INT NOT NULL,
    	id_conteneur INT NOT NULL,
    	id_objet INT NOT NULL,
    	date_demande DATETIME DEFAULT CURRENT_TIMESTAMP,
    	date_depot DATETIME,
    	date_recuperation DATETIME,
    	statut ENUM('Prevu', 'Effectue', 'Annule'),
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id),
    	FOREIGN KEY (id_conteneur) REFERENCES CONTENEUR(id),
    	FOREIGN KEY (id_objet) REFERENCES OBJET(id)
);

-- 13. TRANSACTION
CREATE TABLE `TRANSACTION` (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_acheteur INT NOT NULL,
    	id_vendeur INT NOT NULL,
    	id_formation INT,
    	id_annonce INT,
    	montant_ht DECIMAL(10,2),
    	commission_upc DECIMAL(10,2),
    	statut_paiement ENUM('En attente', 'Valide', 'Echoue'),
    	type ENUM('Formation', 'Annonce', 'Abonnement'),
    	date_transaction DATETIME DEFAULT CURRENT_TIMESTAMP,
    	stripe_payment_id VARCHAR(255),
    	FOREIGN KEY (id_acheteur) REFERENCES UTILISATEUR(id),
    	FOREIGN KEY (id_vendeur) REFERENCES UTILISATEUR(id),
    	FOREIGN KEY (id_formation) REFERENCES FORMATION(id),
    	FOREIGN KEY (id_annonce) REFERENCES ANNONCE(id)
);

-- 14. FACTURE
CREATE TABLE FACTURE (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_transaction INT NOT NULL,
    	numero_facture VARCHAR(50) UNIQUE,
	type_paiement ENUM('Carte', 'Espece', 'Virement') DEFAULT 'Carte',
	FOREIGN KEY (id_transaction) REFERENCES `TRANSACTION`(id)
);

-- 15. PROJET_UPCYCLING
CREATE TABLE PROJET_UPCYCLING (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_createur INT NOT NULL,
    titre VARCHAR(100) NOT NULL,
    description_courte VARCHAR(255),
    date_creation DATETIME DEFAULT CURRENT_TIMESTAMP,
    score_impact DECIMAL(10,2),
    nb_vues INT DEFAULT 0,
    nb_likes INT DEFAULT 0,
    co2_evite_kg DECIMAL(10,2),
    visible_public TINYINT(1) DEFAULT 1,
    FOREIGN KEY (id_createur) REFERENCES UTILISATEUR(id)
);
-- 16. PROJET_INSCRIPTION
CREATE TABLE PROJET_INSCRIPTION (
    	id_utilisateur INT NOT NULL,
    	id_projet INT NOT NULL,
    	date_inscription DATETIME DEFAULT CURRENT_TIMESTAMP,
    	PRIMARY KEY (id_utilisateur, id_projet),
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id) ON DELETE CASCADE,
    	FOREIGN KEY (id_projet) REFERENCES PROJET_UPCYCLING(id) ON DELETE CASCADE
);

CREATE TABLE ETAPE (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_projet INT NOT NULL,
	numero_ordre INT NOT NULL,
	titre VARCHAR(100),
    	description TEXT,
    	image_url VARCHAR(255),
    	FOREIGN KEY (id_projet) REFERENCES PROJET_UPCYCLING(id)
);

-- 17. EVENEMENT
CREATE TABLE EVENEMENT (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	titre VARCHAR(100),
    	description TEXT,
    	adresse VARCHAR(100),
    	ville VARCHAR(50),
    	code_postal VARCHAR(5),
    	date_creation DATETIME DEFAULT CURRENT_TIMESTAMP,
    	date_evenement DATETIME,
    	type ENUM('Atelier', 'Collecte', 'Conference', 'Echange')
);

-- 18. EVENEMENT_INSCRIPTION
CREATE TABLE EVENEMENT_INSCRIPTION (
    	id_utilisateur INT NOT NULL,
    	id_evenement INT NOT NULL,
    	date_inscription DATETIME DEFAULT CURRENT_TIMESTAMP,
    	PRIMARY KEY (id_utilisateur, id_evenement),
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id) ON DELETE CASCADE,
    	FOREIGN KEY (id_evenement) REFERENCES EVENEMENT(id) ON DELETE CASCADE
);

-- 19. UPCYCLING_SCORE
CREATE TABLE UPCYCLING_SCORE (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_utilisateur INT NOT NULL,
    	ressources_economisees DECIMAL(10,2),
    	co2_total_evite_kg DECIMAL(10,2),
    	nb_objets_recycles INT DEFAULT 0,
    	total_points INT DEFAULT 0,
    	FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id)
);

-- 20. TYPE ABONNEMENT
CREATE TABLE TYPE_ABONNEMENT (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	nom VARCHAR(50) NOT NULL,
    	description TEXT,
    	prix_ht DECIMAL(10, 2) NOT NULL,
    	duree_mois INT NOT NULL
);

-- 21. ABONNEMENT
CREATE TABLE ABONNEMENT (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	id_acheteur INT NOT NULL,
    	id_type_abonnement INT NOT NULL,
    	date_debut DATETIME DEFAULT CURRENT_TIMESTAMP,
    	date_fin DATETIME,
    	statut ENUM('Actif', 'Expire', 'Resilie') DEFAULT 'Actif',
    	stripe_subscription_id VARCHAR(255),
    	FOREIGN KEY (id_acheteur) REFERENCES UTILISATEUR(id),
    	FOREIGN KEY (id_type_abonnement) REFERENCES TYPE_ABONNEMENT(id)
);
