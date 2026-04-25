package models

type Site struct {
	ID         int    `json:"id"`
	Nom        string `json:"nom"`
	Ville      string `json:"ville"`
	CodePostal string `json:"code_postal"`
	Adresse    string `json:"adresse"`
	Telephone  string `json:"telephone"`
	Type       string `json:"type"`
	Actif      bool   `json:"actif"`
}

type Conteneur struct {
	ID                int     `json:"id"`
	IdSite            int     `json:"id_site"`
	TypeDechet        string  `json:"type_dechet"` 
	Statut            string  `json:"statut"`  
	CapaciteMaxKg     float64 `json:"capacite_max_kg"`
	NiveauRemplissage float64 `json:"niveau_remplissage"`
}

type Casier struct {
	ID           int    `json:"id"`
	IdConteneur  int    `json:"id_conteneur"`
	NumeroCasier string `json:"numero_casier"`
	Taille       string `json:"taille"`
	Statut       string `json:"statut"` 
}