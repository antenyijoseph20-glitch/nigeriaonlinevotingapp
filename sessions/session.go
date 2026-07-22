package sessions

import (
	"net/http"
	"strconv"
)

const SessionCookieName = "session"

// CreateSession creates a session cookie using the user's ID.
func CreateSession(w http.ResponseWriter, userID int) {

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    strconv.Itoa(userID),
		Path:     "/",
		HttpOnly: true,
	})
}

// DeleteSession removes the session cookie.
func DeleteSession(w http.ResponseWriter) {

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// GetSession returns the session cookie.
func GetSession(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(SessionCookieName)
}

// GetSessionUserID returns the logged-in user's ID.
func GetSessionUserID(r *http.Request) (int, error) {

	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(cookie.Value)
}
