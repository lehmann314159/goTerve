package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mikebway/goTerve/internal/models"
)

// FlashcardData contains data for flashcard rendering
type FlashcardData struct {
	Word      *models.Word          `json:"word"`
	Flashcard *models.UserFlashcard `json:"flashcard,omitempty"`
	IsNew     bool                  `json:"isNew"`
	Error     string                `json:"error,omitempty"`
}

// FlashcardAnswerResult contains the result of answering a flashcard
type FlashcardAnswerResult struct {
	Correct     bool   `json:"correct"`
	CorrectWord string `json:"correctWord"`
	Category    string `json:"category"`
	Message     string `json:"message"`
	Error       string `json:"error,omitempty"`
}

// GetFlashcard returns the next flashcard for the user
func (h *Handlers) GetFlashcard(w http.ResponseWriter, r *http.Request) {
	// For now, use a guest user ID - in production this would come from session
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "guest"
	}

	cefrLevel := r.URL.Query().Get("level")
	if cefrLevel == "" {
		cefrLevel = "A1"
	}

	word, flashcard, err := h.db.GetNextFlashcard(userID, cefrLevel)
	if err != nil {
		h.renderPartial(w, "flashcard.html", FlashcardData{Error: "Failed to get flashcard"})
		return
	}

	if word == nil {
		// No words available - get a random one
		word, err = h.db.GetRandomWord(cefrLevel)
		if err != nil {
			h.renderPartial(w, "flashcard.html", FlashcardData{Error: "No words available"})
			return
		}
	}

	data := FlashcardData{
		Word:      word,
		Flashcard: flashcard,
		IsNew:     flashcard == nil,
	}

	h.renderPartial(w, "flashcard.html", data)
}

// AnswerFlashcard processes a flashcard answer
func (h *Handlers) AnswerFlashcard(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderPartial(w, "flashcard-result.html", FlashcardAnswerResult{Error: "Invalid form data"})
		return
	}

	userID := r.FormValue("user_id")
	if userID == "" {
		userID = "guest"
	}

	wordIDStr := r.FormValue("word_id")
	wordID, err := strconv.Atoi(wordIDStr)
	if err != nil {
		h.renderPartial(w, "flashcard-result.html", FlashcardAnswerResult{Error: "Invalid word ID"})
		return
	}

	answer := r.FormValue("answer")

	// Get the word to check the answer
	word, err := h.db.GetWordByID(wordID)
	if err != nil {
		h.renderPartial(w, "flashcard-result.html", FlashcardAnswerResult{Error: "Failed to get word"})
		return
	}

	// Check if answer is correct (case-insensitive)
	correct := normalizeAnswer(answer) == normalizeAnswer(word.English)

	// Get or create flashcard record
	flashcard, err := h.db.GetUserFlashcard(userID, wordID)
	if err != nil {
		h.renderPartial(w, "flashcard-result.html", FlashcardAnswerResult{Error: "Database error"})
		return
	}

	now := time.Now()
	if flashcard == nil {
		flashcard = &models.UserFlashcard{
			UserID:   userID,
			WordID:   wordID,
			Category: "new",
		}
	}

	flashcard.TimesReviewed++
	flashcard.LastReviewedAt = &now

	if correct {
		flashcard.TimesCorrect++
		// Update category based on success rate
		flashcard.Category = calculateCategory(flashcard)
		// Schedule next review based on spaced repetition
		nextReview := calculateNextReview(flashcard)
		flashcard.NextReviewAt = &nextReview
	} else {
		// Wrong answer - reset to learning
		if flashcard.Category == "mastered" || flashcard.Category == "review" {
			flashcard.Category = "learning"
		}
		// Review sooner
		nextReview := now.Add(5 * time.Minute)
		flashcard.NextReviewAt = &nextReview
	}

	if err := h.db.CreateOrUpdateFlashcard(flashcard); err != nil {
		h.renderPartial(w, "flashcard-result.html", FlashcardAnswerResult{Error: "Failed to save progress"})
		return
	}

	result := FlashcardAnswerResult{
		Correct:     correct,
		CorrectWord: word.English,
		Category:    flashcard.Category,
	}

	if correct {
		result.Message = "Correct! Great job!"
	} else {
		result.Message = "Not quite. The answer is: " + word.English
	}

	h.renderPartial(w, "flashcard-result.html", result)
}

// normalizeAnswer normalizes an answer for comparison
func normalizeAnswer(s string) string {
	// Simple normalization - lowercase and trim
	return strings.TrimSpace(strings.ToLower(s))
}

// calculateCategory determines the flashcard category based on performance
func calculateCategory(fc *models.UserFlashcard) string {
	if fc.TimesReviewed == 0 {
		return "new"
	}

	successRate := float64(fc.TimesCorrect) / float64(fc.TimesReviewed)

	if fc.TimesReviewed >= 10 && successRate >= 0.9 {
		return "mastered"
	}
	if fc.TimesReviewed >= 5 && successRate >= 0.7 {
		return "review"
	}
	return "learning"
}

// calculateNextReview calculates the next review time based on spaced repetition
func calculateNextReview(fc *models.UserFlashcard) time.Time {
	now := time.Now()

	// Simple spaced repetition intervals
	switch fc.Category {
	case "new":
		return now.Add(1 * time.Hour)
	case "learning":
		return now.Add(4 * time.Hour)
	case "review":
		return now.Add(24 * time.Hour)
	case "mastered":
		return now.Add(7 * 24 * time.Hour)
	default:
		return now.Add(1 * time.Hour)
	}
}