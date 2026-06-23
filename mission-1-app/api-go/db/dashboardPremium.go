package db

import (
	"fmt"
	"strings"
	"upcycleconnect/api-go/models"
)

func GetEcoStats(userID int) (*models.EcoStats, error) {
	stats := &models.EcoStats{CO2ParMois: []models.CO2ParMois{}}

	err := Conn.QueryRow(`
		SELECT COALESCE(SUM(co2_evite_kg), 0)
		FROM PROJET
		WHERE id_createur = ?
	`, userID).Scan(&stats.CO2Total)
	if err != nil {
		return nil, fmt.Errorf("GetEcoStats co2_total: %v", err)
	}

	var co2MoisActuel, co2MoisPrecedent float64
	err = Conn.QueryRow(`
		SELECT COALESCE(SUM(co2_evite_kg), 0)
		FROM PROJET
		WHERE id_createur = ?
		  AND YEAR(date_creation) = YEAR(CURDATE())
		  AND MONTH(date_creation) = MONTH(CURDATE())
	`, userID).Scan(&co2MoisActuel)
	if err != nil {
		return nil, fmt.Errorf("GetEcoStats co2_mois_actuel: %v", err)
	}

	err = Conn.QueryRow(`
		SELECT COALESCE(SUM(co2_evite_kg), 0)
		FROM PROJET
		WHERE id_createur = ?
		  AND YEAR(date_creation) = YEAR(CURDATE() - INTERVAL 1 MONTH)
		  AND MONTH(date_creation) = MONTH(CURDATE() - INTERVAL 1 MONTH)
	`, userID).Scan(&co2MoisPrecedent)
	if err != nil {
		return nil, fmt.Errorf("GetEcoStats co2_mois_precedent: %v", err)
	}

	if co2MoisPrecedent > 0 {
		stats.CO2Trend = ((co2MoisActuel - co2MoisPrecedent) / co2MoisPrecedent) * 100
	} else if co2MoisActuel > 0 {
		stats.CO2Trend = 100
	}

	stats.EauEconomisee = stats.CO2Total * 12

	err = Conn.QueryRow(`
		SELECT COUNT(*)
		FROM PROJET
		WHERE id_createur = ?
	`, userID).Scan(&stats.MateriauxValorises)
	if err != nil {
		return nil, fmt.Errorf("GetEcoStats materiaux_valorises: %v", err)
	}

	err = Conn.QueryRow(`
		SELECT COALESCE(AVG(score_impact), 0)
		FROM PROJET
		WHERE id_createur = ?
	`, userID).Scan(&stats.ScoreImpactMoyen)
	if err != nil {
		return nil, fmt.Errorf("GetEcoStats score_impact_moyen: %v", err)
	}

	rows, err := Conn.Query(`
		SELECT
			DATE_FORMAT(date_creation, '%b %Y') AS mois,
			COALESCE(SUM(co2_evite_kg), 0) AS total
		FROM PROJET
		WHERE id_createur = ?
		  AND date_creation >= CURDATE() - INTERVAL 6 MONTH
		GROUP BY YEAR(date_creation), MONTH(date_creation)
		ORDER BY YEAR(date_creation) ASC, MONTH(date_creation) ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("GetEcoStats co2_par_mois: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.CO2ParMois
		if err := rows.Scan(&m.Mois, &m.Valeur); err != nil {
			return nil, fmt.Errorf("GetEcoStats scan co2_par_mois: %v", err)
		}
		stats.CO2ParMois = append(stats.CO2ParMois, m)
	}

	return stats, nil
}

func GetMateriauxStats() ([]models.MateriauStat, error) {
	rows, err := Conn.Query(`
		SELECT
			COALESCE(NULLIF(TRIM(type_materiau), ''), 'Non précisé') AS type_materiau,
			COUNT(*) AS total,
			COALESCE(NULLIF(TRIM(ville), ''), 'Multi-zones') AS zone
		FROM ANNONCE
		WHERE statut = 'Disponible'
		  AND est_valide = 'Valide'
		GROUP BY type_materiau, ville
		ORDER BY total DESC
		LIMIT 20
	`)
	if err != nil {
		return nil, fmt.Errorf("GetMateriauxStats: %v", err)
	}
	defer rows.Close()

	stats := []models.MateriauStat{}
	for rows.Next() {
		var s models.MateriauStat
		if err := rows.Scan(&s.TypeMateriau, &s.Count, &s.Zone); err != nil {
			return nil, fmt.Errorf("GetMateriauxStats scan: %v", err)
		}
		stats = append(stats, s)
	}

	return mergeMateriauxByType(stats), nil
}

func mergeMateriauxByType(raw []models.MateriauStat) []models.MateriauStat {
	merged := make(map[string]*models.MateriauStat)
	var order []string

	for _, item := range raw {
		key := strings.ToLower(item.TypeMateriau)
		if existing, ok := merged[key]; ok {
			existing.Count += item.Count
		} else {
			copyItem := item
			merged[key] = &copyItem
			order = append(order, key)
		}
	}

	result := make([]models.MateriauStat, 0, len(order))
	for _, key := range order {
		result = append(result, *merged[key])
	}
	return result
}

func GetAlertesPrioritaires(userID int) ([]models.AlertePrioritaire, error) {
	var materiauxRecherches, villeUser string
	err := Conn.QueryRow(`
		SELECT COALESCE(materiaux_recherches, ''), COALESCE(ville, '')
		FROM UTILISATEUR
		WHERE id = ?
	`, userID).Scan(&materiauxRecherches, &villeUser)
	if err != nil {
		return nil, fmt.Errorf("GetAlertesPrioritaires user: %v", err)
	}

	keywords := splitKeywords(materiauxRecherches)

	var likeConditions []string
	var likeArgs []interface{}
	for _, kw := range keywords {
		cleaned := strings.TrimSpace(kw)
		if cleaned == "" {
			continue
		}
		likeConditions = append(likeConditions, `(
			LOWER(a.type_materiau) LIKE ?
			OR ? LIKE CONCAT('%', LOWER(a.type_materiau), '%')
		)`)
		pattern := "%" + strings.ToLower(cleaned) + "%"
		likeArgs = append(likeArgs, pattern, strings.ToLower(cleaned))
	}

	matchClause := "0"
	if len(likeConditions) > 0 {
		matchClause = "(" + strings.Join(likeConditions, " OR ") + ")"
	}

	query := fmt.Sprintf(`
		SELECT
			a.id,
			a.titre,
			COALESCE(NULLIF(TRIM(a.type_materiau), ''), 'Non précisé') AS type_materiau,
			COALESCE(s.nom, 'Non assigné') AS nom_site,
			COALESCE(a.ville, '') AS ville,
			a.date_creation,
			(
				CASE WHEN %s THEN 1 ELSE 0 END
				+ CASE WHEN a.ville IS NOT NULL AND a.ville <> '' AND a.ville = ? THEN 1 ELSE 0 END
			) AS match_score
		FROM ANNONCE a
		LEFT JOIN CONTENEUR c ON c.id = a.id_conteneur
		LEFT JOIN SITE s ON s.id = c.id_site
		WHERE a.statut = 'Disponible'
		  AND a.est_valide = 'Valide'
		  AND a.id_vendeur != ?
		HAVING match_score > 0
		ORDER BY match_score DESC, a.date_creation DESC
		LIMIT 15
	`, matchClause)

	args := append(likeArgs, villeUser, userID)

	rows, err := Conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAlertesPrioritaires query: %v", err)
	}
	defer rows.Close()

	alertes := []models.AlertePrioritaire{}
	for rows.Next() {
		var a models.AlertePrioritaire
		var ville string
		if err := rows.Scan(
			&a.ID, &a.TitreAnnonce, &a.TypeMateriau, &a.NomSite,
			&ville, &a.DateCreation, &a.MatchScore,
		); err != nil {
			return nil, fmt.Errorf("GetAlertesPrioritaires scan: %v", err)
		}
		a.Ville = ville
		if ville != "" && ville == villeUser {
			a.DistanceLabel = "Même ville"
		} else if ville != "" {
			a.DistanceLabel = ville
		} else {
			a.DistanceLabel = "?"
		}
		alertes = append(alertes, a)
	}

	return alertes, nil
}

func splitKeywords(raw string) []string {
	raw = strings.ReplaceAll(raw, ",", " ")
	parts := strings.Fields(raw)
	keywords := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if len(trimmed) >= 3 { 
			keywords = append(keywords, trimmed)
		}
	}
	return keywords
}