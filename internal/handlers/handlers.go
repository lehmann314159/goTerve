package handlers

import (
	"html/template"
	"net/http"

	"github.com/mikebway/goTerve/internal/claude"
	"github.com/mikebway/goTerve/internal/finnish"
	"github.com/mikebway/goTerve/internal/store"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	db        *store.Store
	templates *template.Template
	finnish   *finnish.Finnish
	claude    *claude.Client
}

// New creates a new Handlers instance
func New(db *store.Store, templates *template.Template, anthropicKey string) *Handlers {
	return &Handlers{
		db:        db,
		templates: templates,
		finnish:   finnish.New(),
		claude:    claude.New(anthropicKey),
	}
}

// render is a helper to render templates
func (h *Handlers) render(w http.ResponseWriter, templateName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderPartial renders a partial template for HTMX responses
func (h *Handlers) renderPartial(w http.ResponseWriter, templateName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// jsonError sends a JSON error response
func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"error":"` + message + `"}`))
}