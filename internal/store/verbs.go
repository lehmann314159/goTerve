package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mikebway/goTerve/internal/models"
)

// GetVerbs retrieves all verbs
func (s *Store) GetVerbs() ([]models.Verb, error) {
	query := `SELECT id, infinitive, type, stem, translation, examples, frequency, cefr_level, created_at, updated_at
		FROM verbs ORDER BY frequency ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get verbs: %w", err)
	}
	defer rows.Close()

	var verbs []models.Verb
	for rows.Next() {
		var verb models.Verb
		var examplesJSON string

		err := rows.Scan(&verb.ID, &verb.Infinitive, &verb.Type, &verb.Stem, &verb.Translation,
			&examplesJSON, &verb.Frequency, &verb.CEFRLevel, &verb.CreatedAt, &verb.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan verb: %w", err)
		}

		if err := json.Unmarshal([]byte(examplesJSON), &verb.Examples); err != nil {
			verb.Examples = []string{}
		}
		verbs = append(verbs, verb)
	}

	return verbs, nil
}

// GetVerbByInfinitive retrieves a verb by its infinitive form
func (s *Store) GetVerbByInfinitive(infinitive string) (*models.Verb, error) {
	query := `SELECT id, infinitive, type, stem, translation, examples, frequency, cefr_level, created_at, updated_at
		FROM verbs WHERE infinitive = ?`
	row := s.db.QueryRow(query, infinitive)

	var verb models.Verb
	var examplesJSON string

	err := row.Scan(&verb.ID, &verb.Infinitive, &verb.Type, &verb.Stem, &verb.Translation,
		&examplesJSON, &verb.Frequency, &verb.CEFRLevel, &verb.CreatedAt, &verb.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get verb: %w", err)
	}

	if err := json.Unmarshal([]byte(examplesJSON), &verb.Examples); err != nil {
		verb.Examples = []string{}
	}
	return &verb, nil
}

// GetRandomVerb retrieves a random verb
func (s *Store) GetRandomVerb() (*models.Verb, error) {
	query := `SELECT id, infinitive, type, stem, translation, examples, frequency, cefr_level, created_at, updated_at
		FROM verbs ORDER BY RANDOM() LIMIT 1`
	row := s.db.QueryRow(query)

	var verb models.Verb
	var examplesJSON string

	err := row.Scan(&verb.ID, &verb.Infinitive, &verb.Type, &verb.Stem, &verb.Translation,
		&examplesJSON, &verb.Frequency, &verb.CEFRLevel, &verb.CreatedAt, &verb.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get random verb: %w", err)
	}

	if err := json.Unmarshal([]byte(examplesJSON), &verb.Examples); err != nil {
		verb.Examples = []string{}
	}
	return &verb, nil
}
