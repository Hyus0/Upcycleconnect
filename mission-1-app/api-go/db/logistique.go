package db

import (
	"fmt"
	"database/sql"
	"upcycleconnect/api-go/models"
	"math/rand"
	"time"
)

func ReserverUnCasier(annonceID int, siteID int) (string, error) {
    var poidsAnnonce float64
    err := Conn.QueryRow("SELECT poids_estime_kg FROM ANNONCE WHERE id = ?", annonceID).Scan(&poidsAnnonce)
    if err != nil {
        return "", fmt.Errorf("impossible de trouver le poids de l'annonce")
    }

    query := `
        SELECT c.id FROM CASIER c
        JOIN CONTENEUR co ON c.id_conteneur = co.id
        WHERE co.id_site = ? 
        AND c.statut = 'Libre' 
        AND co.statut = 'Operationnel'
        AND (co.niveau_remplissage + ?) <= co.capacite_max_kg
        LIMIT 1`

    var casierID int
    err = Conn.QueryRow(query, siteID, poidsAnnonce).Scan(&casierID)
    if err != nil {
        return "", fmt.Errorf("Site complet ou limite de poids atteinte")
    }

    rand.Seed(time.Now().UnixNano())
    pin := fmt.Sprintf("%06d", rand.Intn(1000000))

    _, err = Conn.Exec("UPDATE CASIER SET statut = 'Reserve' WHERE id = ?", casierID)
    if err != nil {
        return "", err
    }

    _, err = Conn.Exec(`
        UPDATE ANNONCE SET 
            id_casier = ?, 
            id_site = ?, 
            code_pin_depot = ?, 
            statut = 'Reserve' 
        WHERE id = ?`, casierID, siteID, pin, annonceID)

    return pin, err
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