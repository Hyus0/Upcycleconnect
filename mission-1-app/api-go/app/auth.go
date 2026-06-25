package app

import (
	"net/http"
	"strconv"
	"strings"

	"upcycleconnect/api-go/db"
)

// AuthUser porte l'identite resolue a partir du token de session.
type AuthUser struct {
	ID   int
	Role string
}

// extractToken lit le token depuis l'en-tete Authorization.
// Supporte les formes "Bearer <token>" et "<token>".
func extractToken(r *http.Request) string {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	if h == "" {
		return ""
	}
	if len(h) >= 7 && strings.EqualFold(h[:7], "bearer ") {
		return strings.TrimSpace(h[7:])
	}
	return h
}

// authenticate resout l'utilisateur courant a partir du token.
func authenticate(r *http.Request) (*AuthUser, bool) {
	token := extractToken(r)
	if token == "" {
		return nil, false
	}
	id, role, err := db.GetUserByToken(token)
	if err != nil || id == 0 {
		return nil, false
	}
	return &AuthUser{ID: id, Role: role}, true
}

// RequireRole renvoie un middleware exigeant un token valide et, si des roles
// sont fournis, un role autorise. Sans role fourni, exige seulement un token valide.
//   - 401 unauthorized : token absent ou invalide
//   - 403 forbidden    : token valide mais role non autorise
func RequireRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := authenticate(r)
			if !ok {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			if len(roles) > 0 {
				allowed := false
				for _, role := range roles {
					if strings.EqualFold(user.Role, role) {
						allowed = true
						break
					}
				}
				if !allowed {
					writeError(w, http.StatusForbidden, "forbidden")
					return
				}
			}
			next(w, r)
		}
	}
}

// RequireAuth exige uniquement un token valide (tout role).
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return RequireRole()(next)
}

// RequireSelf exige un token valide ET que l'utilisateur authentifie soit le
// proprietaire de la ressource (param de chemin == id du token). Les Admin
// sont autorises sur toutes les ressources.
//   - 401 unauthorized : token absent ou invalide
//   - 400 invalid_id   : param de chemin non numerique
//   - 403 forbidden    : l'utilisateur tente d'acceder a une ressource d'autrui
func RequireSelf(param string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := authenticate(r)
			if !ok {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			if strings.EqualFold(user.Role, "Admin") {
				next(w, r)
				return
			}
			targetID, err := strconv.Atoi(strings.TrimSpace(r.PathValue(param)))
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid_id")
				return
			}
			if targetID != user.ID {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
			next(w, r)
		}
	}
}
