package handlers

import (
	"net/http"

	"github.com/mikebway/goTerve/internal/models"
)

// ConjugationRequest represents a conjugation request
type ConjugationRequest struct {
	Verb    string `json:"verb"`
	Tense   int    `json:"tense"`
	Person  int    `json:"person"`  // 1, 2, 3
	Plural  bool   `json:"plural"`
}

// ConjugationResult represents the result of a conjugation
type ConjugationResult struct {
	Verb         string   `json:"verb"`
	Tense        string   `json:"tense"`
	Conjugation  string   `json:"conjugation"`
	Person       string   `json:"person"`
	Translation  string   `json:"translation"`
	VerbType     int      `json:"verbType"`
	AllForms     []string `json:"allForms"`
	IsImperative bool     `json:"isImperative"`
	Error        string   `json:"error,omitempty"`
}

// Conjugate handles verb conjugation requests
func (h *Handlers) Conjugate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderPartial(w, "conjugation-result.html", ConjugationResult{Error: "Invalid form data"})
		return
	}

	verbInfinitive := r.FormValue("verb")
	if verbInfinitive == "" {
		// Get a random verb if none specified
		verb, err := h.db.GetRandomVerb()
		if err != nil {
			h.renderPartial(w, "conjugation-result.html", ConjugationResult{Error: "Failed to get verb"})
			return
		}
		verbInfinitive = verb.Infinitive
	}

	// Get verb from database
	verb, err := h.db.GetVerbByInfinitive(verbInfinitive)
	if err != nil {
		h.renderPartial(w, "conjugation-result.html", ConjugationResult{Error: "Database error"})
		return
	}
	if verb == nil {
		h.renderPartial(w, "conjugation-result.html", ConjugationResult{Error: "Verb not found"})
		return
	}

	// Parse tense
	tenseStr := r.FormValue("tense")
	tense := models.TensePresent
	switch tenseStr {
	case "imperfect":
		tense = models.TenseImperfect
	case "perfect":
		tense = models.TensePerfect
	case "conditional":
		tense = models.TenseConditional
	case "imperative":
		tense = models.TenseImperative
	}

	// Get all conjugations for this verb and tense
	allForms := h.finnish.ConjugateAll(verb.Infinitive, verb.Type, tense)

	result := ConjugationResult{
		Verb:         verb.Infinitive,
		Tense:        tense.Name(),
		Translation:  verb.Translation,
		VerbType:     int(verb.Type),
		AllForms:     allForms,
		IsImperative: tense == models.TenseImperative,
	}

	h.renderPartial(w, "conjugation-result.html", result)
}