package server

import (
	"net/http"
)

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Pages
	s.router.Get("/", s.handlers.Home)
	s.router.Get("/conjugation", s.handlers.ConjugationPage)
	s.router.Get("/declension", s.handlers.DeclensionPage)
	s.router.Get("/flashcards", s.handlers.FlashcardsPage)
	s.router.Get("/reading", s.handlers.ReadingPage)
	s.router.Get("/test", s.handlers.TestPage)
	s.router.Get("/login", s.handlers.LoginPage)

	// HTMX API endpoints
	s.router.Post("/api/conjugate", s.handlers.Conjugate)
	s.router.Post("/api/decline", s.handlers.Decline)
	s.router.Get("/api/flashcard", s.handlers.GetFlashcard)
	s.router.Post("/api/flashcard/answer", s.handlers.AnswerFlashcard)
	s.router.Post("/api/story/generate", s.handlers.GenerateStory)
	s.router.Get("/api/test/question", s.handlers.GetTestQuestion)
	s.router.Post("/api/test/answer", s.handlers.AnswerTestQuestion)

	// Auth routes
	s.router.Post("/auth/login", s.handlers.Login)
	s.router.Post("/auth/logout", s.handlers.Logout)
}
