package store

import (
	"fmt"
	"log"
)

// seedWords populates the words table with initial vocabulary
func (s *Store) seedWords() error {
	// Check if already seeded
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check words count: %w", err)
	}
	if count > 0 {
		log.Printf("Words already seeded (%d words)", count)
		return nil
	}

	log.Println("Seeding words...")

	words := []struct {
		Finnish   string
		English   string
		Phonetic  string
		CEFRLevel string
		Frequency int
		Category  string
	}{
		// A1 Level - Most Essential Words
		{"hei", "hello", "hey", "A1", 1, "greetings"},
		{"kiitos", "thank you", "kee-tos", "A1", 2, "greetings"},
		{"anteeksi", "excuse me/sorry", "an-teek-si", "A1", 3, "greetings"},
		{"kyllä", "yes", "kyl-lah", "A1", 4, "basic"},
		{"ei", "no", "ay", "A1", 5, "basic"},
		{"minä", "I", "mi-nah", "A1", 6, "pronouns"},
		{"sinä", "you", "si-nah", "A1", 7, "pronouns"},
		{"hän", "he/she", "hahn", "A1", 8, "pronouns"},
		{"olla", "to be", "ol-la", "A1", 9, "verbs"},
		{"yksi", "one", "yk-si", "A1", 10, "numbers"},
		{"kaksi", "two", "kak-si", "A1", 11, "numbers"},
		{"kolme", "three", "kol-me", "A1", 12, "numbers"},
		{"vesi", "water", "veh-si", "A1", 13, "food"},
		{"ruoka", "food", "ruo-ka", "A1", 14, "food"},
		{"kahvi", "coffee", "kah-vi", "A1", 15, "food"},
		{"koti", "home", "ko-ti", "A1", 16, "places"},
		{"työ", "work", "tyo", "A1", 17, "daily_life"},
		{"koulu", "school", "kou-lu", "A1", 18, "places"},
		{"aika", "time", "ai-ka", "A1", 19, "basic"},
		{"päivä", "day", "pai-va", "A1", 20, "time"},
		{"yö", "night", "yo", "A1", 21, "time"},
		{"mies", "man", "mi-es", "A1", 22, "people"},
		{"nainen", "woman", "nai-nen", "A1", 23, "people"},
		{"lapsi", "child", "lap-si", "A1", 24, "people"},
		{"isä", "father", "i-sah", "A1", 25, "family"},
		{"äiti", "mother", "ai-ti", "A1", 26, "family"},
		{"iso", "big", "i-so", "A1", 27, "adjectives"},
		{"pieni", "small", "pi-e-ni", "A1", 28, "adjectives"},
		{"hyvä", "good", "hy-vah", "A1", 29, "adjectives"},
		{"paha", "bad", "pa-ha", "A1", 30, "adjectives"},
		{"auto", "car", "au-to", "A1", 46, "transport"},
		{"bussi", "bus", "bus-si", "A1", 47, "transport"},
		{"juna", "train", "ju-na", "A1", 48, "transport"},
		{"leipä", "bread", "lei-pa", "A1", 50, "food"},
		{"maito", "milk", "mai-to", "A1", 51, "food"},
		{"liha", "meat", "li-ha", "A1", 52, "food"},
		{"kala", "fish", "ka-la", "A1", 53, "food"},
		{"kissa", "cat", "kis-sa", "A1", 54, "animals"},
		{"koira", "dog", "koi-ra", "A1", 55, "animals"},
		{"lintu", "bird", "lin-tu", "A1", 56, "animals"},
		{"kukka", "flower", "kuk-ka", "A1", 57, "nature"},
		{"puu", "tree", "puu", "A1", 58, "nature"},
		{"järvi", "lake", "jar-vi", "A1", 60, "nature"},
		{"metsä", "forest", "met-sa", "A1", 61, "nature"},
		{"talvi", "winter", "tal-vi", "A1", 65, "seasons"},
		{"kevät", "spring", "ke-vat", "A1", 66, "seasons"},
		{"kesä", "summer", "ke-sa", "A1", 67, "seasons"},
		{"syksy", "autumn", "syk-sy", "A1", 68, "seasons"},
		{"lumi", "snow", "lu-mi", "A1", 69, "weather"},
		{"sade", "rain", "sa-de", "A1", 70, "weather"},
		{"aurinko", "sun", "au-rin-ko", "A1", 72, "weather"},
		{"kirja", "book", "kir-ja", "A1", 75, "education"},
		{"pöytä", "table", "poy-ta", "A1", 78, "furniture"},
		{"tuoli", "chair", "tuo-li", "A1", 79, "furniture"},
		{"ikkuna", "window", "ik-ku-na", "A1", 81, "house"},
		{"ovi", "door", "o-vi", "A1", 82, "house"},

		// A2 Level
		{"opiskella", "to study", "o-pis-kel-la", "A2", 31, "verbs"},
		{"puhua", "to speak", "pu-hua", "A2", 32, "verbs"},
		{"ymmärtää", "to understand", "ym-mar-taa", "A2", 33, "verbs"},
		{"tietää", "to know", "ti-e-taa", "A2", 34, "verbs"},
		{"kieli", "language", "ki-e-li", "A2", 35, "education"},
		{"ystävä", "friend", "ys-ta-va", "A2", 36, "people"},
		{"kaupunki", "city", "kau-pun-ki", "A2", 37, "places"},
		{"maa", "country", "maa", "A2", 38, "places"},
		{"kaunis", "beautiful", "kau-nis", "A2", 39, "adjectives"},
		{"viikko", "week", "viik-ko", "A2", 40, "time"},
		{"syödä", "to eat", "syo-da", "A2", 98, "verbs"},
		{"juoda", "to drink", "juo-da", "A2", 99, "verbs"},
		{"nukkua", "to sleep", "nuk-ku-a", "A2", 100, "verbs"},
		{"katsoa", "to watch/look", "kat-so-a", "A2", 102, "verbs"},
		{"kuunnella", "to listen", "kuun-nel-la", "A2", 103, "verbs"},
		{"lukea", "to read", "lu-ke-a", "A2", 104, "verbs"},
		{"kirjoittaa", "to write", "kir-joit-taa", "A2", 105, "verbs"},
		{"ostaa", "to buy", "os-taa", "A2", 106, "verbs"},
		{"raha", "money", "ra-ha", "A2", 109, "finance"},
		{"kauppa", "store/shop", "kaup-pa", "A2", 110, "places"},
		{"ravintola", "restaurant", "ra-vin-to-la", "A2", 111, "places"},

		// B1 Level
		{"keskustella", "to discuss", "kes-kus-tel-la", "B1", 41, "verbs"},
		{"päättää", "to decide", "paat-taa", "B1", 42, "verbs"},
		{"kokea", "to experience", "ko-ke-a", "B1", 43, "verbs"},
		{"yhteiskunta", "society", "yh-teis-kun-ta", "B1", 44, "society"},
		{"kulttuuri", "culture", "kult-tu-u-ri", "B1", 45, "society"},
	}

	stmt, err := s.db.Prepare(`INSERT INTO words (finnish, english, phonetic, cefr_level, frequency, category) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare word insert: %w", err)
	}
	defer stmt.Close()

	for _, w := range words {
		_, err := stmt.Exec(w.Finnish, w.English, w.Phonetic, w.CEFRLevel, w.Frequency, w.Category)
		if err != nil {
			return fmt.Errorf("failed to insert word %s: %w", w.Finnish, err)
		}
	}

	log.Printf("Seeded %d words", len(words))
	return nil
}

// seedVerbs populates the verbs table
func (s *Store) seedVerbs() error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM verbs").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check verbs count: %w", err)
	}
	if count > 0 {
		log.Printf("Verbs already seeded (%d verbs)", count)
		return nil
	}

	log.Println("Seeding verbs...")

	verbs := []struct {
		Infinitive  string
		Type        int
		Stem        string
		Translation string
		Frequency   int
		CEFRLevel   string
	}{
		{"olla", 3, "ole", "to be", 1, "A1"},
		{"tehdä", 2, "tee", "to do/make", 2, "A1"},
		{"sanoa", 1, "sano", "to say", 3, "A1"},
		{"mennä", 3, "mene", "to go", 4, "A1"},
		{"tulla", 3, "tule", "to come", 5, "A1"},
		{"antaa", 1, "anna", "to give", 6, "A1"},
		{"ottaa", 1, "ota", "to take", 7, "A1"},
		{"nähdä", 2, "näe", "to see", 8, "A1"},
		{"tietää", 1, "tiedä", "to know", 9, "A2"},
		{"voida", 2, "voi", "to be able/can", 10, "A1"},
		{"haluta", 3, "halua", "to want", 11, "A1"},
		{"tavata", 4, "tapaa", "to meet", 12, "A2"},
		{"tarvita", 5, "tarvitse", "to need", 13, "A2"},
		{"vanheta", 6, "vanhene", "to age", 14, "B1"},
	}

	stmt, err := s.db.Prepare(`INSERT INTO verbs (infinitive, type, stem, translation, examples, frequency, cefr_level) VALUES (?, ?, ?, ?, '[]', ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare verb insert: %w", err)
	}
	defer stmt.Close()

	for _, v := range verbs {
		_, err := stmt.Exec(v.Infinitive, v.Type, v.Stem, v.Translation, v.Frequency, v.CEFRLevel)
		if err != nil {
			return fmt.Errorf("failed to insert verb %s: %w", v.Infinitive, err)
		}
	}

	log.Printf("Seeded %d verbs", len(verbs))
	return nil
}

// seedNouns populates the nouns table
func (s *Store) seedNouns() error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM nouns").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check nouns count: %w", err)
	}
	if count > 0 {
		log.Printf("Nouns already seeded (%d nouns)", count)
		return nil
	}

	log.Println("Seeding nouns...")

	nouns := []struct {
		Nominative     string
		Translation    string
		DeclensionType int
		Stem           string
		Frequency      int
		CEFRLevel      string
	}{
		{"talo", "house", 1, "talo", 1, "A1"},
		{"katu", "street", 1, "katu", 2, "A1"},
		{"auto", "car", 1, "auto", 3, "A1"},
		{"koti", "home", 1, "koti", 4, "A1"},
		{"lintu", "bird", 1, "lintu", 5, "A1"},
		{"käsi", "hand", 5, "käte", 6, "A1"},
		{"sydän", "heart", 3, "sydäme", 7, "A2"},
		{"nainen", "woman", 4, "naise", 8, "A1"},
		{"mies", "man", 6, "miehe", 9, "A1"},
		{"lapsi", "child", 5, "lapse", 10, "A1"},
		{"vesi", "water", 5, "vete", 11, "A1"},
		{"kirja", "book", 1, "kirja", 12, "A1"},
		{"ruoka", "food", 1, "ruoka", 13, "A1"},
		{"pää", "head", 1, "pää", 14, "A1"},
		{"yö", "night", 6, "yö", 15, "A1"},
		{"päivä", "day", 1, "päivä", 16, "A1"},
		{"aika", "time", 1, "aika", 17, "A2"},
		{"työ", "work", 6, "työ", 18, "A2"},
		{"koulu", "school", 1, "koulu", 19, "A1"},
		{"maa", "country/land", 1, "maa", 20, "A1"},
	}

	stmt, err := s.db.Prepare(`INSERT INTO nouns (nominative, translation, examples, declension_type, stem, frequency, cefr_level) VALUES (?, ?, '[]', ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare noun insert: %w", err)
	}
	defer stmt.Close()

	for _, n := range nouns {
		_, err := stmt.Exec(n.Nominative, n.Translation, n.DeclensionType, n.Stem, n.Frequency, n.CEFRLevel)
		if err != nil {
			return fmt.Errorf("failed to insert noun %s: %w", n.Nominative, err)
		}
	}

	log.Printf("Seeded %d nouns", len(nouns))
	return nil
}
