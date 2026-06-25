-- 002_admin_extensions.sql
-- Les extensions admin (ANNONCE.provider, EVENEMENT.capacite_max/statut,
-- NOTIFICATION.canal/audience/statut/scheduled_at, TRANSACTION.montant_ht/type)
-- sont desormais integrees directement dans 001_schema.sql.
-- Ce fichier est conserve volontairement vide pour preserver l'ordre des scripts.
SELECT 'admin extensions deja incluses dans 001_schema.sql' AS info;
