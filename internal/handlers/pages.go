package handlers

import (
	"net/http"
)

// PageData contains common data for page templates
type PageData struct {
	Title       string
	ActivePage  string
	User        interface{}
	Verbs       interface{}
	Nouns       interface{}
	Words       interface{}
	CEFRLevels  []string
	Error       string
}

// Home renders the home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Terve - Finnish Language Learning",
		ActivePage: "home",
		CEFRLevels: []string{"A1", "A2", "B1", "B2", "C1", "C2"},
	}
	h.render(w, "home.html", data)
}

// ConjugationPage renders the verb conjugation practice page
func (h *Handlers) ConjugationPage(w http.ResponseWriter, r *http.Request) {
	verbs, err := h.db.GetVerbs()
	if err != nil {
		http.Error(w, "Failed to load verbs", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:      "Verb Conjugation - Terve",
		ActivePage: "conjugation",
		Verbs:      verbs,
		CEFRLevels: []string{"A1", "A2", "B1", "B2", "C1", "C2"},
	}
	h.render(w, "conjugation.html", data)
}

// DeclensionPage renders the noun declension practice page
func (h *Handlers) DeclensionPage(w http.ResponseWriter, r *http.Request) {
	nouns, err := h.db.GetNouns()
	if err != nil {
		http.Error(w, "Failed to load nouns", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:      "Noun Declension - Terve",
		ActivePage: "declension",
		Nouns:      nouns,
		CEFRLevels: []string{"A1", "A2", "B1", "B2", "C1", "C2"},
	}
	h.render(w, "declension.html", data)
}

// FlashcardsPage renders the flashcards page
func (h *Handlers) FlashcardsPage(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Flashcards - Terve",
		ActivePage: "flashcards",
		CEFRLevels: []string{"A1", "A2", "B1", "B2", "C1", "C2"},
	}
	h.render(w, "flashcards.html", data)
}

// ReadingPage renders the reading comprehension page
func (h *Handlers) ReadingPage(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Reading Practice - Terve",
		ActivePage: "reading",
		CEFRLevels: []string{"A1", "A2", "B1", "B2", "C1", "C2"},
	}
	h.render(w, "reading.html", data)
}

// TestPage renders the grammar test page
func (h *Handlers) TestPage(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Grammar Test - Terve",
		ActivePage: "test",
		CEFRLevels: []string{"A1", "A2", "B1", "B2", "C1", "C2"},
	}
	h.render(w, "test.html", data)
}

// LoginPage renders the login page
func (h *Handlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Login - Terve",
		ActivePage: "login",
	}
	h.render(w, "login.html", data)
}