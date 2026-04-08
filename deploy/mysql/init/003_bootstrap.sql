INSERT INTO LANGUE (code, nom_langue)
SELECT 'fr', 'Francais'
WHERE NOT EXISTS (SELECT 1 FROM LANGUE WHERE code = 'fr');

INSERT INTO UTILISATEUR (
  prenom, nom, password, mail, adresse, ville, code_postal, date_naissance,
  role, statut, id_langue
)
SELECT 'Admin', 'UpcycleConnect', 'AdminTemp1!', 'admin@upcycleconnect.local', '', 'Paris', '75011', NULL, 'Admin', 'Actif', 1
WHERE NOT EXISTS (SELECT 1 FROM UTILISATEUR WHERE mail = 'admin@upcycleconnect.local');

INSERT INTO UTILISATEUR (
  prenom, nom, password, mail, adresse, ville, code_postal, date_naissance,
  role, statut, id_langue
)
SELECT 'Marie', 'Lambert', 'Particulier1!', 'marie.lambert@upcycleconnect.local', '', 'Paris', '75011', NULL, 'Particulier', 'Actif', 1
WHERE NOT EXISTS (SELECT 1 FROM UTILISATEUR WHERE mail = 'marie.lambert@upcycleconnect.local');

INSERT INTO UTILISATEUR (
  prenom, nom, password, mail, adresse, ville, code_postal, date_naissance,
  role, statut, id_langue
)
SELECT 'Atelier', 'Renouveau', 'Prestataire1!', 'atelier.renouveau@upcycleconnect.local', '', 'Lyon', '69002', NULL, 'Prestataire', 'Actif', 1
WHERE NOT EXISTS (SELECT 1 FROM UTILISATEUR WHERE mail = 'atelier.renouveau@upcycleconnect.local');
