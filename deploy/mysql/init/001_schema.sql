-- 001_schema.sql
-- Schema complet UpcycleConnect, genere depuis la base reelle (upcycletest).
-- Genere le 2026-06-24. Source de verite du schema (39 tables).
-- Remplace l'ancien schema partiel. NE PAS editer a la main : regenerer via mysqldump --no-data.


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
DROP TABLE IF EXISTS `ABONNEMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ABONNEMENT` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_acheteur` int NOT NULL,
  `id_type_abonnement` int NOT NULL,
  `date_debut` datetime DEFAULT CURRENT_TIMESTAMP,
  `date_fin` datetime DEFAULT NULL,
  `statut` enum('Actif','Expire','Resilie') COLLATE utf8mb4_unicode_ci DEFAULT 'Actif',
  `stripe_subscription_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_acheteur` (`id_acheteur`),
  KEY `id_type_abonnement` (`id_type_abonnement`),
  CONSTRAINT `ABONNEMENT_ibfk_1` FOREIGN KEY (`id_acheteur`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `ABONNEMENT_ibfk_2` FOREIGN KEY (`id_type_abonnement`) REFERENCES `TYPE_ABONNEMENT` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ABONNEMENT_UTILISATEUR`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ABONNEMENT_UTILISATEUR` (
  `id_abonne` int NOT NULL,
  `id_suivi` int NOT NULL,
  `date_abonnement` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_abonne`,`id_suivi`),
  KEY `id_suivi` (`id_suivi`),
  CONSTRAINT `ABONNEMENT_UTILISATEUR_ibfk_1` FOREIGN KEY (`id_abonne`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `ABONNEMENT_UTILISATEUR_ibfk_2` FOREIGN KEY (`id_suivi`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ANNONCE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ANNONCE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_vendeur` int NOT NULL,
  `id_acheteur` int DEFAULT NULL,
  `id_casier` int DEFAULT NULL,
  `id_categorie` int DEFAULT NULL,
  `id_site` int DEFAULT NULL,
  `titre` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type_materiau` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `poids_estime_kg` decimal(10,2) DEFAULT NULL,
  `prix` decimal(10,2) DEFAULT NULL,
  `etat_objet` enum('Neuf','Bon etat','Usage') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `statut` enum('Disponible','Reserve','Depose','Paye','Recupere','Annule') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `est_valide` enum('En attente','Valide','Refuse') COLLATE utf8mb4_unicode_ci DEFAULT 'En attente',
  `code_barre_depot` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `code_barre_retrait` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  `date_depot_effective` datetime DEFAULT NULL,
  `date_achat` datetime DEFAULT NULL,
  `date_recuperation_effective` datetime DEFAULT NULL,
  `type` enum('Don','Vente','Service') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ville` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `code_postal` varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `provider` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `adresse` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_vendeur` (`id_vendeur`),
  KEY `id_acheteur` (`id_acheteur`),
  KEY `id_casier` (`id_casier`),
  KEY `id_categorie` (`id_categorie`),
  KEY `id_site` (`id_site`),
  CONSTRAINT `ANNONCE_ibfk_1` FOREIGN KEY (`id_vendeur`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `ANNONCE_ibfk_2` FOREIGN KEY (`id_acheteur`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `ANNONCE_ibfk_3` FOREIGN KEY (`id_casier`) REFERENCES `CASIER` (`id`),
  CONSTRAINT `ANNONCE_ibfk_4` FOREIGN KEY (`id_categorie`) REFERENCES `CATEGORIE` (`id`),
  CONSTRAINT `ANNONCE_ibfk_5` FOREIGN KEY (`id_site`) REFERENCES `SITE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `AVIS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AVIS` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_auteur` int NOT NULL,
  `id_cible` int NOT NULL,
  `note` int NOT NULL,
  `commentaire` text COLLATE utf8mb4_unicode_ci,
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_auteur` (`id_auteur`),
  KEY `id_cible` (`id_cible`),
  CONSTRAINT `AVIS_ibfk_1` FOREIGN KEY (`id_auteur`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `AVIS_ibfk_2` FOREIGN KEY (`id_cible`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `AVIS_chk_1` CHECK ((`note` between 1 and 5))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `CASIER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CASIER` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_conteneur` int NOT NULL,
  `numero_casier` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
  `taille` enum('Petit','Moyen','Grand') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `statut` enum('Libre','Reserve','Occupe','Maintenance') COLLATE utf8mb4_unicode_ci DEFAULT 'Libre',
  PRIMARY KEY (`id`),
  KEY `id_conteneur` (`id_conteneur`),
  CONSTRAINT `CASIER_ibfk_1` FOREIGN KEY (`id_conteneur`) REFERENCES `CONTENEUR` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `CATEGORIE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CATEGORIE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nom` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `id_parent` int DEFAULT NULL,
  `statut` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_parent` (`id_parent`),
  CONSTRAINT `CATEGORIE_ibfk_1` FOREIGN KEY (`id_parent`) REFERENCES `CATEGORIE` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `COMMANDE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMMANDE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_utilisateur` int NOT NULL,
  `montant_total` decimal(10,2) NOT NULL,
  `statut` enum('En attente','Payee','Annulee') COLLATE utf8mb4_unicode_ci DEFAULT 'En attente',
  `date_commande` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_utilisateur` (`id_utilisateur`),
  CONSTRAINT `COMMANDE_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `COMMENTAIRE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMMENTAIRE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_utilisateur` int NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_utilisateur` (`id_utilisateur`),
  CONSTRAINT `COMMENTAIRE_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `CONTENEUR`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CONTENEUR` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_site` int NOT NULL,
  `type_dechet` enum('Verre','Plastique','Metal','Papier','Electronique') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `statut` enum('Operationnel','Plein','Maintenance') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `capacite_max_kg` decimal(10,2) DEFAULT NULL,
  `niveau_remplissage` decimal(10,2) DEFAULT '0.00',
  PRIMARY KEY (`id`),
  KEY `id_site` (`id_site`),
  CONSTRAINT `CONTENEUR_ibfk_1` FOREIGN KEY (`id_site`) REFERENCES `SITE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `DM_CONVERSATION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DM_CONVERSATION` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_user_one` int NOT NULL,
  `id_user_two` int NOT NULL,
  `id_annonce` int DEFAULT NULL,
  `initiator_id` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_dm_user_one` (`id_user_one`),
  KEY `idx_dm_user_two` (`id_user_two`),
  KEY `idx_dm_annonce` (`id_annonce`),
  KEY `fk_dm_conversation_initiator` (`initiator_id`),
  CONSTRAINT `fk_dm_conversation_annonce` FOREIGN KEY (`id_annonce`) REFERENCES `ANNONCE` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_dm_conversation_initiator` FOREIGN KEY (`initiator_id`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_conversation_user_one` FOREIGN KEY (`id_user_one`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_conversation_user_two` FOREIGN KEY (`id_user_two`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `DM_MESSAGE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DM_MESSAGE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_conversation` int NOT NULL,
  `id_sender` int NOT NULL,
  `contenu` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `lu` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_dm_message_conversation` (`id_conversation`),
  KEY `fk_dm_message_sender` (`id_sender`),
  CONSTRAINT `fk_dm_message_conversation` FOREIGN KEY (`id_conversation`) REFERENCES `DM_CONVERSATION` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_message_sender` FOREIGN KEY (`id_sender`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `DM_OFFER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DM_OFFER` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_conversation` int NOT NULL,
  `id_annonce` int DEFAULT NULL,
  `id_buyer` int NOT NULL,
  `id_seller` int NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `status` enum('En attente','Acceptee','Refusee','Annulee') COLLATE utf8mb4_unicode_ci DEFAULT 'En attente',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_dm_offer_conversation` (`id_conversation`),
  KEY `fk_dm_offer_buyer` (`id_buyer`),
  KEY `fk_dm_offer_seller` (`id_seller`),
  KEY `fk_dm_offer_annonce` (`id_annonce`),
  CONSTRAINT `fk_dm_offer_annonce` FOREIGN KEY (`id_annonce`) REFERENCES `ANNONCE` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_dm_offer_buyer` FOREIGN KEY (`id_buyer`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_offer_conversation` FOREIGN KEY (`id_conversation`) REFERENCES `DM_CONVERSATION` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_offer_seller` FOREIGN KEY (`id_seller`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `DM_SALE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DM_SALE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_offer` int NOT NULL,
  `id_conversation` int NOT NULL,
  `id_annonce` int DEFAULT NULL,
  `id_buyer` int NOT NULL,
  `id_seller` int NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `status` enum('Offre acceptee','Payee','Recue','Evaluee') COLLATE utf8mb4_unicode_ci DEFAULT 'Offre acceptee',
  `received_at` datetime DEFAULT NULL,
  `reviewed_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_dm_sale_offer` (`id_offer`),
  KEY `fk_dm_sale_conversation` (`id_conversation`),
  KEY `fk_dm_sale_buyer` (`id_buyer`),
  KEY `fk_dm_sale_seller` (`id_seller`),
  KEY `fk_dm_sale_annonce` (`id_annonce`),
  CONSTRAINT `fk_dm_sale_annonce` FOREIGN KEY (`id_annonce`) REFERENCES `ANNONCE` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_dm_sale_buyer` FOREIGN KEY (`id_buyer`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_sale_conversation` FOREIGN KEY (`id_conversation`) REFERENCES `DM_CONVERSATION` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_sale_offer` FOREIGN KEY (`id_offer`) REFERENCES `DM_OFFER` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_dm_sale_seller` FOREIGN KEY (`id_seller`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ETAPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ETAPE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_projet` int NOT NULL,
  `numero_ordre` int NOT NULL,
  `titre` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `image_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_projet` (`id_projet`,`numero_ordre`),
  CONSTRAINT `ETAPE_ibfk_1` FOREIGN KEY (`id_projet`) REFERENCES `PROJET_UPCYCLING` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `EVENEMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `EVENEMENT` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_createur` int NOT NULL,
  `titre` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `adresse` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ville` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `code_postal` varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  `date_evenement` datetime DEFAULT NULL,
  `capacite_max` int DEFAULT '0',
  `statut` enum('planned','published','archived') COLLATE utf8mb4_unicode_ci DEFAULT 'planned',
  `type` enum('Atelier','Collecte','Conference','Echange') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_createur` (`id_createur`),
  CONSTRAINT `EVENEMENT_ibfk_1` FOREIGN KEY (`id_createur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `EVENEMENT_INSCRIPTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `EVENEMENT_INSCRIPTION` (
  `id_utilisateur` int NOT NULL,
  `id_evenement` int NOT NULL,
  `date_inscription` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_utilisateur`,`id_evenement`),
  KEY `id_evenement` (`id_evenement`),
  CONSTRAINT `EVENEMENT_INSCRIPTION_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `EVENEMENT_INSCRIPTION_ibfk_2` FOREIGN KEY (`id_evenement`) REFERENCES `EVENEMENT` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `FACTURE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FACTURE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_transaction` int NOT NULL,
  `numero_facture` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type_paiement` enum('Carte','Espece','Virement') COLLATE utf8mb4_unicode_ci DEFAULT 'Carte',
  `date_emission` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `numero_facture` (`numero_facture`),
  KEY `id_transaction` (`id_transaction`),
  CONSTRAINT `FACTURE_ibfk_1` FOREIGN KEY (`id_transaction`) REFERENCES `TRANSACTION` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `FAVORIS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FAVORIS` (
  `id_utilisateur` int NOT NULL,
  `id_annonce` int NOT NULL,
  `date_ajout` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_utilisateur`,`id_annonce`),
  KEY `id_annonce` (`id_annonce`),
  CONSTRAINT `FAVORIS_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FAVORIS_ibfk_2` FOREIGN KEY (`id_annonce`) REFERENCES `ANNONCE` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `FORMATION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FORMATION` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_formateur` int NOT NULL,
  `type` enum('Atelier','Cours','Webinaire') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `titre` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `capacite_max` int DEFAULT NULL,
  `est_valide` enum('En attente','Valide','Refuse') COLLATE utf8mb4_unicode_ci DEFAULT 'En attente',
  `date_debut` datetime DEFAULT NULL,
  `date_fin` datetime DEFAULT NULL,
  `statut` enum('Ouvert','Complet','Termine','Annule') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `prix_unitaire` decimal(10,2) DEFAULT NULL,
  `adresse` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ville` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `code_postal` varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_formateur` (`id_formateur`),
  CONSTRAINT `FORMATION_ibfk_1` FOREIGN KEY (`id_formateur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `FORMATION_INSCRIPTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FORMATION_INSCRIPTION` (
  `id_utilisateur` int NOT NULL,
  `id_formation` int NOT NULL,
  `date_inscription` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_utilisateur`,`id_formation`),
  KEY `id_formation` (`id_formation`),
  CONSTRAINT `FORMATION_INSCRIPTION_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FORMATION_INSCRIPTION_ibfk_2` FOREIGN KEY (`id_formation`) REFERENCES `FORMATION` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `FORUM`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FORUM` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_utilisateur` int NOT NULL,
  `id_salon` int DEFAULT NULL,
  `titre` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sujet` text COLLATE utf8mb4_unicode_ci,
  `ouvert` tinyint(1) DEFAULT '1',
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_utilisateur` (`id_utilisateur`),
  KEY `id_salon` (`id_salon`),
  CONSTRAINT `FORUM_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `FORUM_ibfk_2` FOREIGN KEY (`id_salon`) REFERENCES `FORUM_SALON` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `FORUM_MESSAGE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FORUM_MESSAGE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_utilisateur` int NOT NULL,
  `id_forum` int NOT NULL,
  `contenu` text COLLATE utf8mb4_unicode_ci,
  `date_envoi` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_utilisateur` (`id_utilisateur`),
  KEY `id_forum` (`id_forum`),
  CONSTRAINT `FORUM_MESSAGE_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `FORUM_MESSAGE_ibfk_2` FOREIGN KEY (`id_forum`) REFERENCES `FORUM` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `FORUM_SALON`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FORUM_SALON` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nom` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `LANGUE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `LANGUE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` char(2) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nom_langue` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `LIGNE_COMMANDE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `LIGNE_COMMANDE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_commande` int NOT NULL,
  `id_vendeur` int DEFAULT NULL,
  `type_item` enum('Formation','Annonce','Abonnement','Evenement') COLLATE utf8mb4_unicode_ci NOT NULL,
  `reference_id` int NOT NULL,
  `prix_unitaire` decimal(10,2) NOT NULL,
  `commission_upc` decimal(10,2) DEFAULT '0.00',
  PRIMARY KEY (`id`),
  KEY `id_commande` (`id_commande`),
  KEY `id_vendeur` (`id_vendeur`),
  CONSTRAINT `LIGNE_COMMANDE_ibfk_1` FOREIGN KEY (`id_commande`) REFERENCES `COMMANDE` (`id`) ON DELETE CASCADE,
  CONSTRAINT `LIGNE_COMMANDE_ibfk_2` FOREIGN KEY (`id_vendeur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `MESSAGE_SIGNALEMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `MESSAGE_SIGNALEMENT` (
  `id_message` int NOT NULL,
  `id_utilisateur` int NOT NULL,
  `motif` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date_signalement` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_message`,`id_utilisateur`),
  KEY `id_utilisateur` (`id_utilisateur`),
  CONSTRAINT `MESSAGE_SIGNALEMENT_ibfk_1` FOREIGN KEY (`id_message`) REFERENCES `FORUM_MESSAGE` (`id`) ON DELETE CASCADE,
  CONSTRAINT `MESSAGE_SIGNALEMENT_ibfk_2` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `NOTIFICATION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `NOTIFICATION` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_utilisateur` int NOT NULL,
  `id_emetteur` int DEFAULT '0',
  `type` enum('Alerte','Message','Rappel','Casier','Like','Avis','Follow') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `canal` enum('email','push','sms') COLLATE utf8mb4_unicode_ci DEFAULT 'email',
  `audience` enum('all','particuliers','prestataires','admins') COLLATE utf8mb4_unicode_ci DEFAULT 'all',
  `statut` enum('draft','scheduled','sent') COLLATE utf8mb4_unicode_ci DEFAULT 'sent',
  `titre` varchar(150) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `message` text COLLATE utf8mb4_unicode_ci,
  `lu` tinyint(1) DEFAULT '0',
  `date_envoi` datetime DEFAULT CURRENT_TIMESTAMP,
  `scheduled_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_utilisateur` (`id_utilisateur`),
  CONSTRAINT `NOTIFICATION_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `PANIER_ITEM`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PANIER_ITEM` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_utilisateur` int NOT NULL,
  `type_item` enum('Formation','Annonce','Abonnement','Evenement') COLLATE utf8mb4_unicode_ci NOT NULL,
  `reference_id` int NOT NULL,
  `prix_unitaire` decimal(10,2) NOT NULL,
  `date_ajout` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_utilisateur` (`id_utilisateur`),
  CONSTRAINT `PANIER_ITEM_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `PROJET_INSCRIPTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROJET_INSCRIPTION` (
  `id_utilisateur` int NOT NULL,
  `id_projet` int NOT NULL,
  `date_inscription` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_utilisateur`,`id_projet`),
  KEY `id_projet` (`id_projet`),
  CONSTRAINT `PROJET_INSCRIPTION_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `PROJET_INSCRIPTION_ibfk_2` FOREIGN KEY (`id_projet`) REFERENCES `PROJET_UPCYCLING` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `PROJET_LIKE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROJET_LIKE` (
  `id_utilisateur` int NOT NULL,
  `id_projet` int NOT NULL,
  `date_like` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_utilisateur`,`id_projet`),
  KEY `id_projet` (`id_projet`),
  CONSTRAINT `PROJET_LIKE_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE,
  CONSTRAINT `PROJET_LIKE_ibfk_2` FOREIGN KEY (`id_projet`) REFERENCES `PROJET_UPCYCLING` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `PROJET_UPCYCLING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROJET_UPCYCLING` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_createur` int NOT NULL,
  `image_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `titre` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description_courte` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  `score_impact` decimal(10,2) DEFAULT NULL,
  `nb_vues` int DEFAULT '0',
  `nb_likes` int DEFAULT '0',
  `co2_evite_kg` decimal(10,2) DEFAULT NULL,
  `visible_public` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `id_createur` (`id_createur`),
  CONSTRAINT `PROJET_UPCYCLING_ibfk_1` FOREIGN KEY (`id_createur`) REFERENCES `UTILISATEUR` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `PROJET_VUE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROJET_VUE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_projet` int NOT NULL,
  `id_utilisateur` int DEFAULT NULL,
  `ip_adresse` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `date_vue` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `id_projet` (`id_projet`),
  CONSTRAINT `PROJET_VUE_ibfk_1` FOREIGN KEY (`id_projet`) REFERENCES `PROJET_UPCYCLING` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `SITE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `SITE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nom` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ville` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `code_postal` varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `adresse` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `telephone` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` enum('Decheterie','Point de collecte','Association') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `actif` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `TIPS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `TIPS` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_createur` int NOT NULL,
  `titre` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `video_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `role_cible` enum('Particulier','Prestataire','Admin','Salarie') COLLATE utf8mb4_unicode_ci NOT NULL,
  `date_creation` datetime DEFAULT CURRENT_TIMESTAMP,
  `actif` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `id_createur` (`id_createur`),
  CONSTRAINT `TIPS_ibfk_1` FOREIGN KEY (`id_createur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `TRADUCTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `TRADUCTION` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_langue` int DEFAULT NULL,
  `cle_traduction` text COLLATE utf8mb4_unicode_ci,
  `text_traduit` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  KEY `id_langue` (`id_langue`),
  CONSTRAINT `TRADUCTION_ibfk_1` FOREIGN KEY (`id_langue`) REFERENCES `LANGUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `TRANSACTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `TRANSACTION` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_acheteur` int NOT NULL,
  `id_commande` int NOT NULL,
  `montant_total` decimal(10,2) DEFAULT NULL,
  `statut_paiement` enum('En attente','Valide','Echoue') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date_transaction` datetime DEFAULT CURRENT_TIMESTAMP,
  `stripe_payment_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `montant_ht` decimal(10,2) DEFAULT NULL,
  `type` enum('Formation','Annonce','Abonnement') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_acheteur` (`id_acheteur`),
  KEY `id_commande` (`id_commande`),
  CONSTRAINT `TRANSACTION_ibfk_1` FOREIGN KEY (`id_acheteur`) REFERENCES `UTILISATEUR` (`id`),
  CONSTRAINT `TRANSACTION_ibfk_2` FOREIGN KEY (`id_commande`) REFERENCES `COMMANDE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `TYPE_ABONNEMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `TYPE_ABONNEMENT` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nom` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `prix_ht` decimal(10,2) NOT NULL,
  `duree_mois` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `UPCYCLING_SCORE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `UPCYCLING_SCORE` (
  `id` int NOT NULL AUTO_INCREMENT,
  `id_utilisateur` int NOT NULL,
  `ressources_economisees` decimal(10,2) DEFAULT NULL,
  `co2_total_evite_kg` decimal(10,2) DEFAULT NULL,
  `nb_objets_recycles` int DEFAULT '0',
  `total_points` int DEFAULT '0',
  `niveau` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_utilisateur` (`id_utilisateur`),
  CONSTRAINT `UPCYCLING_SCORE_ibfk_1` FOREIGN KEY (`id_utilisateur`) REFERENCES `UTILISATEUR` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `UTILISATEUR`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `UTILISATEUR` (
  `id` int NOT NULL AUTO_INCREMENT,
  `prenom` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nom` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `image_profil` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `siret` varchar(14) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `siret_valide` tinyint(1) DEFAULT NULL,
  `mail_valide` tinyint(1) DEFAULT NULL,
  `banniere` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mail` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `adresse` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ville` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_admin` tinyint(1) DEFAULT '0',
  `code_postal` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date_naissance` date DEFAULT NULL,
  `date_inscription` datetime DEFAULT CURRENT_TIMESTAMP,
  `role` enum('Particulier','Prestataire','Admin','Salarie') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `statut` enum('Actif','Inactif','Banni','En attente de validation') COLLATE utf8mb4_unicode_ci DEFAULT 'Actif',
  `id_langue` int DEFAULT NULL,
  `token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ban_forum` tinyint(1) DEFAULT '0',
  `materiaux_recherches` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'Mots-cl??s libres s??par??s par des virgules, ex: bois,m??tal,palette',
  `date_update_password` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `mail` (`mail`),
  UNIQUE KEY `token` (`token`),
  KEY `id_langue` (`id_langue`),
  CONSTRAINT `UTILISATEUR_ibfk_1` FOREIGN KEY (`id_langue`) REFERENCES `LANGUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

