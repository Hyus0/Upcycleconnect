package models

type ProjetUpcycling struct {
	ID                int           `json:"id"`
	IdCreateur        int           `json:"id_createur"`
	Titre             string        `json:"titre"`
	DescriptionCourte string        `json:"description_courte"`
	ImageUrl          string        `json:"image_url"`
	ScoreImpact       float64       `json:"score_impact"`
	Co2EviteKg        float64       `json:"co2_evite_kg"`
	NbVues            int           `json:"nb_vues"`
	NbLikes           int           `json:"nb_likes"`
	VisiblePublic     bool          `json:"visible_public"`
	DateCreation      string        `json:"date_creation"`
	Etapes            []EtapeProjet `json:"etapes,omitempty"`
}

type EtapeProjet struct {
	ID          int    `json:"id"`
	IdProjet    int    `json:"id_projet"`
	NumeroOrdre int    `json:"numero_ordre"`
	Titre       string `json:"titre"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}