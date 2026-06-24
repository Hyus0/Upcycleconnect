package app

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

// CreatePaymentIntentHandler crée un PaymentIntent Stripe et retourne son client_secret
// au frontend pour qu'il puisse confirmer le paiement via Stripe.js.
// POST /paiement/intent
func CreatePaymentIntentHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Montant   int64  `json:"montant"`    // en centimes (ex: 1999 = 19,99 €)
		Mode      string `json:"mode"`       // "panier" | "annonce" | "subscription"
		PlanID    int    `json:"plan_id"`
		AnnonceID int    `json:"annonce_id"`
		UserID    int    `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	if req.Montant <= 0 {
		http.Error(w, "Montant invalide", http.StatusBadRequest)
		return
	}

	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		http.Error(w, "Stripe non configuré (STRIPE_SECRET_KEY manquante)", http.StatusInternalServerError)
		return
	}
	stripe.Key = secretKey

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Montant),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
		Metadata: map[string]string{
			"mode":       req.Mode,
			"user_id":    strconv.Itoa(req.UserID),
			"plan_id":    strconv.Itoa(req.PlanID),
			"annonce_id": strconv.Itoa(req.AnnonceID),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, "Erreur Stripe : "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"client_secret": pi.ClientSecret,
	})
}
