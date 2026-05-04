package handlers

import (
	"context"
	"net/http"
	"roman-sangre/internal/repository"
	"time"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionID, err := GetSessionID(r)
		if err != nil {
			http.Redirect(w, r, "/donante/auth", http.StatusSeeOther)
			return
		}

		sesion, err := repository.GetSessionByID(sessionID)
		if err != nil {
			http.Redirect(w, r, "/donante/auth", http.StatusSeeOther)
			return
		}

		if !sesion.IsActive {
			http.Redirect(w, r, "/donante/auth", http.StatusSeeOther)
			return
		}

		// validar expiración
		if time.Now().After(sesion.ExpiresAt) {
			http.Redirect(w, r, "/donante/auth", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", sesion.UserID)

		next(w, r.WithContext(ctx))
	}
}

// Luego separar los middlewares para cada tipo de userrr
