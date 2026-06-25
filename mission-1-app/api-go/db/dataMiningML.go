package db

import (
	"log"
)

type StatActeur struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

type StatPrestation struct {
	TypeItem string `json:"type_item"`
	Count    int    `json:"count"`
}

type PredictionML struct {
	IDUtilisateur         int     `json:"id_utilisateur"`
	Prenom                string  `json:"prenom"`
	Nom                   string  `json:"nom"`
	Role                  string  `json:"role"`
	PrestationRecommandee string  `json:"prestation_recommandee"`
	Confiance             float64 `json:"confiance"`
}

func GetStatsActeurs() ([]StatActeur, error) {
	rows, err := Conn.Query("SELECT role, COUNT(*) FROM UTILISATEUR GROUP BY role")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []StatActeur
	for rows.Next() {
		var s StatActeur
		rows.Scan(&s.Role, &s.Count)
		stats = append(stats, s)
	}
	return stats, nil
}

func GetStatsPrestations() ([]StatPrestation, error) {
	rows, err := Conn.Query("SELECT type, COUNT(*) FROM ANNONCE GROUP BY type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []StatPrestation
	for rows.Next() {
		var s StatPrestation
		rows.Scan(&s.TypeItem, &s.Count)
		stats = append(stats, s)
	}
	return stats, nil
}

func GetMLPredictions() ([]PredictionML, error) {
	query := `
		SELECT p.id_utilisateur, u.prenom, u.nom, u.role, p.prestation_recommandee, p.confiance 
		FROM ML_PREDICTION p
		JOIN UTILISATEUR u ON p.id_utilisateur = u.id
		ORDER BY p.confiance DESC
		LIMIT 20
	`
	rows, err := Conn.Query(query)
	if err != nil {
		log.Println("Erreur lecture ML_PREDICTION:", err)
		return nil, err
	}
	defer rows.Close()

	var predictions []PredictionML
	for rows.Next() {
		var p PredictionML
		rows.Scan(&p.IDUtilisateur, &p.Prenom, &p.Nom, &p.Role, &p.PrestationRecommandee, &p.Confiance)
		predictions = append(predictions, p)
	}
	return predictions, nil
}