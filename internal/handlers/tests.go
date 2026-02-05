package handlers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/mikebway/goTerve/internal/models"
)

// TestQuestionData contains data for a test question
type TestQuestionData struct {
	QuestionNum    int      `json:"questionNum"`
	TotalQuestions int      `json:"totalQuestions"`
	QuestionType   string   `json:"questionType"` // conjugation, declension, vocabulary
	Question       string   `json:"question"`
	Options        []string `json:"options"`
	CorrectIndex   int      `json:"correctIndex"`
	SessionID      int      `json:"sessionId"`
	Word           string   `json:"word"`
	WordID         int      `json:"wordId"`
	Error          string   `json:"error,omitempty"`
}

// TestAnswerResult contains the result of answering a test question
type TestAnswerResult struct {
	Correct        bool   `json:"correct"`
	CorrectAnswer  string `json:"correctAnswer"`
	Explanation    string `json:"explanation"`
	QuestionNum    int    `json:"questionNum"`
	TotalQuestions int    `json:"totalQuestions"`
	Score          int    `json:"score"`
	SessionID      int    `json:"sessionId"`
	IsComplete     bool   `json:"isComplete"`
	Error          string `json:"error,omitempty"`
}

// GetTestQuestion returns a test question
func (h *Handlers) GetTestQuestion(w http.ResponseWriter, r *http.Request) {
	testType := r.URL.Query().Get("type")
	if testType == "" {
		testType = "vocabulary"
	}

	cefrLevel := r.URL.Query().Get("level")
	if cefrLevel == "" {
		cefrLevel = "A1"
	}

	sessionIDStr := r.URL.Query().Get("session_id")
	questionNumStr := r.URL.Query().Get("question_num")

	sessionID := 0
	questionNum := 1
	if sessionIDStr != "" {
		sessionID, _ = strconv.Atoi(sessionIDStr)
	}
	if questionNumStr != "" {
		questionNum, _ = strconv.Atoi(questionNumStr)
	}

	var question TestQuestionData

	switch testType {
	case "conjugation":
		question = h.generateConjugationQuestion(cefrLevel)
	case "declension":
		question = h.generateDeclensionQuestion(cefrLevel)
	default:
		question = h.generateVocabularyQuestion(cefrLevel)
	}

	question.SessionID = sessionID
	question.QuestionNum = questionNum
	question.TotalQuestions = 10
	question.QuestionType = testType

	h.renderPartial(w, "test-question.html", question)
}

// AnswerTestQuestion processes a test answer
func (h *Handlers) AnswerTestQuestion(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderPartial(w, "test-result.html", TestAnswerResult{Error: "Invalid form data"})
		return
	}

	answerStr := r.FormValue("answer")
	correctIndexStr := r.FormValue("correct_index")
	questionNumStr := r.FormValue("question_num")
	sessionIDStr := r.FormValue("session_id")

	answer, _ := strconv.Atoi(answerStr)
	correctIndex, _ := strconv.Atoi(correctIndexStr)
	questionNum, _ := strconv.Atoi(questionNumStr)
	sessionID, _ := strconv.Atoi(sessionIDStr)

	correct := answer == correctIndex

	// In a full implementation, we'd update the test session in the database
	result := TestAnswerResult{
		Correct:        correct,
		QuestionNum:    questionNum,
		TotalQuestions: 10,
		SessionID:      sessionID,
		IsComplete:     questionNum >= 10,
	}

	if correct {
		result.Explanation = "Correct! Well done!"
	} else {
		result.Explanation = "Not quite right. Keep practicing!"
	}

	h.renderPartial(w, "test-result.html", result)
}

// generateVocabularyQuestion creates a vocabulary multiple choice question
func (h *Handlers) generateVocabularyQuestion(cefrLevel string) TestQuestionData {
	// Get a random word
	word, err := h.db.GetRandomWord(cefrLevel)
	if err != nil {
		return TestQuestionData{Error: "Failed to generate question"}
	}

	// Get some other words for wrong answers
	words, err := h.db.GetWords(cefrLevel, 10)
	if err != nil || len(words) < 4 {
		return TestQuestionData{Error: "Not enough words for test"}
	}

	// Create options (1 correct, 3 wrong)
	options := []string{word.English}
	for _, w := range words {
		if w.ID != word.ID && len(options) < 4 {
			options = append(options, w.English)
		}
	}

	// Shuffle options
	shuffleStrings(options)

	// Find correct index after shuffle
	correctIndex := 0
	for i, opt := range options {
		if opt == word.English {
			correctIndex = i
			break
		}
	}

	return TestQuestionData{
		Question:     "What is the English translation of \"" + word.Finnish + "\"?",
		Options:      options,
		CorrectIndex: correctIndex,
		Word:         word.Finnish,
		WordID:       word.ID,
	}
}

// generateConjugationQuestion creates a verb conjugation question
func (h *Handlers) generateConjugationQuestion(cefrLevel string) TestQuestionData {
	verb, err := h.db.GetRandomVerb()
	if err != nil {
		return TestQuestionData{Error: "Failed to generate question"}
	}

	// Get conjugation for a random person
	persons := []string{"minä", "sinä", "hän", "me", "te", "he"}
	personIdx := rand.Intn(len(persons))
	person := persons[personIdx]

	// Get correct conjugation
	allForms := h.finnish.ConjugateAll(verb.Infinitive, verb.Type, models.TensePresent)
	if len(allForms) <= personIdx {
		return TestQuestionData{Error: "Failed to conjugate"}
	}
	correctForm := allForms[personIdx]

	// Generate wrong options
	options := []string{correctForm}
	// Add some plausible wrong answers
	wrongEndings := []string{"n", "t", "vat", "mme", "tte", ""}
	for _, ending := range wrongEndings {
		if len(options) >= 4 {
			break
		}
		wrongForm := verb.Stem + ending
		if wrongForm != correctForm {
			options = append(options, wrongForm)
		}
	}

	shuffleStrings(options)

	correctIndex := 0
	for i, opt := range options {
		if opt == correctForm {
			correctIndex = i
			break
		}
	}

	return TestQuestionData{
		Question:     "Conjugate \"" + verb.Infinitive + "\" (" + verb.Translation + ") for \"" + person + "\" in present tense:",
		Options:      options,
		CorrectIndex: correctIndex,
		Word:         verb.Infinitive,
	}
}

// generateDeclensionQuestion creates a noun declension question
func (h *Handlers) generateDeclensionQuestion(cefrLevel string) TestQuestionData {
	noun, err := h.db.GetRandomNoun()
	if err != nil {
		return TestQuestionData{Error: "Failed to generate question"}
	}

	// Pick a random case
	cases := []models.NounCase{
		models.CaseGenitive,
		models.CasePartitive,
		models.CaseInessive,
		models.CaseElative,
		models.CaseIllative,
	}
	targetCase := cases[rand.Intn(len(cases))]

	// Get correct declension
	correctForm := h.finnish.Decline(noun.Nominative, noun.DeclensionType, targetCase)

	// Generate wrong options
	options := []string{correctForm}
	wrongEndings := []string{"ssa", "sta", "an", "lla", "lta", "lle", "n", "a"}
	for _, ending := range wrongEndings {
		if len(options) >= 4 {
			break
		}
		wrongForm := noun.Stem + ending
		if wrongForm != correctForm {
			options = append(options, wrongForm)
		}
	}

	shuffleStrings(options)

	correctIndex := 0
	for i, opt := range options {
		if opt == correctForm {
			correctIndex = i
			break
		}
	}

	return TestQuestionData{
		Question:     "Decline \"" + noun.Nominative + "\" (" + noun.Translation + ") in the " + targetCase.Name() + " case:",
		Options:      options,
		CorrectIndex: correctIndex,
		Word:         noun.Nominative,
	}
}

// shuffleStrings shuffles a slice of strings in place
func shuffleStrings(slice []string) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}