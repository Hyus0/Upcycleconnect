package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
)

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	text := value.String
	return &text
}

func nullIntPtr(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	number := int(value.Int64)
	return &number
}

func nullTimePtr(value sql.NullTime) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.Time.Format("2006-01-02 15:04:05")
	return &formatted
}

type annonceRow struct {
	ID                        int
	IdVendeur                 int
	IdAcheteur                sql.NullInt64
	IdCasier                  sql.NullInt64
	IdCategorie               sql.NullInt64
	IdSite                    sql.NullInt64
	Titre                     sql.NullString
	Description               sql.NullString
	TypeMateriau              sql.NullString
	PoidsEstimeKg             sql.NullFloat64
	Prix                      sql.NullFloat64
	EtatObjet                 sql.NullString
	Statut                    sql.NullString
	EstValide                 sql.NullString
	CodePinDepot              sql.NullString
	CodeBarreRetrait          sql.NullString
	DateCreation              sql.NullTime
	DateDepotEffective        sql.NullTime
	DateRecuperationEffective sql.NullTime
	Type                      sql.NullString
	Ville                     sql.NullString
	CodePostal                sql.NullString
	Adresse                   sql.NullString
}

func scanAnnonce(row annonceRow) models.Annonce {
	return models.Annonce{
		ID:                        row.ID,
		IdVendeur:                 row.IdVendeur,
		IdAcheteur:                nullIntPtr(row.IdAcheteur),
		IdCasier:                  nullIntPtr(row.IdCasier),
		IdCategorie:               int(row.IdCategorie.Int64),
		Titre:                     row.Titre.String,
		Description:               row.Description.String,
		TypeMateriau:              row.TypeMateriau.String,
		PoidsEstimeKg:             row.PoidsEstimeKg.Float64,
		Prix:                      row.Prix.Float64,
		EtatObjet:                 row.EtatObjet.String,
		Statut:                    row.Statut.String,
		EstValide:                 row.EstValide.String,
		CodePinDepot:              nullStringPtr(row.CodePinDepot),
		IdSite:                    nullIntPtr(row.IdSite),
		CodeBarreRetrait:          nullStringPtr(row.CodeBarreRetrait),
		DateCreation:              nullTimePtr(row.DateCreation),
		DateDepotEffective:        nullTimePtr(row.DateDepotEffective),
		DateRecuperationEffective: nullTimePtr(row.DateRecuperationEffective),
		Type:                      row.Type.String,
		Ville:                     row.Ville.String,
		CodePostal:                row.CodePostal.String,
		Adresse:                   row.Adresse.String,
	}
}

func CreateAnnonce(a *models.Annonce) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		INSERT INTO ANNONCE (
			id_vendeur, titre, description, prix, etat_objet,
			statut, est_valide, type, ville, code_postal, adresse, provider
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := Conn.Exec(
		query,
		a.IdVendeur,
		a.Titre,
		a.Description,
		a.Prix,
		a.EtatObjet,
		"Disponible",
		"En attente",
		a.Type,
		a.Ville,
		a.CodePostal,
		a.Adresse,
		"",
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = int(id)
	return nil
}

func GetAllAnnonces() ([]models.Annonce, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		SELECT
			a.id,
			a.id_vendeur,
			a.id_acheteur,
			NULL AS id_casier,
			NULL AS id_categorie,
			s.id AS id_site,
			a.titre,
			a.description,
			'' AS type_materiau,
			0 AS poids_estime_kg,
			a.prix,
			a.etat_objet,
			a.statut,
			a.est_valide,
			CASE WHEN dc.id IS NULL THEN NULL ELSE CONCAT('DEPOT-', dc.id) END AS code_pin_depot,
			NULL AS code_barre_retrait,
			a.date_creation,
			dc.date_depot AS date_depot_effective,
			dc.date_recuperation AS date_recuperation_effective,
			a.type,
			a.ville,
			a.code_postal,
			a.adresse
		FROM ANNONCE a
		LEFT JOIN (
			SELECT d1.*
			FROM DEPOT_CONTENEUR d1
			INNER JOIN (
				SELECT id_objet, MAX(id) AS max_id
				FROM DEPOT_CONTENEUR
				GROUP BY id_objet
			) latest ON latest.id_objet = d1.id_objet AND latest.max_id = d1.id
		) dc ON dc.id_objet = a.id
		LEFT JOIN CONTENEUR c ON c.id = dc.id_conteneur
		LEFT JOIN SITE s ON s.id = c.id_site
	`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var annonces []models.Annonce
	for rows.Next() {
		var row annonceRow

		err := rows.Scan(
			&row.ID, &row.IdVendeur, &row.IdAcheteur, &row.IdCasier, &row.IdCategorie, &row.IdSite,
			&row.Titre, &row.Description, &row.TypeMateriau, &row.PoidsEstimeKg, &row.Prix,
			&row.EtatObjet, &row.Statut, &row.EstValide, &row.CodePinDepot, &row.CodeBarreRetrait,
			&row.DateCreation, &row.DateDepotEffective, &row.DateRecuperationEffective,
			&row.Type, &row.Ville, &row.CodePostal, &row.Adresse,
		)
		if err != nil {
			fmt.Println("Erreur Scan GetAllAnnonces:", err)
			return nil, err
		}
		annonces = append(annonces, scanAnnonce(row))
	}

	return annonces, nil
}

func GetAnnoncesByUserID(userID int) ([]models.Annonce, error) {
	query := `
        SELECT
            a.id,
            a.id_vendeur,
            a.id_acheteur,
            NULL AS id_casier,
            NULL AS id_categorie,
            s.id AS id_site,
            a.titre,
            a.description,
            '' AS type_materiau,
            0 AS poids_estime_kg,
            a.prix,
            a.etat_objet,
            a.statut,
            a.est_valide,
            CASE WHEN dc.id IS NULL THEN NULL ELSE CONCAT('DEPOT-', dc.id) END AS code_pin_depot,
            NULL AS code_barre_retrait,
            a.date_creation,
            dc.date_depot AS date_depot_effective,
            dc.date_recuperation AS date_recuperation_effective,
            a.type,
            a.ville,
            a.code_postal,
            a.adresse
        FROM ANNONCE a
        LEFT JOIN (
            SELECT d1.*
            FROM DEPOT_CONTENEUR d1
            INNER JOIN (
                SELECT id_objet, MAX(id) AS max_id
                FROM DEPOT_CONTENEUR
                GROUP BY id_objet
            ) latest ON latest.id_objet = d1.id_objet AND latest.max_id = d1.id
        ) dc ON dc.id_objet = a.id
        LEFT JOIN CONTENEUR c ON c.id = dc.id_conteneur
        LEFT JOIN SITE s ON s.id = c.id_site
        WHERE a.id_vendeur = ?
    `

	rows, err := Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var annonces []models.Annonce
	for rows.Next() {
		var row annonceRow
		err := rows.Scan(
			&row.ID, &row.IdVendeur, &row.IdAcheteur, &row.IdCasier, &row.IdCategorie, &row.IdSite,
			&row.Titre, &row.Description, &row.TypeMateriau, &row.PoidsEstimeKg,
			&row.Prix, &row.EtatObjet, &row.Statut, &row.EstValide,
			&row.CodePinDepot, &row.CodeBarreRetrait,
			&row.DateCreation, &row.DateDepotEffective, &row.DateRecuperationEffective,
			&row.Type, &row.Ville, &row.CodePostal, &row.Adresse,
		)
		if err != nil {
			fmt.Println("Erreur Scan GetAnnoncesByUserID:", err)
			return nil, err
		}
		annonces = append(annonces, scanAnnonce(row))
	}
	return annonces, nil
}

func GetAnnonce(id int) (*models.Annonce, error) {
	query := `
		SELECT
			a.id,
			a.id_vendeur,
			a.id_acheteur,
			NULL AS id_casier,
			NULL AS id_categorie,
			s.id AS id_site,
			a.titre,
			a.description,
			'' AS type_materiau,
			0 AS poids_estime_kg,
			a.prix,
			a.etat_objet,
			a.statut,
			a.est_valide,
			CASE WHEN dc.id IS NULL THEN NULL ELSE CONCAT('DEPOT-', dc.id) END AS code_pin_depot,
			NULL AS code_barre_retrait,
			a.date_creation,
			dc.date_depot AS date_depot_effective,
			dc.date_recuperation AS date_recuperation_effective,
			a.type,
			a.ville,
			a.code_postal,
			a.adresse
		FROM ANNONCE a
		LEFT JOIN (
			SELECT d1.*
			FROM DEPOT_CONTENEUR d1
			INNER JOIN (
				SELECT id_objet, MAX(id) AS max_id
				FROM DEPOT_CONTENEUR
				GROUP BY id_objet
			) latest ON latest.id_objet = d1.id_objet AND latest.max_id = d1.id
		) dc ON dc.id_objet = a.id
		LEFT JOIN CONTENEUR c ON c.id = dc.id_conteneur
		LEFT JOIN SITE s ON s.id = c.id_site
		WHERE a.id = ?
	`

	var row annonceRow
	err := Conn.QueryRow(query, id).Scan(
		&row.ID, &row.IdVendeur, &row.IdAcheteur, &row.IdCasier, &row.IdCategorie, &row.IdSite,
		&row.Titre, &row.Description, &row.TypeMateriau, &row.PoidsEstimeKg, &row.Prix,
		&row.EtatObjet, &row.Statut, &row.EstValide, &row.CodePinDepot, &row.CodeBarreRetrait,
		&row.DateCreation, &row.DateDepotEffective, &row.DateRecuperationEffective,
		&row.Type, &row.Ville, &row.CodePostal, &row.Adresse,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	a := scanAnnonce(row)
	return &a, nil
}

func ModifyAnnonce(id int, a models.Annonce) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		UPDATE ANNONCE SET
			id_vendeur = ?,
			id_acheteur = ?,
			titre = ?,
			description = ?,
			statut = ?,
			est_valide = ?,
			prix = ?,
			etat_objet = ?,
			adresse = ?,
			ville = ?,
			code_postal = ?,
			type = ?,
			provider = ?
		WHERE id = ?
	`

	result, err := Conn.Exec(
		query,
		a.IdVendeur,
		a.IdAcheteur,
		a.Titre,
		a.Description,
		a.Statut,
		a.EstValide,
		a.Prix,
		a.EtatObjet,
		a.Adresse,
		a.Ville,
		a.CodePostal,
		a.Type,
		"",
		id,
	)
	if err != nil {
		return fmt.Errorf("ModifyAnnonce: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ModifyAnnonce RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune annonce trouvee avec l'ID %d", id)
	}
	return nil
}

func DeleteAnnonce(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	result, err := Conn.Exec("DELETE FROM ANNONCE WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteAnnonce: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteAnnonce RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune annonce trouvee avec l'ID %d", id)
	}
	return nil
}
