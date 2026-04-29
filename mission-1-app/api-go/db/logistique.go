package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"upcycleconnect/api-go/models"
	"time"
)

func ReserverUnCasier(annonceID int, siteID int) (string, error) {
    query := `
        SELECT id
        FROM CONTENEUR
        WHERE id_site = ?
        AND statut = 'Operationnel'
        LIMIT 1`

    var conteneurID int
    err := Conn.QueryRow(query, siteID).Scan(&conteneurID)
    if err != nil {
        return "", fmt.Errorf("aucun conteneur operationnel disponible sur ce site")
    }

    rand.Seed(time.Now().UnixNano())
    pin := fmt.Sprintf("%06d", rand.Intn(1000000))

    tx, err := Conn.Begin()
    if err != nil {
        return "", err
    }

    _, err = tx.Exec(`
        INSERT INTO DEPOT_CONTENEUR (id_utilisateur, id_conteneur, id_objet, statut)
        SELECT id_vendeur, ?, id, 'Prevu'
        FROM ANNONCE
        WHERE id = ?
    `, conteneurID, annonceID)
    if err != nil {
        tx.Rollback()
        return "", err
    }

    _, err = tx.Exec(`UPDATE ANNONCE SET statut = 'Reserve' WHERE id = ?`, annonceID)
    if err != nil {
        tx.Rollback()
        return "", err
    }

    if err := tx.Commit(); err != nil {
        return "", err
    }

    return pin, nil
}

func GetAllSites() ([]models.Site, error) {
    if Conn == nil {
        return nil, fmt.Errorf("DB non connectée")
    }

    query := `SELECT id, nom, ville, code_postal, adresse, COALESCE(telephone, ''), type, actif FROM SITE`
    
    rows, err := Conn.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    sites := []models.Site{}

    for rows.Next() {
        var s models.Site
        err := rows.Scan(
            &s.ID, 
            &s.Nom, 
            &s.Ville, 
            &s.CodePostal, 
            &s.Adresse, 
            &s.Telephone, 
            &s.Type, 
            &s.Actif,
        )
        if err != nil {
            fmt.Println("Erreur lors du Scan du site:", err)
            return nil, err
        }
        sites = append(sites, s)
    }

    return sites, nil
}

func GetSiteByID(id int) (map[string]string, error) {
    query := `SELECT nom, ville, code_postal, adresse, COALESCE(telephone, "") FROM SITE WHERE id = ?`
    
    var nom, ville, codePostal, adresse, telephone string
    err := Conn.QueryRow(query, id).Scan(&nom, &ville, &codePostal, &adresse, &telephone)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("site introuvable")
        }
        return nil, err
    }

    return map[string]string{
        "nom":         nom,
        "ville":       ville,
        "code_postal": codePostal,
        "adresse":     adresse,
        "telephone":   telephone,
    }, nil
}

func GetConteneursBySite(siteID int) ([]models.Conteneur, error) {
	query := `SELECT id, id_site, type_dechet, statut, capacite_max_kg, niveau_remplissage 
	          FROM CONTENEUR WHERE id_site = ?`
	rows, err := Conn.Query(query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conteneurs []models.Conteneur
	for rows.Next() {
		var c models.Conteneur
		if err := rows.Scan(&c.ID, &c.IdSite, &c.TypeDechet, &c.Statut, &c.CapaciteMaxKg, &c.NiveauRemplissage); err != nil {
			return nil, err
		}
		conteneurs = append(conteneurs, c)
	}
	return conteneurs, nil
}

func RetireObjetCasier(idAnnonce int) error {
    tx, err := Conn.Begin()
    if err != nil {
        return err
    }

    _, err = tx.Exec(`
        UPDATE DEPOT_CONTENEUR
        SET statut = 'Annule', date_recuperation = NOW()
        WHERE id_objet = ? AND statut IN ('Prevu', 'Effectue')
    `, idAnnonce)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`UPDATE ANNONCE SET statut = 'Disponible' WHERE id = ?`, idAnnonce)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}
