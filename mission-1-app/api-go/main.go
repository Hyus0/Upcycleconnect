package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
)

type scoreRequest struct {
	Category      string  `json:"category"`
	WeightKg      float64 `json:"weightKg,string"`
	ReusedPercent float64 `json:"reusedPercent,string"`
}

type barcodeRequest struct {
	ContainerID string `json:"containerId"`
	Barcode     string `json:"barcode"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/upcycling-score", upcyclingScoreHandler)
	mux.HandleFunc("/api/barcode/validate", barcodeValidateHandler)
	mux.HandleFunc("/api/notifications/preview", notificationPreviewHandler)

	addr := ":8081"
	log.Printf("api-go listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "api-go"})
}

func upcyclingScoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	var req scoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
		return
	}

	if req.WeightKg <= 0 || req.ReusedPercent < 0 || req.ReusedPercent > 100 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_input"})
		return
	}

	categoryFactor := map[string]float64{
		"furniture":   1.2,
		"textile":     0.8,
		"electronics": 1.5,
		"mixed":       1.0,
	}[strings.ToLower(req.Category)]
	if categoryFactor == 0 {
		categoryFactor = 1.0
	}

	base := (req.WeightKg * req.ReusedPercent / 10.0) * categoryFactor
	score := math.Min(100, math.Round(base))
	co2 := math.Round((req.WeightKg*(req.ReusedPercent/100.0)*categoryFactor)*100) / 100

	writeJSON(w, http.StatusOK, map[string]any{
		"score": score,
		"co2Kg": co2,
	})
}

func barcodeValidateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	var req barcodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
		return
	}

	status := "rejected"
	if len(req.Barcode) >= 6 && strings.TrimSpace(req.ContainerID) != "" {
		status = "accepted"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"containerId": req.ContainerID,
		"barcode":     req.Barcode,
		"status":      status,
	})
}

func notificationPreviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"provider": "onesignal_pending_integration",
		"status":   "preview_ok",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
