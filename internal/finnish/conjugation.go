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
// Returns: [minä, sinä, hän, me, te, he] for most tenses
// For imperative: [sinä, te, me, hän, he] (common forms)
func (f *Finnish) ConjugateAll(infinitive string, verbType models.VerbType, tense models.Tense) []string {
	switch tense {
	case models.TensePresent:
		return f.conjugatePresent(infinitive, verbType)
	case models.TenseImperfect:
		return f.conjugateImperfect(infinitive, verbType)
	case models.TensePerfect:
		return f.conjugatePerfect(infinitive, verbType)
	case models.TenseImperative:
		return f.conjugateImperative(infinitive, verbType)
	case models.TenseConditional:
		return f.conjugateConditional(infinitive, verbType)
	case models.TenseNegativePresent:
		return f.conjugateNegativePresent(infinitive, verbType)
	case models.TenseNegativeImperfect:
		return f.conjugateNegativeImperfect(infinitive, verbType)
	case models.TenseNegativePerfect:
		return f.conjugateNegativePerfect(infinitive, verbType)
	case models.TenseNegativeImperative:
		return f.conjugateNegativeImperative(infinitive, verbType)
	case models.TenseNegativeConditional:
		return f.conjugateNegativeConditional(infinitive, verbType)
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

// conjugatePerfect returns perfect tense conjugations
// Perfect = olla (present) + past participle (-nut/-nyt for singular, -neet for plural)
func (f *Finnish) conjugatePerfect(infinitive string, verbType models.VerbType) []string {
	participle := f.getPastParticiple(infinitive, verbType)
	participlePlural := f.getPastParticiplePlural(infinitive, verbType)

	return []string{
		"olen " + participle,       // minä olen puhunut
		"olet " + participle,       // sinä olet puhunut
		"on " + participle,         // hän on puhunut
		"olemme " + participlePlural, // me olemme puhuneet
		"olette " + participlePlural, // te olette puhuneet
		"ovat " + participlePlural,   // he ovat puhuneet
	}
}

// getPastParticiple returns the singular past participle (-nut/-nyt)
func (f *Finnish) getPastParticiple(infinitive string, verbType models.VerbType) string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	// Determine -nut or -nyt based on vowel harmony
	ending := "nut"
	if harmony == 'ä' {
		ending = "nyt"
	}

	switch verbType {
	case models.VerbTypeI:
		// puhua -> puhunut, sanoa -> sanonut
		return stem + ending
	case models.VerbTypeII:
		// syödä -> syönyt, juoda -> juonut
		return stem + ending
	case models.VerbTypeIII:
		// tulla -> tullut, mennä -> mennyt
		// Double the consonant
		runes := []rune(stem)
		if len(runes) > 0 {
			lastChar := runes[len(runes)-1]
			return stem + string(lastChar) + "ut"
		}
		return stem + ending
	case models.VerbTypeIV:
		// tavata -> tavannut, haluta -> halunnut
		return stem + "n" + ending
	case models.VerbTypeV:
		// tarvita -> tarvinnut
		return stem + "n" + ending
	case models.VerbTypeVI:
		// vanheta -> vanhennut
		return stem + "n" + ending
	}

	return stem + ending
}

// getPastParticiplePlural returns the plural past participle (-neet)
func (f *Finnish) getPastParticiplePlural(infinitive string, verbType models.VerbType) string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	// Determine -neet or -neet based on vowel harmony
	ending := "neet"
	if harmony == 'ä' {
		ending = "neet"
	}

	switch verbType {
	case models.VerbTypeI:
		// puhua -> puhuneet
		return stem + ending
	case models.VerbTypeII:
		// syödä -> syöneet
		return stem + ending
	case models.VerbTypeIII:
		// tulla -> tulleet, mennä -> menneet
		runes := []rune(stem)
		if len(runes) > 0 {
			lastChar := runes[len(runes)-1]
			return stem + string(lastChar) + "eet"
		}
		return stem + ending
	case models.VerbTypeIV:
		// tavata -> tavanneet
		return stem + "n" + ending
	case models.VerbTypeV:
		// tarvita -> tarvinneet
		return stem + "n" + ending
	case models.VerbTypeVI:
		// vanheta -> vanhenneet
		return stem + "n" + ending
	}

	return stem + ending
}

// conjugateImperative returns imperative (command) forms
// Returns: [sinä (2sg), te (2pl), me (1pl), hän (3sg), he (3pl)]
func (f *Finnish) conjugateImperative(infinitive string, verbType models.VerbType) []string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	var impStem string
	var secondSg string

	switch verbType {
	case models.VerbTypeI:
		// puhua -> puhu!, puhukaa!
		weakStem := applyConsonantGradation(stem, true)
		secondSg = weakStem
		impStem = stem
	case models.VerbTypeII:
		// syödä -> syö!, syökää!
		secondSg = stem
		impStem = stem
	case models.VerbTypeIII:
		// tulla -> tule!, tulkaa!
		secondSg = stem + "e"
		impStem = stem
	case models.VerbTypeIV:
		// tavata -> tapaa!, tavatkaa!
		secondSg = stem + string(harmony) + string(harmony)
		impStem = stem + string(harmony) + "t"
	case models.VerbTypeV:
		// tarvita -> tarvitse!, tarvitakaa!
		secondSg = stem + "tse"
		impStem = stem + "t"
	case models.VerbTypeVI:
		// vanheta -> vanhene!, vanetkaa!
		secondSg = stem + "ne"
		impStem = stem + "t"
	default:
		secondSg = stem
		impStem = stem
	}

	return []string{
		secondSg + "!",                              // sinä: puhu!
		impStem + "k" + string(harmony) + string(harmony) + "!", // te: puhukaa!
		impStem + "k" + string(harmony) + string(harmony) + "mme!", // me: puhukaamme!
		impStem + "k" + string(harmony) + string(harmony) + "n!", // hän: puhukoon!
		impStem + "k" + string(harmony) + string(harmony) + "t!", // he: puhukoot!
	}
}

// conjugateNegativePresent returns negative present tense conjugations
// Negative uses "ei" auxiliary + verb stem (connegative)
func (f *Finnish) conjugateNegativePresent(infinitive string, verbType models.VerbType) []string {
	stem := getStem(infinitive, verbType)

	// Get the connegative stem (same as 2nd person singular without ending)
	var conneg string

	switch verbType {
	case models.VerbTypeI:
		conneg = applyConsonantGradation(stem, true)
	case models.VerbTypeII:
		conneg = stem
	case models.VerbTypeIII:
		conneg = stem + "e"
	case models.VerbTypeIV:
		harmony := getVowelHarmony(infinitive)
		conneg = stem + string(harmony) + string(harmony)
	case models.VerbTypeV:
		conneg = stem + "tse"
	case models.VerbTypeVI:
		conneg = stem + "ne"
	default:
		conneg = stem
	}

	return []string{
		"en " + conneg,      // minä en puhu
		"et " + conneg,      // sinä et puhu
		"ei " + conneg,      // hän ei puhu
		"emme " + conneg,    // me emme puhu
		"ette " + conneg,    // te ette puhu
		"eivät " + conneg,   // he eivät puhu
	}
}

// conjugateNegativeImperfect returns negative imperfect (past) tense conjugations
// Uses "ei" + past participle
func (f *Finnish) conjugateNegativeImperfect(infinitive string, verbType models.VerbType) []string {
	participle := f.getPastParticiple(infinitive, verbType)
	participlePlural := f.getPastParticiplePlural(infinitive, verbType)

	return []string{
		"en " + participle,        // minä en puhunut
		"et " + participle,        // sinä et puhunut
		"ei " + participle,        // hän ei puhunut
		"emme " + participlePlural, // me emme puhuneet
		"ette " + participlePlural, // te ette puhuneet
		"eivät " + participlePlural, // he eivät puhuneet
	}
}

// conjugateNegativePerfect returns negative perfect tense conjugations
// Uses "ei ole" + past participle
func (f *Finnish) conjugateNegativePerfect(infinitive string, verbType models.VerbType) []string {
	participle := f.getPastParticiple(infinitive, verbType)
	participlePlural := f.getPastParticiplePlural(infinitive, verbType)

	return []string{
		"en ole " + participle,        // minä en ole puhunut
		"et ole " + participle,        // sinä et ole puhunut
		"ei ole " + participle,        // hän ei ole puhunut
		"emme ole " + participlePlural, // me emme ole puhuneet
		"ette ole " + participlePlural, // te ette ole puhuneet
		"eivät ole " + participlePlural, // he eivät ole puhuneet
	}
}

// conjugateNegativeConditional returns negative conditional tense conjugations
// Uses "ei" + conditional connegative
func (f *Finnish) conjugateNegativeConditional(infinitive string, verbType models.VerbType) []string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	// Conditional connegative is stem + isi
	var condConneg string

	switch verbType {
	case models.VerbTypeI:
		weakStem := applyConsonantGradation(stem, true)
		condConneg = weakStem + "isi"
	case models.VerbTypeII:
		condConneg = stem + "isi"
	case models.VerbTypeIII:
		condConneg = stem + "isi"
	case models.VerbTypeIV:
		condConneg = stem + string(harmony) + "isi"
	case models.VerbTypeV:
		condConneg = stem + "tsisi"
	case models.VerbTypeVI:
		condConneg = stem + "nisi"
	default:
		condConneg = stem + "isi"
	}

	return []string{
		"en " + condConneg,      // minä en puhuisi
		"et " + condConneg,      // sinä et puhuisi
		"ei " + condConneg,      // hän ei puhuisi
		"emme " + condConneg,    // me emme puhuisi
		"ette " + condConneg,    // te ette puhuisi
		"eivät " + condConneg,   // he eivät puhuisi
	}
}

// conjugateNegativeImperative returns negative imperative (prohibition) forms
// Uses "älä/älkää" + connegative with -ko/-kö
func (f *Finnish) conjugateNegativeImperative(infinitive string, verbType models.VerbType) []string {
	stem := getStem(infinitive, verbType)
	harmony := getVowelHarmony(infinitive)

	// Get connegative stem for imperative
	var conneg string
	var impConneg string

	switch verbType {
	case models.VerbTypeI:
		conneg = applyConsonantGradation(stem, true)
		impConneg = stem + "k" + string(harmony)
	case models.VerbTypeII:
		conneg = stem
		impConneg = stem + "k" + string(harmony)
	case models.VerbTypeIII:
		conneg = stem + "e"
		impConneg = stem + "k" + string(harmony)
	case models.VerbTypeIV:
		conneg = stem + string(harmony) + string(harmony)
		impConneg = stem + string(harmony) + "tk" + string(harmony)
	case models.VerbTypeV:
		conneg = stem + "tse"
		impConneg = stem + "tk" + string(harmony)
	case models.VerbTypeVI:
		conneg = stem + "ne"
		impConneg = stem + "tk" + string(harmony)
	default:
		conneg = stem
		impConneg = stem + "k" + string(harmony)
	}

	// älä/älkää forms
	ala := "älä"
	alkaa := "älkää"
	alkaamme := "älkäämme"
	alkoon := "älköön"
	alkoot := "älkööt"
	if harmony == 'a' {
		// Back vowel harmony - but älä etc. always use front vowels
		// älä, älkää etc. are fixed forms
	}

	return []string{
		ala + " " + conneg + "!",           // älä puhu!
		alkaa + " " + impConneg + "!",      // älkää puhuko!
		alkaamme + " " + impConneg + "!",   // älkäämme puhuko!
		alkoon + " " + impConneg + "!",     // älköön puhuko!
		alkoot + " " + impConneg + "!",     // älkööt puhuko!
	}
}