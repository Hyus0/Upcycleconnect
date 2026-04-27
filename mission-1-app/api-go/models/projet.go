package models

type ProjetUpcycling struct {
	ID                int       `json:"id"`
	IdCreateur        int       `json:"id_createur"`
	Titre             string    `json:"titre"`
	DescriptionCourte string    `json:"description_courte"`
	DateCreation      string `json:"date_creation"`
	ScoreImpact       float64   `json:"score_impact"`
	NbVues            int       `json:"nb_vues"`
	NbLikes           int       `json:"nb_likes"`
	Co2EviteKg        float64   `json:"co2_evite_kg"`
	VisiblePublic     bool      `json:"visible_public"`
	Etapes            []Etape   `json:"etapes,omitempty"` 
}

type Etape struct {
	ID          int    `json:"id"`
	IdProjet    int    `json:"id_projet"`
	NumeroOrdre int    `json:"numero_ordre"`
	Titre       string `json:"titre"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}
