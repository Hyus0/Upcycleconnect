package app

import (
	"upcycle/db"
	"upcycle/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("school", "esgi")

	users, err := db.GetAllUsers()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur de récupération des utilisateurs", http.StatusInternalServerError)
	}

	res, _ := json.Marshal(users)
	fmt.Fprintf(w, "%s", string(res))
}