package models

type AbonnementUtilisateur struct {
	IdAbonne       int    `json:"id_abonne"`
	IdSuivi        int    `json:"id_suivi"`
	DateAbonnement string `json:"date_abonnement"`
}