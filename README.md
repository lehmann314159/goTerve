# Terve - Finnish Language Learning App

A web application for learning Finnish grammar and vocabulary, built with Go and HTMX.

## Features

### Verb Conjugation
Practice conjugating Finnish verbs across all 6 verb types:
- **Tenses**: Present, Imperfect, Perfect, Conditional, Imperative
- **Negative forms**: All tenses have negative conjugations
- Covers consonant gradation and vowel harmony rules

### Noun Declension
Practice declining Finnish nouns through 10 grammatical cases:
- Nominative, Genitive, Partitive, Inessive, Elative, Illative, Adessive, Ablative, Allative, Accusative

### Vocabulary Flashcards
- Spaced repetition algorithm
- CEFR level filtering (A1-C2)
- Progress tracking (new, learning, review, mastered)

### Reading Practice
- AI-generated Finnish stories using Claude API
- Stories tailored to your CEFR level
- Side-by-side Finnish/English translations

### Grammar Tests
- Multiple choice questions
- Vocabulary, conjugation, and declension tests
- CEFR level selection

## Tech Stack

- **Backend**: Go (standard library + chi router)
- **Frontend**: HTMX + server-side rendered HTML templates
- **Database**: SQLite
- **AI**: Claude API (Haiku) for story generation

## Running Locally

```bash
# Clone the repository
git clone https://github.com/lehmann314159/goTerve.git
cd goTerve

# Run the server
go run ./cmd/server/main.go

# Visit http://localhost:3001
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3001` |
| `DATABASE_PATH` | SQLite database path | `./data/terve.db` |
| `ANTHROPIC_API_KEY` | Claude API key for story generation | (optional) |

## Docker

### Using Docker Compose

```bash
# Create .env file with your API key (optional)
echo "ANTHROPIC_API_KEY=your-key-here" > .env

# Run
docker compose up -d
```

### Pull from Docker Hub

```bash
docker pull lehmann314159/terve:latest
docker run -d -p 3001:3001 -v terve-data:/app/data lehmann314159/terve:latest
```

## Deployment

The app is designed to run behind a reverse proxy (like Caddy). See `docker-compose.yml` for production configuration with external network support.

## Finnish Grammar Reference

### Verb Types
| Type | Ending | Example |
|------|--------|---------|
| 1 | -a/-ä | puhua (to speak), sanoa (to say) |
| 2 | -da/-dä | syödä (to eat), juoda (to drink) |
| 3 | -la/-lä, -na/-nä, -ra/-rä | tulla (to come), mennä (to go) |
| 4 | -ata/-ätä | tavata (to meet), haluta (to want) |
| 5 | -ita/-itä | tarvita (to need), valita (to choose) |
| 6 | -eta/-etä | vanheta (to age), kylmetä (to get cold) |

### Noun Cases
| Case | Use | Example (talo = house) |
|------|-----|------------------------|
| Nominative | Subject | talo |
| Genitive | Possession | talon |
| Partitive | Partial/ongoing | taloa |
| Inessive | Inside | talossa |
| Elative | From inside | talosta |
| Illative | Into | taloon |
| Adessive | At/on | talolla |
| Ablative | From (surface) | talolta |
| Allative | To (surface) | talolle |

## License

MIT