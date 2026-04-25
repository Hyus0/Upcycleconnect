package models

type Annonce struct {
	ID                        int        `json:"id"`
	IdVendeur                 int        `json:"id_vendeur"`
	IdAcheteur                *int       `json:"id_acheteur,omitempty"`
	IdCasier                  *int       `json:"id_casier,omitempty"`
	IdCategorie               int        `json:"id_categorie"`
	Titre                     string     `json:"titre"`
	Description               string     `json:"description"`
	TypeMateriau              string     `json:"type_materiau"`
	PoidsEstimeKg             float64    `json:"poids_estime_kg"`
	Prix                      float64    `json:"prix"`
	EtatObjet                 string     `json:"etat_objet"`
	Statut                    string     `json:"statut"`
	EstValide                 string     `json:"est_valide"`
	CodePinDepot              *string     `json:"code_pin_depot"`
	IdSite       			  int    `json:"id_site"`
	CodeBarreRetrait          *string     `json:"code_barre_retrait"`
	DateCreation              *string `json:"date_creation,omitempty"`
	DateDepotEffective        *string `json:"date_depot_effective,omitempty"`
	DateRecuperationEffective *string `json:"date_recuperation_effective,omitempty"`
	Type                      string     `json:"type"`
	Ville                     string     `json:"ville"`
	CodePostal                string     `json:"code_postal"`
	Adresse                   string     `json:"adresse"`
}
