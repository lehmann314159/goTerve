package finnish

import (
	"strings"

	"github.com/mikebway/goTerve/internal/models"
)

// Decline declines a Finnish noun to the specified case
func (f *Finnish) Decline(nominative string, declType models.DeclensionType, nounCase models.NounCase) string {
	if nounCase == models.CaseNominative {
		return nominative
	}

	stem := getNounStem(nominative, declType)
	harmony := getVowelHarmony(nominative)

	switch nounCase {
	case models.CaseGenitive:
		return f.declineGenitive(nominative, stem, declType)
	case models.CasePartitive:
		return f.declinePartitive(nominative, stem, declType, harmony)
	case models.CaseInessive:
		return stem + "ss" + string(harmony)
	case models.CaseElative:
		return stem + "st" + string(harmony)
	case models.CaseIllative:
		return f.declineIllative(nominative, stem, declType)
	case models.CaseAdessive:
		return stem + "ll" + string(harmony)
	case models.CaseAblative:
		return stem + "lt" + string(harmony)
	case models.CaseAllative:
		return stem + "lle"
	case models.CaseAccusative:
		// Accusative is same as genitive for most nouns
		return f.declineGenitive(nominative, stem, declType)
	}

	return nominative
}

// getNounStem extracts the stem from a Finnish noun
func getNounStem(nominative string, declType models.DeclensionType) string {
	runes := []rune(nominative)
	length := len(runes)

	switch declType {
	case models.DeclensionTypeI:
		// Simple vowel stems: talo -> talo, katu -> kadu
		return nominative
	case models.DeclensionTypeII:
		// Consonant + vowel: lintu -> linnu
		return nominative
	case models.DeclensionTypeIII:
		// Consonant clusters with vowel change
		return nominative
	case models.DeclensionTypeIV:
		// -nen words: nainen -> nais-, suomalainen -> suomalais-
		if strings.HasSuffix(nominative, "nen") && length > 3 {
			return string(runes[:length-3]) + "s"
		}
		return nominative
	case models.DeclensionTypeV:
		// -si words: käsi -> käde-, lehti -> lehde-
		if strings.HasSuffix(nominative, "si") && length > 2 {
			return string(runes[:length-2]) + "de"
		}
		if strings.HasSuffix(nominative, "ti") && length > 2 {
			return string(runes[:length-2]) + "de"
		}
		return nominative
	case models.DeclensionTypeVI:
		// Special/irregular: mies -> miehe-
		return nominative
	}

	return nominative
}

// declineGenitive returns the genitive form
func (f *Finnish) declineGenitive(nominative, stem string, declType models.DeclensionType) string {
	switch declType {
	case models.DeclensionTypeI:
		// talo -> talon
		return stem + "n"
	case models.DeclensionTypeII:
		// lintu -> linnun (with gradation)
		weakStem := applyConsonantGradation(stem, true)
		return weakStem + "n"
	case models.DeclensionTypeIII:
		return stem + "n"
	case models.DeclensionTypeIV:
		// nainen -> naisen
		return stem + "en"
	case models.DeclensionTypeV:
		// käsi -> käden
		return stem + "n"
	case models.DeclensionTypeVI:
		return stem + "n"
	}
	return stem + "n"
}

// declinePartitive returns the partitive form
func (f *Finnish) declinePartitive(nominative, stem string, declType models.DeclensionType, harmony rune) string {
	runes := []rune(nominative)
	length := len(runes)

	switch declType {
	case models.DeclensionTypeI:
		// Check if ends in a vowel
		if length > 0 && isVowel(runes[length-1]) {
			// talo -> taloa, auto -> autoa
			return nominative + string(harmony)
		}
		// Consonant ending - add ta/tä
		return nominative + "t" + string(harmony)
	case models.DeclensionTypeII:
		// lintu -> lintua
		return nominative + string(harmony)
	case models.DeclensionTypeIII:
		return nominative + "t" + string(harmony)
	case models.DeclensionTypeIV:
		// nainen -> naista
		return stem + "t" + string(harmony)
	case models.DeclensionTypeV:
		// käsi -> kättä
		return stem + "tt" + string(harmony)
	case models.DeclensionTypeVI:
		return stem + "t" + string(harmony)
	}
	return nominative + string(harmony)
}

// declineIllative returns the illative form
func (f *Finnish) declineIllative(nominative, stem string, declType models.DeclensionType) string {
	runes := []rune(nominative)
	length := len(runes)
	harmony := getVowelHarmony(nominative)

	switch declType {
	case models.DeclensionTypeI:
		// If ends in a single vowel, double it and add n
		if length > 0 && isVowel(runes[length-1]) {
			// talo -> taloon, katu -> katuun
			lastVowel := runes[length-1]
			return nominative + string(lastVowel) + "n"
		}
		// Otherwise add -iin or similar
		return nominative + "iin"
	case models.DeclensionTypeII:
		if length > 0 && isVowel(runes[length-1]) {
			lastVowel := runes[length-1]
			return nominative + string(lastVowel) + "n"
		}
		return nominative + "un"
	case models.DeclensionTypeIII:
		return nominative + "een"
	case models.DeclensionTypeIV:
		// nainen -> naiseen
		return stem + "een"
	case models.DeclensionTypeV:
		// käsi -> käteen
		return stem + "en"
	case models.DeclensionTypeVI:
		return stem + "h" + string(harmony) + "n"
	}

	// Default: double last vowel + n
	if length > 0 && isVowel(runes[length-1]) {
		lastVowel := runes[length-1]
		return nominative + string(lastVowel) + "n"
	}
	return nominative + "iin"
}

// DeclineAll returns all case forms for a noun
func (f *Finnish) DeclineAll(nominative string, declType models.DeclensionType) map[string]string {
	cases := []models.NounCase{
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

	result := make(map[string]string)
	for _, c := range cases {
		result[c.Name()] = f.Decline(nominative, declType, c)
	}
	return result
}