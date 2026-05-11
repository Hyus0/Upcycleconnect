package db

import "fmt"

func GetFollowStatus(idProfil int, idConnecte int) (int, int, bool, error) {
	if Conn == nil {
		return 0, 0, false, fmt.Errorf("connexion DB non initialisee")
	}

	var totalFollowers int
	err := Conn.QueryRow("SELECT COUNT(*) FROM ABONNEMENT_UTILISATEUR WHERE id_suivi = ?", idProfil).Scan(&totalFollowers)
	if err != nil {
		return 0, 0, false, err
	}

	var totalFollowing int
	err = Conn.QueryRow("SELECT COUNT(*) FROM ABONNEMENT_UTILISATEUR WHERE id_abonne = ?", idProfil).Scan(&totalFollowing)
	if err != nil {
		return 0, 0, false, err
	}

	var isFollowing bool = false
	if idConnecte > 0 {
		var count int
		err = Conn.QueryRow("SELECT COUNT(*) FROM ABONNEMENT_UTILISATEUR WHERE id_suivi = ? AND id_abonne = ?", idProfil, idConnecte).Scan(&count)
		if err == nil {
			isFollowing = (count > 0)
		}
	}

	return totalFollowers, totalFollowing, isFollowing, nil
}

func ToggleFollowUser(idSuivi int, idAbonne int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	var count int
	err := Conn.QueryRow("SELECT COUNT(*) FROM ABONNEMENT_UTILISATEUR WHERE id_suivi = ? AND id_abonne = ?", idSuivi, idAbonne).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		_, err = Conn.Exec("DELETE FROM ABONNEMENT_UTILISATEUR WHERE id_suivi = ? AND id_abonne = ?", idSuivi, idAbonne)
	} else {
		_, err = Conn.Exec("INSERT INTO ABONNEMENT_UTILISATEUR (id_suivi, id_abonne) VALUES (?, ?)", idSuivi, idAbonne)
	}

	return err
}