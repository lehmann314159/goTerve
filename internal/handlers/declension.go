package handlers

import (
	"net/http"
	"strconv"

	"github.com/mikebway/goTerve/internal/models"
)

// DeclensionResult represents the result of a declension
type DeclensionResult struct {
	Noun           string            `json:"noun"`
	Translation    string            `json:"translation"`
	DeclensionType int               `json:"declensionType"`
	Cases          map[string]string `json:"cases"`
	Error          string            `json:"error,omitempty"`
}

// Decline handles noun declension requests
func (h *Handlers) Decline(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderPartial(w, "declension-result.html", DeclensionResult{Error: "Invalid form data"})
		return
	}

	nounNominative := r.FormValue("noun")
	if nounNominative == "" {
		// Get a random noun if none specified
		noun, err := h.db.GetRandomNoun()
		if err != nil {
			h.renderPartial(w, "declension-result.html", DeclensionResult{Error: "Failed to get noun"})
			return
		}
		nounNominative = noun.Nominative
	}

	// Get noun from database
	noun, err := h.db.GetNounByNominative(nounNominative)
	if err != nil {
		h.renderPartial(w, "declension-result.html", DeclensionResult{Error: "Database error"})
		return
	}
	if noun == nil {
		h.renderPartial(w, "declension-result.html", DeclensionResult{Error: "Noun not found"})
		return
	}

	// Parse requested case (optional - if not provided, return all cases)
	caseStr := r.FormValue("case")
	requestedCase := 0
	if caseStr != "" {
		requestedCase, _ = strconv.Atoi(caseStr)
	}

	// Get declensions
	cases := make(map[string]string)
	if requestedCase > 0 {
		// Single case requested
		nounCase := models.NounCase(requestedCase)
		declined := h.finnish.Decline(noun.Nominative, noun.DeclensionType, nounCase)
		cases[nounCase.Name()] = declined
	} else {
		// All cases
		allCases := []models.NounCase{
			models.CaseNominative,
			models.CaseGenitive,
			models.CasePartitive,
			models.CaseInessive,
			models.CaseElative,
			models.CaseIllative,
			models.CaseAdessive,
			models.CaseAblative,
			models.CaseAllative,
		}
		for _, c := range allCases {
			declined := h.finnish.Decline(noun.Nominative, noun.DeclensionType, c)
			cases[c.Name()] = declined
		}
	}

	result := DeclensionResult{
		Noun:           noun.Nominative,
		Translation:    noun.Translation,
		DeclensionType: int(noun.DeclensionType),
		Cases:          cases,
	}

	h.renderPartial(w, "declension-result.html", result)
}