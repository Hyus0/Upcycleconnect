package models

type Avis struct {
	ID           int    `json:"id"`
	IdAuteur     int    `json:"id_auteur"`
	IdCible      int    `json:"id_cible"`
	Note         int    `json:"note"`
	Commentaire  string `json:"commentaire"`
	DateCreation string `json:"date_creation"`
	PrenomAuteur string `json:"prenom_auteur,omitempty"`
}