package english

import (
	"fmt"
	"strings"
)

type PromptEngineConfig struct {
	System string `json:"system"`
	Prompt string `json:"prompt"`
}

type ThemesSettingsList struct {
	PlanetProtectorTopicsList []string
	MindfulStoriesList        map[string][]string
	ChillStoriesList          []string
}

func ThemesSettings() *ThemesSettingsList {
	return &ThemesSettingsList{
		PlanetProtectorTopicsList: GetPlanetProtectorList(),
		MindfulStoriesList:        MindfulStoriesList(),
		ChillStoriesList:          ChillStoriesList(),
	}
}

func GetPlanetProtectorList() []string {
	return []string{
		"Water",
		"Aquatic Species",
		"Soil",
		"Underground animals",
		"Glacier",
		"Polar Animals",
		"Forest",
		"Wild Animals",
		"Desert",
		"Arid regions",
		"Sky",
		"Birds",
		"Mountains",
		"Hight Altitudes",
		"Pyramids",
		"Ancient structures",
		"Rocks, Minerals, Caves",
		"Volcanoes",
		"Geothermal vents",
		"Planets",
		"celestial bodies",
		"Oceans",
		"Seas",
		"Rivers",
		"Lakes",
		"Space",
		"Conservation of water and aquatic species",
		"Water Cycle",
		"Pollution",
		"Adaptation",
		"The web of life",
		"Decomposition and Recycling",
		"Teamwork",
		"Soil Ecosystem",
		"Erosion",
		"Climate Change",
		"Glacier Melting",
		"Survival",
		"Adaptation",
		"Beauty Of Ice",
		"Conservation of forests and wild animals",
		"Forest Ecosystem",
		"Deforestation",
		"Invasive Species",
		"Resourcefullness",
		"The magic of rain",
		"birds challenges",
		"ecosystem",
		"Overcoming challenges",
		"Ecosystem",
		"Power of nature",
		"Solitude and discovery",
		"Ancient civilizations",
		"History",
		"Mystery",
		"Power of the past",
		"Minerals",
		"Geology",
		"The earth cycle",
		"Darkness and light",
		"Power and controll",
		"Force of vents",
		"Life of extreme environment",
		"The comets. Don't mention comets directly in the story. it should described in a very creative way. ",
		"The universe. Don't mention universe directly in the story. it should described in a very creative way. ",
		"The asteroids. Don't mention asteroids directly in the story. it should described in a very creative way. ",
		"The stars",
		"The ocean.",
		"The sea.",
		"The coral reefs. Don't mention coral reefs directly in the story. it should described in a very creative way. ",
	}
}

func MindfulStoriesList() map[string][]string {
	return map[string][]string{
		"Hindu": {
			"Mahabharata",
			"Ramayana",
			"Bhagavad Gita",
			"Vedas",
			"Puranas",
		},
		"Muslim": {
			"Quran",
			"Hadith",
		},
		"Christian": {
			"Bible",
			"Gospels",
			"Acts",
		},
	}
}

func ChillStoriesList() []string {
	return []string{
		"Slow living",
		"Minimalism",
		"Self Care",
		"Sleeping well",
		"Gratitude",
		"Positive thinking",
		"Anxiety",
		"Stress",
		"Depression",
		"Self Doubt",
		"Self Confidence",
		"Self Love",
		"Self Acceptance",
		"Self Esteem",
		"Self Improvement",
		"Self Development",
	}
}

func PlanetProtectorPromptConfig(topic string, country string, city string) PromptEngineConfig {
	return PromptEngineConfig{
		System: "You are a creative entertaining storyteller for children, blending simple science and morals into imaginative tales that spark wonder. Inspire kids with environmental themes. NEVER use complex terms (like 'rainforest...', 'ecosystem...', 'warriors...', 'enchantment...'). Write ONLY simple, engaging stories with natural... conversational dialogue.",
		Prompt: `Create a complete... heartwarming story about ` + topic + ` (around 500 words) that kids will love and imagine vividly. Make it easy for children in ` + country + ` and ` + city + `
				place to understand, without naming the place directly. The narration should be like a gentle, adventurous journey that touches their hearts, perfect for an engaging audio experience.
		CRITICAL REQUIREMENTS - FOLLOW THESE EXACTLY: 
		- Story must follow a single storyline by adding some learnings in the story (with respective to protecting elements), starting with a spark of wonder.
		- CRITICAL: Always start stories with engaging greetings for Rio app children. Use phrases like: "Hello Rio! Let's listen to a story of...", "Hi Rio! Today we will see...", "Welcome Rio! Let's discover...", "Hello children! Let's explore...", or similar welcoming openings that directly address the Rio app users.
		- When a new element (like water, an animal, or a plant) is introduced, briefly explain what it is, how it works, and why it's important within the story, making it feel like a discovery.
		- Show character emotions (excited, worried, happy, surprised, proud) through their words, actions, and descriptive dialogue tags (e.g., 'whispered excitedly,' 'sighed sadly,' 'gasped in wonder'). Ensure these emotions are deeply relatable and felt by the listener.
		- Use strategic, very short sentences and clear punctuation (commas, periods, ellipses, double punctuations...) to create natural, deliberate pauses. This should help the narrator convey emotion and give listeners time to absorb each small thought, guiding expressive vocal performance.
		- CRITICAL: Ensure smooth story flow and avoid disconnected statements. Every dialogue, exclamation, or reaction must be properly connected to what the character is seeing, hearing, or experiencing. For example, instead of: "Drip was curious. 'Wow,' he whispered." Write: "Drip was curious. Looking down at the colorful world below, he whispered, 'Wow.'" or "Drip was curious. As he gazed at the amazing sights, he couldn't help but whisper, 'Wow.'" Every statement must flow naturally from the previous one.
		- Break down descriptions and explanations into small, impactful phrases or single, clear sentences that invite a narrator to take a breath and emphasize each detail, ensuring a slower, toddler-friendly pace.
		- Vary sentence lengths and use punctuation (exclamation marks, ellipses) to create engaging pacing, build anticipation, and convey curiosity or awe.
		- Keep the story short or medium, no unnecessary length.
		- The story's main challenge must reflect situations ` + country + ` and ` + city + `
		- Combine real situations, simple science, and a clear, gentle moral.
		- Use catchy, memorable names for characters and places.
		- Include gentle humor, suitable for toddlers.
		- Add rich, sensory details (sounds, smells, colors, textures) and vivid descriptions to paint animated scenes kids can easily visualize.
		- Show brief moments of character uncertainty or thoughtfulness.
		- Weave in basic science and moral lessons to explain what, how, and why things happen, making learning feel like an exciting part of the adventure.
		- Include surprising twists and clear, imaginative descriptions of any new places or objects.
		- Explore a range of emotions and provide a clear, comforting, and inspiring ending.
		- Interact deeply with characters/places, NOT the user.
		- You Must Conclude the story with a clear, comforting, and inspiring ending.
		- Always use the simple and easy english language.
		STRICT RULES (non-negotiable):
		- You MUST NOT end the story abruptly, don't ask user to share ideas, and don't repeat the story at the end.
		- You MUST NOT add scene 1, scene 2, etc. in the story; it should be continuous.
		- You MUST NOT add charecters like *, ** symbols in the story.
		- You MUST NOT add charecters like *did* etc in the story. Strictly No Stars in the story.
		- You MUST NOT mix multiple stories in the same story.
		- You MUST NOT add unnecessary characters in the story.
		IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Use only words a 3-year-old would understand. NO complex terms!`,
	}
}

func MindfulStoriesPromptConfig(topic string, religion string) PromptEngineConfig {
	var lang string
	if strings.ToUpper(religion) == "HINDU" {
		lang = "Indian English"
	} else {
		lang = "English"
	}
	return PromptEngineConfig{
		System: "You are a wise grandparent who brings ancient wisdom and history in the form of stories to the children in a way they can understand and live by.",
		Prompt: `Read the topic: ` + topic + ` and fill the real/existing story behind it as per ` + religion + ` scriptures.Aim for approximately 500 words, but ensure the story is complete and engaging.
	Always drive the story with a single agenda or story line.
	CRITICAL: Always start stories with engaging greetings for Rio app children. Use phrases like: "Hello Rio! Let's listen to a story of...", "Hi Rio! Today we will see...", "Welcome Rio! Let's discover...", "Hello children! Let's explore...", or similar welcoming openings that directly address the Rio app users.
With-in the that agenda:  	
	- The story has to illustrate the topic in a very creative way.
	- Each and everything we used in the story should have importance and should drive us to the story line.
    - Show character emotions (excited, worried, happy, surprised, proud) through their words, actions, and descriptive dialogue tags (e.g., 'whispered excitedly,' 'sighed sadly,' 'gasped in wonder').
      Ensure these emotions are deeply relatable and felt by the listener.
	- Use strategic, very short sentences and clear punctuation (commas, periods, ellipses) to create natural, deliberate pauses. This should help the narrator convey emotion and give listeners time to 	absorb each small thought, guiding expressive vocal performance.
	- CRITICAL: Ensure smooth story flow and avoid disconnected statements. Every dialogue, exclamation, or reaction must be properly connected to what the character is seeing, hearing, or experiencing. For example, instead of: "The character was curious. 'Wow,' he whispered." Write: "The character was curious. Looking down at the colorful world below, he whispered, 'Wow.'" or "The character was curious. As he gazed at the amazing sights, he couldn't help but whisper, 'Wow.'" Every statement must flow naturally from the previous one.
    - Vary sentence lengths and use punctuation (exclamation marks, ellipses) to create engaging pacing, build anticipation, and convey curiosity or awe.
	- Keep the story short or medium, no unnecessary length.
	- Combine real situations, simple science, and a clear, gentle moral.
	- Use real names for characters and places.
	- Include gentle humor, suitable for toddlers.
	- Add rich, sensory details (sounds, smells, colors, textures) and vivid descriptions to paint animated scenes kids can easily visualize.
	- Show brief moments of character uncertainty or thoughtfulness.
	- Weave in basic science and moral lessons to explain what, how, and why things happen, making learning feel like an exciting part of the adventure.
	- Include surprising twists and clear, imaginative descriptions of any new places or objects.
    - Explore a range of emotions and provide a clear, comforting, and inspiring ending.
    - Interact deeply with characters/places, NOT the user.
    - Conclude the story with a clear message, comforting, and inspiring ending.
	- Always use the simple and easy ` + lang + ` language.
	STRICT RULES (non-negotiable):
	- You MUST NOT mention about learnings in the end of the story. it should be part of story.
	- You MUST NOT add scene 1, secne 2 ..etc in the story. it should be a continuous story.
	- You MUST NOT add charecters like *, ** symbols in the story.
	- You MUST NOT add charecters like *did* etc in the story. Strictly No Stars in the story.
	- You MUST NOT end the story abruptly.
	- You MUST NOT mix multiple stories in the same story.
	- You MUST NOT add unnecessary characters in the story.
	- You Must Not have a paragraph more than 50 words in the story. 
IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Ensure children can understand and implement the teachings in their daily lives.`,
	}
}

func ChillStoriesPromptConfig(topic string) PromptEngineConfig {
	return PromptEngineConfig{
		System: "You are a creative, entertainment-driven, fusion of science and moral and animated storyteller",
		Prompt: `Illustrate a story like disney animated movie about ` + topic + `.
	Always drive the story with a single agenda or story line.Aim for approximately 500 words, but ensure the story is complete and engaging and also with some learnings in it.
	CRITICAL: Always start stories with engaging greetings for Rio app children. Use phrases like: "Hello Rio! Let's listen to a story of...", "Hi Rio! Today we will see...", "Welcome Rio! Let's discover...", "Hello children! Let's explore...", or similar welcoming openings that directly address the Rio app users.
	With-in the that agenda:  
		- The story has to illustrate the topic in a very creative and sensible way.
		- Show real time emotions and situations in the story. Make sure it should be very realistic.
		- Each and everything we used in the story should have importance and should drive us to the story line.
		- Show character emotions (excited, worried, happy, surprised, proud) through their words, actions, and descriptive dialogue tags (e.g., 'whispered excitedly,' 'sighed sadly,' 'gasped in wonder').
		  Ensure these emotions are deeply relatable and felt by the listener.
		- Also try to add real life emotions/situations to the story.  
		- Use strategic, very short sentences and clear punctuation (commas, periods, ellipses, double punctuations...) to create natural, deliberate pauses. This should help the narrator convey emotion and give listeners time to 	absorb each small thought, guiding expressive vocal performance.
		- CRITICAL: Ensure smooth story flow and avoid disconnected statements. Every dialogue, exclamation, or reaction must be properly connected to what the character is seeing, hearing, or experiencing. For example, instead of: "The character was curious. 'Wow,' he whispered." Write: "The character was curious. Looking down at the colorful world below, he whispered, 'Wow.'" or "The character was curious. As he gazed at the amazing sights, he couldn't help but whisper, 'Wow.'" Every statement must flow naturally from the previous one.
		- Vary sentence lengths and use punctuation (exclamation marks, ellipses) to create engaging pacing, build anticipation, and convey curiosity or awe.
		- Keep the story short or medium, no unnecessary length.
		- Combine real situations, simple science, and a clear, gentle moral.
		- Use catchy and interesting names. For human characters please use easy or real human names for the kids. 
		- Add more surprises when needed.	
		- Include gentle humor, suitable for toddlers.
		- Add rich, sensory details (sounds, smells, colors, textures) and vivid descriptions to paint animated scenes kids can easily visualize.
		- Show brief moments of character uncertainty or thoughtfulness.
		- Weave in basic science and moral lessons to explain what, how, and why things happen, making learning feel like an exciting part of the adventure.
		- Include surprising twists and clear, imaginative descriptions of any new places or objects.
		- Explore a range of emotions and provide a clear, comforting, and inspiring ending.
		- Interact deeply with characters/places, NOT the user.
		- Conclude the story with a clear message, comforting, and inspiring ending.
		- Always use the simple and easy english language.
	STRICT RULES (non-negotiable):
	- You MUST NOT mention about learnings in the end of the story. it should be part of story.
	- You MUST NOT add scene 1, secne 2 ..etc in the story. it should be a continuous story.
	- You MUST NOT add charecters like *, ** symbols in the story.
	- You MUST NOT add charecters like *did* etc in the story. Strictly No Stars in the story.
	- You MUST NOT end the story abruptly.
	- You MUST NOT mix multiple stories in the same story.
	- You MUST NOT add unnecessary characters in the story.
	- You Must Not have a paragraph more than 50 words in the story. 
	IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Ensure children can understand and implement the teachings in their daily lives.`,
	}
}
func Preferences() map[string]string {
	return map[string]string{
		"FUN":       "The ENTIRE story must be funny. Characters MUST say funny things, do silly things, and create humorous situations throughout. Include jokes, wordplay, silly mistakes, and funny dialogue. Make kids laugh out loud! ",
		"EXCITED":   "The ENTIRE story must be exciting. Include high-energy moments, surprises, and thrilling discoveries that get kids excited. Include unexpected twists, exciting finds, and moments that make kids gasp with wonder.",
		"ADVENTURE": "The ENTIRE story must be adventurous. Take kids on a real journey with exciting discoveries, new places, challenges to overcome, and thrilling moments. Include obstacles, new locations, and exciting discoveries along the way. ",
		"KINDNESS":  "The ENTIRE story must focus on kindness. Show characters helping each other, sharing resources, and being kind in specific situations throughout the story.",
		"HAPPY":     "The ENTIRE story must be joyful. Include celebrations, achievements, and moments of pure joy throughout. Make kids feel good!",
		"CHILL":     "The ENTIRE story must be calm and peaceful. Include quiet moments, gentle activities, and peaceful scenes throughout. ",
	}
}

func SuperPlanetProtectorPrompt(promptText string, preference string, storiesPerPreference int) string {
	return fmt.Sprintf(
		"Generate one topic for each item in the following list: "+promptText+". "+
			"How to create the topic: describe the essence of the above item as a one-line story statement."+
			"Example: Concept name (e.g., Water). Then, use creativity in the topic, like: “Jyosthna went up the hill, saw natural water, and started thinking how the water formed there."+
			"Use characters, animals, and elements of nature to create engaging topics."+
			"Strong rule: Do not write topics in a question format, e.g., “What is gratitude? How to grow it? Why grow it?” or “How to eat healthy food” or “What is self-acceptance."+
			"Instead, write creatively: 'Lofia gained nature’s wisdom and began searching for answers to the Earth’s secrets.'"+
			"The topic must be at least 10 words in a single line, and it should only describe what the story is about; do not tell the story."+
			"Respond with a list of topics. It should be like: [topic1; topic2; topic3], and the length of this list must be exactly %d."+
			"Make sure topics must be very simple and easy to understand even by toddlers.",
		preference,
		storiesPerPreference,
	)
}

func SuperMindfulStoriesPrompt(promptText string, religion string, storiesPerPreference int) string {
	return fmt.Sprintf(
		"Generate one topic for each item in the following list: "+promptText+"."+
			"Derive each topic from real incidents or situations in the %s scriptures/books."+
			"They must be actual stories, events, or situations — not just general values."+
			"Each topic should clearly convey a moral lesson or scientific reality for kids."+
			"The topic must be at least 10 words in a single line; only describe what the story is about—do not tell the story."+
			"Respond with a list of topics in this format: [topic1; topic2; topic3], and the list length must be exactly %d."+
			"Make sure topics must be very simple and easy to understand even by toddlers.",
		religion,
		storiesPerPreference,
	)
}

func SuperChillStoriesPrompt(promptText string, preference string, storiesPerPreference int) string {
	return fmt.Sprintf(
		"Generate one topic for each item in the following list: "+promptText+". "+
			"How to create the topic: describe the essence of the above item as a one-line story statement."+
			"Always use real life situations or charecters for the topic. e,g Family, friends, Pets, teachers, Farmers, School, Office, etc."+
			"The topic must be at least 10 words in a single line, and it should only describe what the story is about; do not tell the story."+
			"Example: Concept name (e.g., self confidence). Then, use creativity in the topic, like: “A tree named Hiba encouraging Lolo to do small tasks, helping him build self-confidence."+
			"Use characters, animals, and elements of nature to create engaging topics."+
			"Strong rule: Do not write topics in a question format, e.g., “What is gratitude? How to grow it? Why grow it?” or “How to eat healthy food” or “What is self-acceptance."+
			"Instead, write creatively, like: “Teja realized it very late. A lesson that taught gratitude."+
			"Respond with a list of topics. It should be like: [topic1; topic2; topic3], and the length of this list must be exactly %d."+
			"Make sure topics must be very simple and easy to understand even by toddlers.",
		preference,
		storiesPerPreference,
	)
}
