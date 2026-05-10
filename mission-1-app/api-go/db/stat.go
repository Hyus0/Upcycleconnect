package db

import (
	"database/sql"
	"upcycleconnect/api-go/models"
)

func GetPlatformStats() (models.PlatformStats, error) {
	var stats models.PlatformStats

	err := Conn.QueryRow("SELECT COALESCE(SUM(co2_total_evite_kg), 0) FROM UPCYCLING_SCORE").Scan(&stats.Co2Evite)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	err = Conn.QueryRow("SELECT COALESCE(SUM(nb_objets_recycles), 0) FROM UPCYCLING_SCORE").Scan(&stats.ObjetsUpcycles)
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