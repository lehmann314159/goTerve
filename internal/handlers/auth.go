package handlers

import (
	"net/http"
)

// Login handles user login
// Note: This is a placeholder - full implementation would use OAuth
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// For now, just redirect to home
	// In production, this would:
	// 1. Validate OAuth tokens
	// 2. Create/update user in database
	// 3. Create session
	// 4. Set session cookie

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout handles user logout
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}