package models

type CO2ParMois struct {
	Mois   string  `json:"mois"`
	Valeur float64 `json:"valeur"`
}

type EcoStats struct {
	CO2Total           float64      `json:"co2_total"`
	CO2Trend           float64      `json:"co2_trend"`
	EauEconomisee      float64      `json:"eau_economisee"`
	MateriauxValorises int          `json:"materiaux_valorises"`
	ScoreImpactMoyen   float64      `json:"score_impact_moyen"`
	CO2ParMois         []CO2ParMois `json:"co2_par_mois"`
}


type MateriauStat struct {
	TypeMateriau string `json:"type_materiau"`
	Count        int    `json:"count"`
	Zone         string `json:"zone"`
}


type AlertePrioritaire struct {
	ID            int     `json:"id"`
	TitreAnnonce  string  `json:"titre_annonce"`
	TypeMateriau  string  `json:"type_materiau"`
	NomSite       string  `json:"nom_site"`
	Ville         string  `json:"ville"`
	DateCreation  string  `json:"date_creation"`
	MatchScore    int     `json:"match_score"`
	DistanceLabel string  `json:"distance_km"`
}