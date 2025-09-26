package configs

import (
	"errors"
	"math/rand"
	"time"
)

type PromptEngineConfig struct {
	System string `json:"system"`
	Prompt string `json:"prompt"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFrom returns a random index from the provided list length.
func RandomFrom(list []string) (int, error) {
	if len(list) == 0 {
		return 0, errors.New("list is empty")
	}
	return rand.Intn(len(list)), nil
}

type ThemesSettingsList struct {
	PlanetProtectorTopicsList []string
	MindfulStoriesList map[string][]string
	ChillStoriesList []string
}

func ThemesSettings() *ThemesSettingsList {
	return &ThemesSettingsList{
		PlanetProtectorTopicsList: GetPlanetProtectorList(),
		MindfulStoriesList: MindfulStoriesList(),
		ChillStoriesList: ChillStoriesList(),
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
		"The comets",
		"The universe",
		"The asteroids",
		"The stars",
		"The ocean",
		"The sea",
		"The coral reefs",
	}
}

func MindfulStoriesList() map[string][]string {
	return map[string][]string{
		"Hindu":{
			"Mahabharata",
			"Ramayana",
			"Bhagavad Gita",
			"Vedas",
			"Puranas",
			"Upanishads",
		},
		"Muslim":{
			"Quran",
			"Hadith",
		},
		"Christian":{
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
		"Eating healthy",
		"Sleeping well",
		"Meditation",
		"Yoga",
		"Gratitude",
		"Positive thinking",
		"Anxiety",
		"Stress",
		"Depression",
		"Self Doubt",
		"Self Confidence",
		"Self Love",
		"Self Acceptance",
	}
}

func PlanetProtectorPromptConfig(topic string, country string, city string) PromptEngineConfig {
	return PromptEngineConfig{
		System: "You are a creative entertainment-driven , fusion of science and moral and animated imaginative storyteller who weaves magical tales that inspire children to think innovatively about environmental themes. NEVER use complex terms like 'rainforest', 'ecosystem', 'warriors', or 'enchantment'. Write ONLY simple, engaging stories with natural dialogue.",
		Prompt: `Create a VERY ELABORATE and enchanting story about ` + topic + ` that can be easily understandable by people staying in ` + country + ` and ` + city + `. but dont use country and city directly in the story.
		CRITICAL REQUIREMENTS - FOLLOW THESE EXACTLY: 
		- Always drive the story by choosing a challenge based on the the real time situations at ` + country + ` and ` + city + `.
		- The story should always be based on the topic with a litle fusion of real time situations, science and moral.
		- Always use catchy names for the characters and places.
		- Generat only Medium size story.
		- The story should include humour and should be understandable by kids of 3-5 years old
		- Write this as a VERY ELABORATE STORY with rich details and have at least 8 to 10 illustrated/animated detailed scenes with vivid descriptions
		- Add interactivity, challenges & choices to the story by having deep character development
		- Make the story non-linear with more opportunities for kids to interact, pick solutions or answer questions
		- Include character emotions, thoughts, and reactions throughout by having natural dialogue that sounds like real time conversations
		- Use catchy names for kids to understand and imagine
		- Whenever needed Add rich sensory details (sounds, smells, colors, textures, tastes) and support illustration/animation
		- Add brief moments of character when they are having uncertain emotions
		- Add more educational elements, STEM learning, nature learning, scientific learning etc.
		- Include some basic science and moral in the story. so that they can learn why, how and what is happening in the story.
		- Whenever needed Add surprising twists and discoveries and illustration too
		- Add clear descriptions of any new places or objects to make kids imagine like an animation movie
		- Explore more emotions & ending resolution
		- The story has to interact more deeply with characters/places but not with the user.
		- Don't end the story abruptly, don't ask user to share ideas. and also don't repeat the story at the end.
		- Don't add scene 1, secne 2 ..etc in the story. it should be a continuous story.
		- Don't add ** symbols in the story.
		- Don't use country, city directly names in the story.
		IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Use only words a 3-year-old would understand. NO complex terms!`,
	}
}

func MindfulStoriesPromptConfig(topic string, religion string) PromptEngineConfig {
	return PromptEngineConfig{
		System: "You are a wise grandparent who brings ancient wisdom and history in the form of stories to the children in a way they can understand and live by.",
		Prompt: `Read the topic: ` + topic + ` and fill the real/existing story behind it as per ` + religion + ` scriptures.
Add more details, emotions, and interactions to the story.
Add more illustrations and animations and fusion of science and moral and make the story more engaging and interactive and understandable for the kids of 3-9 years old
Illustrate the story with more educational elements, STEM learning, nature learning, scientific learning etc. these learnings should be part of story.
Kids should learn the story by understanding the science and moral in it.
Add more surprises and discoveries to the story.
Add more interactions with characters and places to the story.
- Whenever needed Add surprising twists and discoveries and illustration too
- The story has to interact more deeply with characters/places but not with the user.
- Don't mention about learnings in the end of the story. it should be part of story.
- Don't add scene 1, secne 2 ..etc in the story. it should be a continuous story.
- Don't add ** symbols in the story.
- Don't end the story abruptly.
IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Ensure children can understand and implement the teachings in their daily lives.`,
	}
}

func ChillStoriesPromptConfig(topic string) PromptEngineConfig {
	return PromptEngineConfig {
		System: "You are a creative, entertainment-driven and animated storyteller",
		Prompt: `Illustrate a story like disney animated movie about ` + topic + `, by explaining how important it is to live the life slowly and observe what life is giving us with an example related to the story.
Also explain the science and moral behind the story by adding more details like multiple scenes with more interactions having beautiful emotions
Each scene should be very engaging and give surprising illustrations and animations
Use catchy and interesting names. For human characters please use easy or real human names for the kids. 
Add more surprises when needed.
Make the story more engaging and interactive and understandable for the kids of 3-9 years old
Kids should learn the story by understanding the science and moral in it.
Add more interactions with characters and places to the story.
- The story has to interact more deeply with characters/places but not with the user.
- Don't mention about learnings in the end of the story. it should be part of story.
- Don't end the story abruptly.
- Don't add scene 1, secne 2 ..etc in the story. it should be a continuous story.
- Don't add ** symbols in the story.
- Don't end the story abruptly.
IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Ensure children can understand and implement the teachings in their daily lives.`,
	}
}
