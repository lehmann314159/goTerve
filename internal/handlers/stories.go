package handlers

import (
	"net/http"
)

// StoryData contains data for story rendering
type StoryData struct {
	Story       string `json:"story"`
	Translation string `json:"translation"`
	CEFRLevel   string `json:"cefrLevel"`
	Topic       string `json:"topic"`
	Error       string `json:"error,omitempty"`
}

// GenerateStory generates a Finnish reading story
// Note: This is a placeholder - full implementation would use Claude API
func (h *Handlers) GenerateStory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderPartial(w, "story.html", StoryData{Error: "Invalid form data"})
		return
	}

	cefrLevel := r.FormValue("level")
	if cefrLevel == "" {
		cefrLevel = "A1"
	}

	topic := r.FormValue("topic")
	if topic == "" {
		topic = "daily life"
	}

	// For now, return a sample story
	// In production, this would call the Claude API
	story := getSampleStory(cefrLevel, topic)

	h.renderPartial(w, "story.html", story)
}

// getSampleStory returns a sample story based on level
func getSampleStory(level, topic string) StoryData {
	stories := map[string]StoryData{
		"A1": {
			Story: `Minun nimeni on Anna. Minä asun Helsingissä.
Minä olen opiskelija. Minä opiskelen suomea.
Joka päivä minä käyn koulussa.
Minä pidän suomen kielestä. Se on kaunis kieli.`,
			Translation: `My name is Anna. I live in Helsinki.
I am a student. I study Finnish.
Every day I go to school.
I like the Finnish language. It is a beautiful language.`,
			CEFRLevel: "A1",
			Topic:     topic,
		},
		"A2": {
			Story: `Eilen menin kauppaan ostamaan ruokaa.
Ostin leipää, maitoa ja juustoa.
Kaupassa tapusin vanhan ystäväni Mikon.
Me puhuimme hetken ja sovimme tapaamisen kahvilaan ensi viikolla.
Oli mukava nähdä hänet pitkästä aikaa.`,
			Translation: `Yesterday I went to the store to buy food.
I bought bread, milk and cheese.
At the store I met my old friend Mikko.
We talked for a while and agreed to meet at a café next week.
It was nice to see him after a long time.`,
			CEFRLevel: "A2",
			Topic:     topic,
		},
		"B1": {
			Story: `Suomessa on neljä vuodenaikaa: kevät, kesä, syksy ja talvi.
Talvi on pisin vuodenaika, ja se kestää marraskuusta huhtikuuhun.
Talvella on paljon lunta ja pakkasta, mutta suomalaiset rakastavat talviurheilua.
Monet käyvät hiihtämässä ja luistelemassa.
Kesällä taas yöt ovat valoisia, ja ihmiset viettävät aikaa mökeillä järvien rannalla.`,
			Translation: `Finland has four seasons: spring, summer, autumn and winter.
Winter is the longest season, and it lasts from November to April.
In winter there is a lot of snow and frost, but Finns love winter sports.
Many go skiing and skating.
In summer, on the other hand, nights are bright, and people spend time at cottages by the lakeside.`,
			CEFRLevel: "B1",
			Topic:     topic,
		},
	}

	if story, ok := stories[level]; ok {
		return story
	}
	return stories["A1"]
}