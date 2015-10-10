package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var CookieStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(16),
	securecookie.GenerateRandomKey(16))

func ServeAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" || !validateReferer(r) {
		w.Write([]byte(`{"error": "invalid request"}`))
		return
	}

	password := r.PostFormValue("password")
	hash := HashPassword(password)
	if hash == GlobalDb.Config().PasswordHash {
		s, _ := CookieStore.Get(r, "sessid")
		s.Values["authenticated"] = true
		s.Save(r, w)
		w.Write([]byte(`{}`))
		return
	}
	w.Write([]byte(`{"error": "invalid login credentials"}`))
}

// IsAuthenticated returns whether or not a request was authenticated.
func IsAuthenticated(r *http.Request) bool {
	s, _ := CookieStore.Get(r, "sessid")
	val, ok := s.Values["authenticated"].(bool)
	return ok && val
}

// HashPassword returns the hash of a password string.
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

// validateReferer makes sure the Referer's host is the same as the current host.
func validateReferer(r *http.Request) bool {
	host := r.Host
	if forwardHost := r.Header.Get("X-Forwarded-Host"); forwardHost != "" {
		parts := strings.Split(forwardHost, ",")
		host = strings.TrimSpace(parts[len(parts)-1])
	}
	referer := r.Referer()
	u, err := url.Parse(referer)
	return err != nil && u.Host == host
}
