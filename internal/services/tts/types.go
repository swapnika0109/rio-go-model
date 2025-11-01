package tts

type VoiceType string
type StoryType string

const (
	Chirp3HD VoiceType = "Chirp3"
	// WaveNet  VoiceType = "WaveNet"
	Standard VoiceType = "Standard"
)

const (
	StoryPremium StoryType = "Premium"
	StoryFree    StoryType = "Standard"
)

func (v VoiceType) String() string {
	return string(v)
}

func (s StoryType) String() string {
	return string(s)
}
