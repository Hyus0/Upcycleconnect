ALTER TABLE ANNONCE
  MODIFY COLUMN type ENUM('Don', 'Vente', 'Service');

ALTER TABLE ANNONCE
  ADD COLUMN provider VARCHAR(150) NULL AFTER code_postal;

ALTER TABLE EVENEMENT
  ADD COLUMN capacite_max INT DEFAULT 0 AFTER date_evenement;

ALTER TABLE EVENEMENT
  ADD COLUMN statut ENUM('planned', 'published', 'archived') DEFAULT 'planned' AFTER capacite_max;

ALTER TABLE NOTIFICATION
  ADD COLUMN canal ENUM('email', 'push', 'sms') DEFAULT 'email' AFTER type;

ALTER TABLE NOTIFICATION
  ADD COLUMN audience ENUM('all', 'particuliers', 'prestataires', 'admins') DEFAULT 'all' AFTER canal;

ALTER TABLE NOTIFICATION
  ADD COLUMN statut ENUM('draft', 'scheduled', 'sent') DEFAULT 'sent' AFTER audience;

ALTER TABLE NOTIFICATION
  ADD COLUMN scheduled_at DATETIME NULL AFTER date_envoi;
