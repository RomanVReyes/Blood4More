package handlers

import (
	"net/http"
	"time"

	"roman-sangre/internal/models"
	"roman-sangre/internal/repository"
	"roman-sangre/pkg/utils"
)

const SessionCookieName = "roman_session"

func SetSessionCookie(w http.ResponseWriter, sessionID string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func DeleteSessionCookie(w http.ResponseWriter) {

	http.SetCookie(w, &http.Cookie{
		Name:   SessionCookieName,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

}

func StartUserSession(w http.ResponseWriter, userID string, r *http.Request) error {
	sessionID := utils.GenerateSessionID()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	sesion := models.Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsActive:  true,
		IP:        r.RemoteAddr,
		UserAgent: r.UserAgent(),
	}

	err := repository.CreateSession(sesion)
	if err != nil {
		return err
	}

	SetSessionCookie(w, sessionID, expiresAt)

	return nil
}

func DeactivateUserSession(w http.ResponseWriter, sessionID string) error {

	err := repository.DeactivateSession(sessionID)
	if err != nil {
		return err
	}

	DeleteSessionCookie(w)

	return nil
}

func GetSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
