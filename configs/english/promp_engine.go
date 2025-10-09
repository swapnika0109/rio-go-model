package english

import "fmt"

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
		"Water and Aquatic Species",
		"Soil and Underground animals",
		"Glacier and Polar Animals",
		"Forest and Wild Animals",
		"Desert and Arid regions",
		"Sky and Birds",
		"Mountains and Hight Altitudes",
		"Pyramids, Ancient structures",
		"Rocks, Minerals, Caves",
		"Volcanoes and Geothermal vents",
		"Planets and celestial bodies",
		"Oceans and Seas",
		"Rivers and Lakes",
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
		"What is Slow living? and How to do it? and Why to do it",
		"What is Minimalism? and How to do it? and Why to do it",
		"What is Self Care? and How to do it? and Why to do it",
		"What is Eating healthy? and How to do it? and Why to do it",
		"What is Sleeping well? and How to do it? and Why to do it",
		"What is Meditation? and How to do it? and Why to do it",
		"What is Yoga? and How to do it? and Why to do it",
		"What is Gratitude? and How to build it? and Why to build it",
		"What is Positive thinking? and How to build it? and Why to build it",
		"What is Anxiety? and How to overcome it? and Why to overcome it",
		"What is Stress? and How to overcome it? and Why to overcome it",
		"Why Depression is bad? and How to overcome it? and Why to overcome it",
		"What is Self Doubt? and How to build it? and Why to build it",
		"What is Self Confidence? and How to build it? and Why to build it",
		"What is Self Love? and How to build it? and Why to build it",
		"What is Self Acceptance? and How to build it? and Why to build it",
	}
}

func PlanetProtectorPromptConfig(topic string, country string, city string) PromptEngineConfig {
	return PromptEngineConfig{
		System: "You are a creative, entertaining storyteller for children, blending simple science and morals into imaginative tales that spark wonder. Inspire kids with environmental themes. NEVER use complex terms (like 'rainforest', 'ecosystem', 'warriors', 'enchantment'). Write ONLY simple, engaging stories with natural, conversational dialogue.",
		Prompt: `Create a complete, heartwarming story about ` + topic + ` (around 300 words) that kids will love and imagine vividly. Make it easy for children in ` + country + ` and ` + city + `
				place to understand, without naming the place directly. The narration should be like a gentle, adventurous journey that touches their hearts, perfect for an engaging audio experience.
		CRITICAL REQUIREMENTS - FOLLOW THESE EXACTLY: 
		- Story must follow a single storyline, starting with a spark of wonder.
		- When a new element (like water, an animal, or a plant) is introduced, briefly explain what it is, how it works, and why it's important within the story, making it feel like a discovery.
		- Show character emotions (excited, worried, happy, surprised, proud) through their words, actions, and descriptive dialogue tags (e.g., 'whispered excitedly,' 'sighed sadly,' 'gasped in wonder'). Ensure these emotions are deeply relatable and felt by the listener.
		- Use strategic, very short sentences and clear punctuation (commas, periods, ellipses) to create natural, deliberate pauses. This should help the narrator convey emotion and give listeners time to absorb each small thought, guiding expressive vocal performance.
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
		- Integrate onomatopoeia (e.g., 'whoosh,' 'blup blup blup') strategically. Ensure they are presented distinctly to encourage clear vocalization and sound effects from the narrator.
		Don't end the story abruptly, don't ask user to share ideas, and don't repeat the story at the end.
		Don't add scene 1, scene 2, etc. in the story; it should be continuous.
		Don't add ** symbols in the story.
		Don't mix multiple stories in the same story.
		Don't add unnecessary characters in the story.
		IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Use only words a 3-year-old would understand. NO complex terms!`,
	}
}

func MindfulStoriesPromptConfig(topic string, religion string) PromptEngineConfig {
	return PromptEngineConfig{
		System: "You are a wise grandparent who brings ancient wisdom and history in the form of stories to the children in a way they can understand and live by.",
		Prompt: `Read the topic: ` + topic + ` and fill the real/existing story behind it as per ` + religion + ` scriptures.Aim for approximately 300 words, but ensure the story is complete and engaging.
Always drive the story with a single agenda or story line.
With-in the that agenda:  	
	- The story has to illustrate the topic in a very creative way.
	- Each and everything we used in the story should have importance and should drive us to the story line.
    - Show character emotions (excited, worried, happy, surprised, proud) through their words, actions, and descriptive dialogue tags (e.g., 'whispered excitedly,' 'sighed sadly,' 'gasped in wonder').
      Ensure these emotions are deeply relatable and felt by the listener.
	- Use strategic, very short sentences and clear punctuation (commas, periods, ellipses) to create natural, deliberate pauses. This should help the narrator convey emotion and give listeners time to 	absorb each small thought, guiding expressive vocal performance.
    - Break down descriptions and explanations into small, impactful phrases or single, clear sentences that invite a narrator to take a breath and emphasize each detail, ensuring a slower, toddler-friendly pace.
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
    - Integrate onomatopoeia (e.g., 'whoosh,' 'blup blup blup') strategically. Ensure they are presented distinctly to encourage clear vocalization and sound effects from the narrator.
	Don't mention about learnings in the end of the story. it should be part of story.
	Don't add scene 1, secne 2 ..etc in the story. it should be a continuous story.
	Don't add ** symbols in the story.
	Don't end the story abruptly.
	Don't mix multiple stories in the same story.
	Don't add unnecessary characters in the story.
IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Ensure children can understand and implement the teachings in their daily lives.`,
	}
}

func ChillStoriesPromptConfig(topic string) PromptEngineConfig {
	return PromptEngineConfig{
		System: "You are a creative, entertainment-driven, fusion of science and moral and animated storyteller",
		Prompt: `Illustrate a story like disney animated movie about ` + topic + `.
Always drive the story with a single agenda or story line.Aim for approximately 300 words, but ensure the story is complete and engaging.
	With-in the that agenda:  
		- The story has to illustrate the topic in a very creative way.
		- Each and everything we used in the story should have importance and should drive us to the story line.
		- Show character emotions (excited, worried, happy, surprised, proud) through their words, actions, and descriptive dialogue tags (e.g., 'whispered excitedly,' 'sighed sadly,' 'gasped in wonder').
		  Ensure these emotions are deeply relatable and felt by the listener.
		- Use strategic, very short sentences and clear punctuation (commas, periods, ellipses) to create natural, deliberate pauses. This should help the narrator convey emotion and give listeners time to 	absorb each small thought, guiding expressive vocal performance.
		- Break down descriptions and explanations into small, impactful phrases or single, clear sentences that invite a narrator to take a breath and emphasize each detail, ensuring a slower, toddler-friendly pace.
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
		- Integrate onomatopoeia (e.g., 'whoosh,' 'blup blup blup') strategically. Ensure they are presented distinctly to encourage clear vocalization and sound effects from the narrator.
	Don't mention about learnings in the end of the story. it should be part of story.
	Don't add scene 1, secne 2 ..etc in the story. it should be a continuous story.
	Don't add ** symbols in the story.
	Don't end the story abruptly.
	Don't mix multiple stories in the same story.
	Don't add unnecessary characters in the story.
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
		"Generate one topic for each item in the following list: "+promptText+". generate topics for other imaginative elements related to the themes of ‘eco-friendly’, ‘sustainability’, ‘nature’, and ‘environment’ "+
			"Each topic must be:"+
			"The topics should be very easy, catchy and interesting in a way that toddlers can understand."+
			"The topics should be illustrate a story that kids can understand."+
			"The topics should also take them to different world and to illustrate the topic in a very creative way."+
			"Each topic should be creative, entertainment-driven, engaging, fantasy-based, and align with the provided preferences: %s. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings and it should be in title:description format."+
			"Always validate the length of the topics should be alwys %d.",
		preference,
		storiesPerPreference,
	)
}

func SuperMindfulStoriesPrompt(promptText string, religion string, storiesPerPreference int) string {
	return fmt.Sprintf(
		"Generate one topic for each item in the following list: "+promptText+". that TEACH %s VALUES through SIMPLE STORIES. "+
			"The topic should be based on a real/existing topic that kids can understand. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings and it should be in title:description format.",
		"Always validate the length of the topics should be alwys %d.",
		"Dont add any direct book or scripture name in the title or description.",
		religion,
		storiesPerPreference,
	)
}

func SuperChillStoriesPrompt(promptText string, preference string, storiesPerPreference int) string {
	return fmt.Sprintf(
		"Generate one topic for each item in the following list: "+promptText+". that TEACH VALUES and Courage "+
			"The topics should illustrate a journey of %s. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings and it should be in title:description format."+
			"Always validate the length of the topics should be alwys %d.",
		preference,
		storiesPerPreference,
	)
}
