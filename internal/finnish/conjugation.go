package finnish

import (
	"github.com/mikebway/goTerve/internal/models"
)

// PersonEndings contains verb endings for each person
type PersonEndings struct {
	First   string // minä
	Second  string // sinä
	Third   string // hän
	FirstPl string // me
	SecondPl string // te
	ThirdPl string // he
}

// Conjugate conjugates a Finnish verb for a specific person and tense
func (f *Finnish) Conjugate(infinitive string, verbType models.VerbType, tense models.Tense, person int, plural bool) string {
	forms := f.ConjugateAll(infinitive, verbType, tense)

	index := person - 1
	if plural {
		index += 3
	}

	if index >= 0 && index < len(forms) {
		return forms[index]
	}
	return infinitive
}

// ConjugateAll returns all conjugated forms for a verb in a given tense
// Returns: [minä, sinä, hän, me, te, he]
func (f *Finnish) ConjugateAll(infinitive string, verbType models.VerbType, tense models.Tense) []string {
	switch tense {
	case models.TensePresent:
		return f.conjugatePresent(infinitive, verbType)
	case models.TenseImperfect:
		return f.conjugateImperfect(infinitive, verbType)
	case models.TenseConditional:
		return f.conjugateConditional(infinitive, verbType)
	default:
		return f.conjugatePresent(infinitive, verbType)
	}
}

// conjugatePresent returns present tense conjugations
func (f *Finnish) conjugatePresent(infinitive string, verbType models.VerbType) []string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	var weakStem string
	var strongStem string

	switch verbType {
	case models.VerbTypeI:
		// Type 1: puhua, sanoa - weak grade for 1st/2nd person, strong for 3rd
		weakStem = applyConsonantGradation(stem, true)
		strongStem = stem
		return []string{
			weakStem + "n",                    // minä puhun
			weakStem + "t",                    // sinä puhut
			strongStem + strongStem[len(strongStem)-1:], // hän puhuu (double last vowel)
			weakStem + "mme",                  // me puhumme
			weakStem + "tte",                  // te puhutte
			strongStem + "vat",                // he puhuvat (if back) / -vät (if front)
		}

	case models.VerbTypeII:
		// Type 2: syödä, juoda - stem ends in long vowel or diphthong
		return []string{
			stem + "n",                   // minä syön
			stem + "t",                   // sinä syöt
			stem,                         // hän syö
			stem + "mme",                 // me syömme
			stem + "tte",                 // te syötte
			stem + "v" + string(harmony) + "t", // he syövät
		}

	case models.VerbTypeIII:
		// Type 3: tulla, mennä, purra, nousta
		// These verbs have -e- stem in present
		eStem := stem + "e"
		return []string{
			eStem + "n",                   // minä tulen
			eStem + "t",                   // sinä tulet
			eStem + "e",                   // hän tulee
			eStem + "mme",                 // me tulemme
			eStem + "tte",                 // te tulette
			eStem + "v" + string(harmony) + "t", // he tulevat
		}

	case models.VerbTypeIV:
		// Type 4: tavata, haluta - stem with double vowel
		// tavata -> tapaan, haluta -> haluan
		doubleStem := stem + string(harmony) // adds a or ä
		return []string{
			doubleStem + string(harmony) + "n", // minä tapaan
			doubleStem + string(harmony) + "t", // sinä tapaat
			doubleStem + string(harmony),       // hän tapaa
			doubleStem + string(harmony) + "mme", // me tapaamme
			doubleStem + string(harmony) + "tte", // te tapaatte
			doubleStem + string(harmony) + "v" + string(harmony) + "t", // he tapaavat
		}

	case models.VerbTypeV:
		// Type 5: tarvita, valita - stem with -tse-
		tseStem := stem + "tse"
		return []string{
			tseStem + "n",                   // minä tarvitsen
			tseStem + "t",                   // sinä tarvitset
			tseStem + "e",                   // hän tarvitsee
			tseStem + "mme",                 // me tarvitsemme
			tseStem + "tte",                 // te tarvitsette
			tseStem + "v" + string(harmony) + "t", // he tarvitsevat
		}

	case models.VerbTypeVI:
		// Type 6: vanheta, kylmetä - stem with -ne-
		neStem := stem + "ne"
		return []string{
			neStem + "n",                   // minä vanhenen
			neStem + "t",                   // sinä vanhenet
			neStem + "e",                   // hän vanhenee
			neStem + "mme",                 // me vanhenemme
			neStem + "tte",                 // te vanhenette
			neStem + "v" + string(harmony) + "t", // he vanhenevat
		}
	}

	// Fallback
	return []string{stem + "n", stem + "t", stem, stem + "mme", stem + "tte", stem + "vat"}
}

// conjugateImperfect returns imperfect (past) tense conjugations
func (f *Finnish) conjugateImperfect(infinitive string, verbType models.VerbType) []string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	// Imperfect marker is -i- (changes to -si- after certain vowels)
	var imperfectStem string

	switch verbType {
	case models.VerbTypeI:
		// puhua -> puhui, sanoa -> sanoi
		weakStem := applyConsonantGradation(stem, true)
		imperfectStem = weakStem + "i"
	case models.VerbTypeII:
		// syödä -> söi, juoda -> joi
		imperfectStem = stem + "i"
	case models.VerbTypeIII:
		// tulla -> tuli, mennä -> meni
		imperfectStem = stem + "i"
	case models.VerbTypeIV:
		// tavata -> tapasi, haluta -> halusi
		imperfectStem = stem + string(harmony) + "si"
	case models.VerbTypeV:
		// tarvita -> tarvitsi
		imperfectStem = stem + "tsi"
	case models.VerbTypeVI:
		// vanheta -> vanheni
		imperfectStem = stem + "ni"
	default:
		imperfectStem = stem + "i"
	}

	return []string{
		imperfectStem + "n",                        // minä puhuin
		imperfectStem + "t",                        // sinä puhuit
		imperfectStem,                              // hän puhui
		imperfectStem + "mme",                      // me puhuimme
		imperfectStem + "tte",                      // te puhuitte
		imperfectStem + "v" + string(harmony) + "t", // he puhuivat
	}
}

// conjugateConditional returns conditional tense conjugations
func (f *Finnish) conjugateConditional(infinitive string, verbType models.VerbType) []string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	// Conditional marker is -isi-
	var condStem string

	switch verbType {
	case models.VerbTypeI:
		weakStem := applyConsonantGradation(stem, true)
		condStem = weakStem + "isi"
	case models.VerbTypeII:
		condStem = stem + "isi"
	case models.VerbTypeIII:
		condStem = stem + "isi"
	case models.VerbTypeIV:
		condStem = stem + string(harmony) + "isi"
	case models.VerbTypeV:
		condStem = stem + "tsisi"
	case models.VerbTypeVI:
		condStem = stem + "nisi"
	default:
		condStem = stem + "isi"
	}

	return []string{
		condStem + "n",                        // minä puhuisin
		condStem + "t",                        // sinä puhuisit
		condStem,                              // hän puhuisi
		condStem + "mme",                      // me puhuisimme
		condStem + "tte",                      // te puhuisitte
		condStem + "v" + string(harmony) + "t", // he puhuisivat
	}
}