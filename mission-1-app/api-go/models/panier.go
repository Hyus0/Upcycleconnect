package models

type PanierItem struct {
	ID            int     `json:"id"`
	IdUtilisateur int     `json:"id_utilisateur"`
	TypeItem      string  `json:"type_item"`
	ReferenceID   int     `json:"reference_id"`
	PrixUnitaire  float64 `json:"prix_unitaire"`
	DateAjout     string  `json:"date_ajout"`
	Titre string `json:"titre,omitempty"`
}

type Commande struct {
	ID            int     `json:"id"`
	IdUtilisateur int     `json:"id_utilisateur"`
	MontantTotal  float64 `json:"montant_total"`
	Statut        string  `json:"statut"`
	DateCommande  string  `json:"date_commande"`
}

type Transaction struct {
	ID              int     `json:"id"`
	IdAcheteur      int     `json:"id_acheteur"`
	IdCommande      int     `json:"id_commande"`
	MontantTotal    float64 `json:"montant_total"`
	StatutPaiement  string  `json:"statut_paiement"`
	DateTransaction string  `json:"date_transaction"`
	StripePaymentID string  `json:"stripe_payment_id,omitempty"`
}