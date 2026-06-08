package db

import (
	"database/sql"
	"upcycleconnect/api-go/models"
)

func GetPlatformStats() (models.PlatformStats, error) {
	var stats models.PlatformStats
	
	err := Conn.QueryRow(`
		SELECT COALESCE(SUM(s.co2_total_evite_kg), 0) 
		FROM UPCYCLING_SCORE s
		JOIN UTILISATEUR u ON s.id_utilisateur = u.id
		WHERE u.role = 'Particulier'
	`).Scan(&stats.Co2Evite)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	err = Conn.QueryRow(`
		SELECT COALESCE(SUM(poids_estime_kg * 1.5), 0) 
		FROM ANNONCE 
		WHERE statut = 'Recupere' 
		AND date_recuperation_effective >= DATE_SUB(NOW(), INTERVAL 30 DAY)
	`).Scan(&stats.Co2EviteMois)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	err = Conn.QueryRow(`
		SELECT COALESCE(SUM(s.nb_objets_recycles), 0) 
		FROM UPCYCLING_SCORE s
		JOIN UTILISATEUR u ON s.id_utilisateur = u.id
		WHERE u.role = 'Particulier'
	`).Scan(&stats.ObjetsUpcycles)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	err = Conn.QueryRow("SELECT COUNT(*) FROM UTILISATEUR WHERE role = 'Prestataire' AND statut = 'Actif'").Scan(&stats.ArtisansActifs)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	err = Conn.QueryRow("SELECT COUNT(*) FROM SITE WHERE actif = 1").Scan(&stats.SitesActifs)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	return stats, nil
}