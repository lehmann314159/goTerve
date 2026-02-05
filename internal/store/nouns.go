package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mikebway/goTerve/internal/models"
)

// GetNouns retrieves all nouns
func (s *Store) GetNouns() ([]models.Noun, error) {
	query := `SELECT id, nominative, translation, examples, declension_type, stem, cefr_level, frequency, created_at, updated_at
		FROM nouns ORDER BY frequency ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get nouns: %w", err)
	}
	defer rows.Close()

	var nouns []models.Noun
	for rows.Next() {
		var noun models.Noun
		var examplesJSON string

		err := rows.Scan(&noun.ID, &noun.Nominative, &noun.Translation, &examplesJSON,
			&noun.DeclensionType, &noun.Stem, &noun.CEFRLevel, &noun.Frequency, &noun.CreatedAt, &noun.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan noun: %w", err)
		}

		if err := json.Unmarshal([]byte(examplesJSON), &noun.Examples); err != nil {
			noun.Examples = []string{}
		}
		nouns = append(nouns, noun)
	}

	return nouns, nil
}

// GetNounByNominative retrieves a noun by its nominative form
func (s *Store) GetNounByNominative(nominative string) (*models.Noun, error) {
	query := `SELECT id, nominative, translation, examples, declension_type, stem, cefr_level, frequency, created_at, updated_at
		FROM nouns WHERE nominative = ?`
	row := s.db.QueryRow(query, nominative)

	var noun models.Noun
	var examplesJSON string

	err := row.Scan(&noun.ID, &noun.Nominative, &noun.Translation, &examplesJSON,
		&noun.DeclensionType, &noun.Stem, &noun.CEFRLevel, &noun.Frequency, &noun.CreatedAt, &noun.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get noun: %w", err)
	}

	if err := json.Unmarshal([]byte(examplesJSON), &noun.Examples); err != nil {
		noun.Examples = []string{}
	}
	return &noun, nil
}

// GetRandomNoun retrieves a random noun
func (s *Store) GetRandomNoun() (*models.Noun, error) {
	query := `SELECT id, nominative, translation, examples, declension_type, stem, cefr_level, frequency, created_at, updated_at
		FROM nouns ORDER BY RANDOM() LIMIT 1`
	row := s.db.QueryRow(query)

	var noun models.Noun
	var examplesJSON string

	err := row.Scan(&noun.ID, &noun.Nominative, &noun.Translation, &examplesJSON,
		&noun.DeclensionType, &noun.Stem, &noun.CEFRLevel, &noun.Frequency, &noun.CreatedAt, &noun.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get random noun: %w", err)
	}

	if err := json.Unmarshal([]byte(examplesJSON), &noun.Examples); err != nil {
		noun.Examples = []string{}
	}
	return &noun, nil
}
