package models

type GetFormation struct {
	ID            int     `json:"id"`
	ID_formateur  int     `json:"id_formateur"`
	Prenom_formateur string    `json:"prenom_formateur"`
    Nom_formateur    string    `json:"nom_formateur"`
    Image_formateur	string 	`json:"image_formateur"`
    IsRegistered    bool	`json:"is_registered"`
	Type          string  `json:"type"`
	Titre         string  `json:"titre"`
	Description   string  `json:"description"`
	Capacite_max  int     `json:"capacite_max"`
	Est_valide	  string  `json: "est_valide"`
	Nb_inscrit    int     `json:"nb_inscrit"`
	Date_debut    string  `json:"date_debut"`
	Date_fin      string  `json:"date_fin"`
	Statut        string  `json:"statut"`
	Prix_unitaire float64 `json:"prix_unitaire"`
	Adresse       string  `json:"adresse"`
	Ville         string  `json:"ville"`
	CodePostal    string  `json:"code_postal"`
	Sessions []FormationSession `json:"sessions"`
}

type Formation struct {
	ID            int     `json:"id"`
	ID_formateur  int     `json:"id_formateur"`
	Type          string  `json:"type"`
	Titre         string  `json:"titre"`
	Description   string  `json:"description"`
	Capacite_max  int     `json:"capacite_max"`
	Est_valide	  string  `json: "est_valide"`
	Date_debut    string  `json:"date_debut"`
	Date_fin      string  `json:"date_fin"`
	Statut        string  `json:"statut"`
	Prix_unitaire float64 `json:"prix_unitaire"`
	Adresse       string  `json:"adresse"`
	Ville         string  `json:"ville"`
	Code_postal   string  `json:"code_postal"`
	Sessions []FormationSession `json:"sessions"`
}

type Participant struct {
	ID          int    `json:"id"`
	Prenom      string `json:"prenom"`
	Nom         string `json:"nom"`
	Mail		string `json:"mail"`
	ImageProfil string `json:"image_profil"`
	Role        string `json:"role"`
}

type FormationSession struct {
	ID        int    `json:"id"`
	Nom       string `json:"nom"`
	DateDebut string `json:"date_debut"`
	DateFin   string `json:"date_fin"`
	Statut    string `json:"statut"`
}