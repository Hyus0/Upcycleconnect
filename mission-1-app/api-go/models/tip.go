package models

type Tip struct {
	ID           int    `json:"id"`
	ID_createur  int    `json:"id_createur"`
	Titre        string `json:"titre"`
	Description  string `json:"description"`
	Video_url     string `json:"video_url"`
	Role_cible   string `json:"role_cible"`
	Date_creation string `json:"date_creation"`
	Actif        bool   `json:"actif"`
}