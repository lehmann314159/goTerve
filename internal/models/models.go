package models

import "time"

// CEFRLevel represents the Common European Framework of Reference for Languages level
type CEFRLevel string

const (
	CEFRLevelA1 CEFRLevel = "A1"
	CEFRLevelA2 CEFRLevel = "A2"
	CEFRLevelB1 CEFRLevel = "B1"
	CEFRLevelB2 CEFRLevel = "B2"
	CEFRLevelC1 CEFRLevel = "C1"
	CEFRLevelC2 CEFRLevel = "C2"
)

// User represents a user in the system
type User struct {
	ID                    string    `json:"id"`
	Email                 string    `json:"email"`
	Name                  string    `json:"name"`
	Avatar                string    `json:"avatar,omitempty"`
	GoogleID              string    `json:"googleId,omitempty"`
	CEFRLevel             CEFRLevel `json:"cefrLevel"`
	HasCompletedOnboarding bool      `json:"hasCompletedOnboarding"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// Word represents a vocabulary word
type Word struct {
	ID        int       `json:"id"`
	Finnish   string    `json:"finnish"`
	English   string    `json:"english"`
	Phonetic  string    `json:"phonetic,omitempty"`
	CEFRLevel CEFRLevel `json:"cefrLevel"`
	Frequency int       `json:"frequency"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// VerbType represents Finnish verb conjugation types
type VerbType int

const (
	VerbTypeI   VerbType = 1 // -a/-ä verbs (sanoa, puhua)
	VerbTypeII  VerbType = 2 // -da/-dä verbs (syödä, juoda)
	VerbTypeIII VerbType = 3 // -la/-lä, -na/-nä, -ra/-rä, -ta/-tä verbs
	VerbTypeIV  VerbType = 4 // -ata/-ätä verbs (tavata, hypätä)
	VerbTypeV   VerbType = 5 // -ita/-itä verbs (tarvita, valita)
	VerbTypeVI  VerbType = 6 // -eta/-etä verbs (vanheta, kylmetä)
)

// Verb represents a Finnish verb
type Verb struct {
	ID          int       `json:"id"`
	Infinitive  string    `json:"infinitive"`
	Type        VerbType  `json:"type"`
	Stem        string    `json:"stem"`
	Translation string    `json:"translation"`
	Examples    []string  `json:"examples"`
	Frequency   int       `json:"frequency"`
	CEFRLevel   CEFRLevel `json:"cefrLevel"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NounCase represents Finnish grammatical cases
type NounCase int

const (
	CaseNominative NounCase = 1  // perusmuoto - talo
	CaseGenitive   NounCase = 2  // omanto - talon
	CasePartitive  NounCase = 3  // osanto - taloa
	CaseAccusative NounCase = 4  // kohdanto - talon
	CaseInessive   NounCase = 5  // sisäolento - talossa (in)
	CaseElative    NounCase = 6  // sisäeronto - talosta (from inside)
	CaseIllative   NounCase = 7  // sisätulento - taloon (into)
	CaseAdessive   NounCase = 8  // ulkoolento - talolla (at/on)
	CaseAblative   NounCase = 9  // ulkoeronto - talolta (from)
	CaseAllative   NounCase = 10 // ulkotulento - talolle (to)
)

// CaseName returns the English name of a noun case
func (c NounCase) Name() string {
	names := map[NounCase]string{
		CaseNominative: "nominative",
		CaseGenitive:   "genitive",
		CasePartitive:  "partitive",
		CaseAccusative: "accusative",
		CaseInessive:   "inessive",
		CaseElative:    "elative",
		CaseIllative:   "illative",
		CaseAdessive:   "adessive",
		CaseAblative:   "ablative",
		CaseAllative:   "allative",
	}
	return names[c]
}

// DeclensionType represents Finnish noun declension patterns
type DeclensionType int

const (
	DeclensionTypeI   DeclensionType = 1 // vowel stems: talo, katu
	DeclensionTypeII  DeclensionType = 2 // consonant + vowel: katu, lintu
	DeclensionTypeIII DeclensionType = 3 // consonant clusters: sydän, käsi
	DeclensionTypeIV  DeclensionType = 4 // -nen words: nainen, suomalainen
	DeclensionTypeV   DeclensionType = 5 // -si/-ti words: käsi, lehti
	DeclensionTypeVI  DeclensionType = 6 // special/irregular: mies, yö
)

// Noun represents a Finnish noun
type Noun struct {
	ID             int            `json:"id"`
	Nominative     string         `json:"nominative"`
	Translation    string         `json:"translation"`
	Examples       []string       `json:"examples"`
	DeclensionType DeclensionType `json:"declensionType"`
	Stem           string         `json:"stem"`
	CEFRLevel      CEFRLevel      `json:"cefrLevel"`
	Frequency      int            `json:"frequency"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

// UserFlashcard represents a user's progress on a vocabulary word
type UserFlashcard struct {
	ID            int       `json:"id"`
	UserID        string    `json:"userId"`
	WordID        int       `json:"wordId"`
	Category      string    `json:"category"` // new, learning, review, mastered
	TimesReviewed int       `json:"timesReviewed"`
	TimesCorrect  int       `json:"timesCorrect"`
	LastReviewedAt *time.Time `json:"lastReviewedAt,omitempty"`
	NextReviewAt  *time.Time `json:"nextReviewAt,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// SavedStory represents a generated story saved by a user
type SavedStory struct {
	ID          int       `json:"id"`
	UserID      string    `json:"userId"`
	Story       string    `json:"story"`
	Translation string    `json:"translation"`
	CEFRLevel   CEFRLevel `json:"cefrLevel"`
	Topic       string    `json:"topic,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// TestSession represents a grammar test session
type TestSession struct {
	ID          int       `json:"id"`
	UserID      string    `json:"userId,omitempty"`
	TestType    string    `json:"testType"` // conjugation, declension, vocabulary
	CEFRLevel   CEFRLevel `json:"cefrLevel"`
	TotalQuestions int    `json:"totalQuestions"`
	CorrectAnswers int    `json:"correctAnswers"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// Session represents a user session for authentication
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// Tense represents verb tenses
type Tense int

const (
	TensePresent             Tense = 1
	TenseImperfect           Tense = 2
	TensePerfect             Tense = 3
	TenseImperative          Tense = 4
	TenseConditional         Tense = 5
	TenseNegativePresent     Tense = 11
	TenseNegativeImperfect   Tense = 12
	TenseNegativePerfect     Tense = 13
	TenseNegativeImperative  Tense = 14
	TenseNegativeConditional Tense = 15
)

// TenseName returns the English name of a tense
func (t Tense) Name() string {
	names := map[Tense]string{
		TensePresent:             "present",
		TenseImperfect:           "imperfect",
		TensePerfect:             "perfect",
		TenseImperative:          "imperative",
		TenseConditional:         "conditional",
		TenseNegativePresent:     "negative present",
		TenseNegativeImperfect:   "negative imperfect",
		TenseNegativePerfect:     "negative perfect",
		TenseNegativeImperative:  "negative imperative",
		TenseNegativeConditional: "negative conditional",
	}
	return names[t]
}
