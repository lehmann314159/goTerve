package server

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikebway/goTerve/internal/handlers"
	"github.com/mikebway/goTerve/internal/store"
)

// Server represents the HTTP server
type Server struct {
	db        *store.Store
	port      string
	router    *chi.Mux
	templates *template.Template
	handlers  *handlers.Handlers
}

// New creates a new server instance
func New(db *store.Store, port string) *Server {
	s := &Server{
		db:     db,
		port:   port,
		router: chi.NewRouter(),
	}

	// Parse templates
	s.parseTemplates()

	// Create handlers
	s.handlers = handlers.New(db, s.templates)

	// Setup middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))

	// Setup routes
	s.setupRoutes()

	return s
}

// parseTemplates loads all HTML templates
func (s *Server) parseTemplates() {
	// Define template functions
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"add": func(a, b int) int {
			return a + b
		},
	}

	// Parse all templates
	tmpl := template.New("").Funcs(funcMap)

	// Parse layout first
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("templates", "layouts", "*.html")))

	// Parse pages
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("templates", "pages", "*.html")))

	// Parse partials
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join("templates", "partials", "*.html")))

	s.templates = tmpl
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return http.ListenAndServe(":"+s.port, s.router)
}
