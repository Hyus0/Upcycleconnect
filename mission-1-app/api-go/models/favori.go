package models

type Favori struct {
	IdUtilisateur int    `json:"id_utilisateur"`
	IdAnnonce     int    `json:"id_annonce"`
	DateAjout     string `json:"date_ajout"`
}

type FavoriStatusResponse struct {
	Total       int  `json:"total"`      
	IsFavorited bool `json:"is_favorited"`
}