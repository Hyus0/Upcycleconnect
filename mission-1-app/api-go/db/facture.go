package db

import (
	"database/sql"
	"fmt"
	"html"
	"strings"
	"time"
)

type FactureSummary struct {
	ID              int     `json:"id"`
	TransactionID   int     `json:"id_transaction"`
	NumeroFacture   string  `json:"numero_facture"`
	TypePaiement    string  `json:"type_paiement"`
	CommandeID      int     `json:"commande_id"`
	MontantTotal    float64 `json:"montant_total"`
	StatutPaiement  string  `json:"statut_paiement"`
	DateTransaction string  `json:"date_transaction"`
	AcheteurEmail   string  `json:"acheteur_email,omitempty"`
	AcheteurNom     string  `json:"acheteur_nom,omitempty"`
}

type FactureLine struct {
	TypeItem      string  `json:"type_item"`
	ReferenceID   int     `json:"reference_id"`
	PrixUnitaire  float64 `json:"prix_unitaire"`
	CommissionUPC float64 `json:"commission_upc"`
}

func CreateFactureForTransaction(tx *sql.Tx, transactionID int64) (int, string, error) {
	numero := fmt.Sprintf("FAC-%s-%06d", time.Now().Format("20060102"), transactionID)
	res, err := tx.Exec(
		"INSERT INTO FACTURE (id_transaction, numero_facture, type_paiement) VALUES (?, ?, 'Carte')",
		transactionID,
		numero,
	)
	if err != nil {
		return 0, "", err
	}

	factureID, err := res.LastInsertId()
	if err != nil {
		return 0, "", err
	}

	return int(factureID), numero, nil
}

func GetFacturesByUser(userID int) ([]FactureSummary, error) {
	rows, err := Conn.Query(
		"SELECT f.id, f.id_transaction, f.numero_facture, f.type_paiement, "+
			"t.id_commande, t.montant_total, t.statut_paiement, "+
			"DATE_FORMAT(t.date_transaction, '%Y-%m-%d %H:%i:%s'), "+
			"COALESCE(u.mail, ''), TRIM(CONCAT(COALESCE(u.prenom, ''), ' ', COALESCE(u.nom, ''))) "+
			"FROM FACTURE f "+
			"JOIN `TRANSACTION` t ON t.id = f.id_transaction "+
			"JOIN UTILISATEUR u ON u.id = t.id_acheteur "+
			"WHERE t.id_acheteur = ? "+
			"ORDER BY t.date_transaction DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	factures := []FactureSummary{}
	for rows.Next() {
		var facture FactureSummary
		if err := rows.Scan(
			&facture.ID,
			&facture.TransactionID,
			&facture.NumeroFacture,
			&facture.TypePaiement,
			&facture.CommandeID,
			&facture.MontantTotal,
			&facture.StatutPaiement,
			&facture.DateTransaction,
			&facture.AcheteurEmail,
			&facture.AcheteurNom,
		); err != nil {
			return nil, err
		}
		factures = append(factures, facture)
	}

	return factures, rows.Err()
}

func GetFactureByIDForUser(userID int, factureID int) (*FactureSummary, error) {
	var facture FactureSummary
	err := Conn.QueryRow(
		"SELECT f.id, f.id_transaction, f.numero_facture, f.type_paiement, "+
			"t.id_commande, t.montant_total, t.statut_paiement, "+
			"DATE_FORMAT(t.date_transaction, '%Y-%m-%d %H:%i:%s'), "+
			"COALESCE(u.mail, ''), TRIM(CONCAT(COALESCE(u.prenom, ''), ' ', COALESCE(u.nom, ''))) "+
			"FROM FACTURE f "+
			"JOIN `TRANSACTION` t ON t.id = f.id_transaction "+
			"JOIN UTILISATEUR u ON u.id = t.id_acheteur "+
			"WHERE f.id = ? AND t.id_acheteur = ?",
		factureID,
		userID,
	).Scan(
		&facture.ID,
		&facture.TransactionID,
		&facture.NumeroFacture,
		&facture.TypePaiement,
		&facture.CommandeID,
		&facture.MontantTotal,
		&facture.StatutPaiement,
		&facture.DateTransaction,
		&facture.AcheteurEmail,
		&facture.AcheteurNom,
	)
	if err != nil {
		return nil, err
	}

	return &facture, nil
}

func GetFactureLines(factureID int, userID int) ([]FactureLine, error) {
	rows, err := Conn.Query(
		"SELECT lc.type_item, lc.reference_id, lc.prix_unitaire, lc.commission_upc "+
			"FROM LIGNE_COMMANDE lc "+
			"JOIN `TRANSACTION` t ON t.id_commande = lc.id_commande "+
			"JOIN FACTURE f ON f.id_transaction = t.id "+
			"WHERE f.id = ? AND t.id_acheteur = ? "+
			"ORDER BY lc.id",
		factureID,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines := []FactureLine{}
	for rows.Next() {
		var line FactureLine
		if err := rows.Scan(&line.TypeItem, &line.ReferenceID, &line.PrixUnitaire, &line.CommissionUPC); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, rows.Err()
}

func BuildFactureHTML(userID int, factureID int) (string, *FactureSummary, error) {
	facture, err := GetFactureByIDForUser(userID, factureID)
	if err != nil {
		return "", nil, err
	}

	lines, err := GetFactureLines(factureID, userID)
	if err != nil {
		return "", nil, err
	}

	var rows strings.Builder
	for _, line := range lines {
		rows.WriteString(fmt.Sprintf(
			"<tr><td>%s #%d</td><td>%.2f EUR</td><td>%.2f EUR</td></tr>",
			html.EscapeString(line.TypeItem),
			line.ReferenceID,
			line.PrixUnitaire,
			line.CommissionUPC,
		))
	}

	document := fmt.Sprintf(`<!doctype html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <title>Facture %s</title>
  <style>
    body { font-family: Arial, sans-serif; color: #16221c; margin: 40px; }
    .invoice { max-width: 820px; margin: auto; border: 1px solid #d8e5dc; border-radius: 20px; padding: 32px; }
    .brand { color: #2f8f58; font-size: 14px; font-weight: 800; letter-spacing: 4px; text-transform: uppercase; }
    h1 { font-size: 42px; margin: 12px 0 8px; }
    .meta, .client { color: #64736a; line-height: 1.6; }
    table { width: 100%%; border-collapse: collapse; margin-top: 28px; }
    th, td { border-bottom: 1px solid #e6eee9; padding: 14px; text-align: left; }
    th { color: #2f8f58; font-size: 12px; text-transform: uppercase; letter-spacing: 2px; }
    .total { text-align: right; font-size: 24px; font-weight: 800; margin-top: 28px; }
  </style>
</head>
<body>
  <section class="invoice">
    <div class="brand">Upcycle Connect</div>
    <h1>Facture</h1>
    <p class="meta">
      Numero : <strong>%s</strong><br>
      Date : %s<br>
      Paiement : %s - %s<br>
      Commande : #%d
    </p>
    <p class="client">
      Client : <strong>%s</strong><br>
      Email : %s
    </p>
    <table>
      <thead><tr><th>Article</th><th>Prix</th><th>Commission UPC</th></tr></thead>
      <tbody>%s</tbody>
    </table>
    <div class="total">Total TTC : %.2f EUR</div>
  </section>
</body>
</html>`,
		html.EscapeString(facture.NumeroFacture),
		html.EscapeString(facture.NumeroFacture),
		html.EscapeString(facture.DateTransaction),
		html.EscapeString(facture.TypePaiement),
		html.EscapeString(facture.StatutPaiement),
		facture.CommandeID,
		html.EscapeString(facture.AcheteurNom),
		html.EscapeString(facture.AcheteurEmail),
		rows.String(),
		facture.MontantTotal,
	)

	return document, facture, nil
}
