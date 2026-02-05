package finnish

import (
	"github.com/mikebway/goTerve/internal/models"
)

// Finnish provides Finnish language grammar operations
type Finnish struct{}

// New creates a new Finnish grammar instance
func New() *Finnish {
	return &Finnish{}
}

// isBackVowel checks if a character is a back vowel (a, o, u)
func isBackVowel(r rune) bool {
	return r == 'a' || r == 'o' || r == 'u' || r == 'A' || r == 'O' || r == 'U'
}

// isFrontVowel checks if a character is a front vowel (ä, ö, y)
func isFrontVowel(r rune) bool {
	return r == 'ä' || r == 'ö' || r == 'y' || r == 'Ä' || r == 'Ö' || r == 'Y'
}

// hasBackVowel checks if a word contains any back vowels
func hasBackVowel(word string) bool {
	for _, r := range word {
		if isBackVowel(r) {
			return true
		}
	}
	return false
}

// getVowelHarmony returns the appropriate vowel based on vowel harmony rules
// For a given word, returns 'a' or 'ä' depending on the word's vowels
func getVowelHarmony(word string) rune {
	if hasBackVowel(word) {
		return 'a'
	}
	return 'ä'
}

// getVowelHarmonyO returns 'o' or 'ö' based on vowel harmony
func getVowelHarmonyO(word string) rune {
	if hasBackVowel(word) {
		return 'o'
	}
	return 'ö'
}

// getStem extracts the stem from a Finnish verb infinitive
func getStem(infinitive string, verbType models.VerbType) string {
	runes := []rune(infinitive)
	length := len(runes)

	switch verbType {
	case models.VerbTypeI:
		// Type 1: -a/-ä verbs (puhua -> puhu, sanoa -> sano)
		if length > 1 {
			return string(runes[:length-1])
		}
	case models.VerbTypeII:
		// Type 2: -da/-dä verbs (syödä -> syö, juoda -> juo)
		if length > 2 {
			return string(runes[:length-2])
		}
	case models.VerbTypeIII:
		// Type 3: -la/-lä, -na/-nä, -ra/-rä, -sta/-stä verbs
		if length > 2 {
			return string(runes[:length-2])
		}
	case models.VerbTypeIV:
		// Type 4: -ata/-ätä verbs (tavata -> tapaa, haluta -> halua)
		if length > 2 {
			// Remove -ta/-tä, the stem often changes
			return string(runes[:length-2])
		}
	case models.VerbTypeV:
		// Type 5: -ita/-itä verbs (tarvita -> tarvitse)
		if length > 2 {
			return string(runes[:length-2])
		}
	case models.VerbTypeVI:
		// Type 6: -eta/-etä verbs (vanheta -> vanhene)
		if length > 2 {
			return string(runes[:length-2])
		}
	}

	return infinitive
}

// applyConsonantGradation applies Finnish consonant gradation rules
// strong -> weak gradation (used in certain conjugation forms)
func applyConsonantGradation(stem string, toWeak bool) string {
	if len(stem) < 2 {
		return stem
	}

	runes := []rune(stem)
	length := len(runes)

	// Find the gradation point (usually near the end of the stem)
	for i := length - 2; i >= 0; i-- {
		// Check for double consonants (kk -> k, pp -> p, tt -> t)
		if i > 0 && runes[i] == runes[i-1] {
			switch runes[i] {
			case 'k':
				if toWeak {
					return string(runes[:i]) + string(runes[i+1:])
				}
			case 'p':
				if toWeak {
					return string(runes[:i]) + string(runes[i+1:])
				}
			case 't':
				if toWeak {
					return string(runes[:i]) + string(runes[i+1:])
				}
			}
		}

		// Check for single consonant gradation
		if toWeak {
			switch runes[i] {
			case 'k':
				// k -> - (disappears) in some contexts
				if i > 0 && isVowel(runes[i-1]) && i < length-1 && isVowel(runes[i+1]) {
					return string(runes[:i]) + string(runes[i+1:])
				}
			case 'p':
				// p -> v
				if i > 0 && isVowel(runes[i-1]) {
					return string(runes[:i]) + "v" + string(runes[i+1:])
				}
			case 't':
				// t -> d
				if i > 0 && isVowel(runes[i-1]) {
					return string(runes[:i]) + "d" + string(runes[i+1:])
				}
			}
		}
	}

	return stem
}

// isVowel checks if a rune is a vowel
func isVowel(r rune) bool {
	vowels := "aeiouäöy"
	for _, v := range vowels {
		if r == v {
			return true
		}
	}
	return false
}