package models

type Notification struct {
	ID            int    `json:"id"`
	IDUtilisateur int    `json:"id_utilisateur"`
	Type          string `json:"type"`
	Titre         string `json:"titre"`
	Message       string `json:"message"`
	Lu            bool   `json:"lu"`
	DateEnvoi     string `json:"date_envoi"`
}