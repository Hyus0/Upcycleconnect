package db

import "fmt"

func RepartitionRoles() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT role, COUNT(*) AS total
		FROM UTILISATEUR
		WHERE statut = 'Actif'
		GROUP BY role
		ORDER BY total DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("RepartitionRoles: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var role string
		var total int
		if err := rows.Scan(&role, &total); err == nil {
			result = append(result, map[string]interface{}{
				"role":  role,
				"total": total,
			})
		}
	}
	return result, nil
}

func InscriptionsParMois() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT
			DATE_FORMAT(date_inscription, '%b %Y') AS mois,
			COUNT(*) AS total
		FROM UTILISATEUR
		WHERE date_inscription >= CURDATE() - INTERVAL 12 MONTH
		GROUP BY YEAR(date_inscription), MONTH(date_inscription)
		ORDER BY YEAR(date_inscription) ASC, MONTH(date_inscription) ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("InscriptionsParMois: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var mois string
		var total int
		if err := rows.Scan(&mois, &total); err == nil {
			result = append(result, map[string]interface{}{
				"mois":  mois,
				"total": total,
			})
		}
	}
	return result, nil
}

func VillesActeurs() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT
			COALESCE(NULLIF(TRIM(ville), ''), 'Non renseigné') AS ville,
			COUNT(*) AS total
		FROM UTILISATEUR
		WHERE statut = 'Actif'
		GROUP BY ville
		ORDER BY total DESC
		LIMIT 8
	`)
	if err != nil {
		return nil, fmt.Errorf("VillesActeurs: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var ville string
		var total int
		if err := rows.Scan(&ville, &total); err == nil {
			result = append(result, map[string]interface{}{
				"ville": ville,
				"total": total,
			})
		}
	}
	return result, nil
}

func StatutsUtilisateurs() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT statut, COUNT(*) AS total
		FROM UTILISATEUR
		GROUP BY statut
	`)
	if err != nil {
		return nil, fmt.Errorf("StatutsUtilisateurs: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var statut string
		var total int
		if err := rows.Scan(&statut, &total); err == nil {
			result = append(result, map[string]interface{}{
				"statut": statut,
				"total":  total,
			})
		}
	}
	return result, nil
}

func StatsFormations() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT
			f.titre,
			f.type,
			COUNT(fi.id_formation) AS nb_inscrits,
			f.capacite_max
		FROM FORMATION f
		LEFT JOIN FORMATION_INSCRIPTION fi ON fi.id_formation = f.id
		WHERE f.est_valide = 'Valide'
		GROUP BY f.id, f.titre, f.type, f.capacite_max
		ORDER BY nb_inscrits DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("StatsFormations: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var titre, typeF string
		var nbInscrits, capacite int
		if err := rows.Scan(&titre, &typeF, &nbInscrits, &capacite); err == nil {
			tauxRemplissage := 0.0
			if capacite > 0 {
				tauxRemplissage = float64(nbInscrits) / float64(capacite) * 100
			}
			result = append(result, map[string]interface{}{
				"titre":            titre,
				"type":             typeF,
				"nb_inscrits":      nbInscrits,
				"capacite_max":     capacite,
				"taux_remplissage": tauxRemplissage,
			})
		}
	}
	return result, nil
}

func StatsEvenements() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT
			e.titre,
			e.type,
			COUNT(ei.id_evenement) AS nb_inscrits
		FROM EVENEMENT e
		LEFT JOIN EVENEMENT_INSCRIPTION ei ON ei.id_evenement = e.id
		GROUP BY e.id, e.titre, e.type
		ORDER BY nb_inscrits DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("StatsEvenements: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var titre, typeE string
		var nbInscrits int
		if err := rows.Scan(&titre, &typeE, &nbInscrits); err == nil {
			result = append(result, map[string]interface{}{
				"titre":       titre,
				"type":        typeE,
				"nb_inscrits": nbInscrits,
			})
		}
	}
	return result, nil
}

func StatsAnnonces() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT
			type,
			statut,
			COUNT(*) AS total
		FROM ANNONCE
		GROUP BY type, statut
		ORDER BY type, total DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("StatsAnnonces: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var typeA, statut string
		var total int
		if err := rows.Scan(&typeA, &statut, &total); err == nil {
			result = append(result, map[string]interface{}{
				"type":   typeA,
				"statut": statut,
				"total":  total,
			})
		}
	}
	return result, nil
}

func ChiffreAffairesParMois() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT
			DATE_FORMAT(date_commande, '%b %Y') AS mois,
			SUM(montant_total) AS ca
		FROM COMMANDE
		WHERE statut = 'Payee'
		  AND date_commande >= CURDATE() - INTERVAL 12 MONTH
		GROUP BY YEAR(date_commande), MONTH(date_commande)
		ORDER BY YEAR(date_commande) ASC, MONTH(date_commande) ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("ChiffreAffairesParMois: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var mois string
		var ca float64
		if err := rows.Scan(&mois, &ca); err == nil {
			result = append(result, map[string]interface{}{
				"mois": mois,
				"ca":   ca,
			})
		}
	}
	return result, nil
}

func TypesFormations() ([]map[string]interface{}, error) {
	rows, err := Conn.Query(`
		SELECT type, COUNT(*) AS total
		FROM FORMATION
		WHERE est_valide = 'Valide'
		GROUP BY type
	`)
	if err != nil {
		return nil, fmt.Errorf("TypesFormations: %v", err)
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var typeF string
		var total int
		if err := rows.Scan(&typeF, &total); err == nil {
			result = append(result, map[string]interface{}{
				"type":  typeF,
				"total": total,
			})
		}
	}
	return result, nil
}