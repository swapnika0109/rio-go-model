package tokens

import (
	"rio-go-model/configs"
	"rio-go-model/internal/util"

	"regexp"
	"unicode/utf8"

	"google.golang.org/genai"
)

type StoryCharacters struct {
	Characters      []string `json:"characters"`
	AudioCharacters []string `json:"audio_characters"`
	logger          *util.CustomLogger
}

func NewStoryCharacters() *StoryCharacters {
	logger := util.GetLogger("story.characters", configs.GetSettings())
	return &StoryCharacters{
		logger: logger,
	}
}

func (sc *StoryCharacters) GetGoogleStoryTokens(model string, resp *genai.GenerateContentResponse) int32 {
	sc.logger.Infof("Gemini usage for model=%s", model)

	if resp.UsageMetadata == nil {
		sc.logger.Warnf("UsageMetadata not available on response")
		return 0
	}

	totalTokens := resp.UsageMetadata.TotalTokenCount

	if totalTokens > 0 {
		sc.logger.Infof("Total tokens=%d", totalTokens)
	} else {
		sc.logger.Warnf("Total tokens not available")
	}

	return totalTokens
}

// CountAudioChars computes chargeable character counts for TTS requests.
// For SSML, tags are stripped before counting visible runes.
func (sc *StoryCharacters) CountAudioChars(text string, ssml string) int {
	if ssml != "" {
		tagRe := regexp.MustCompile(`<[^>]+>`)
		visible := tagRe.ReplaceAllString(ssml, "")
		chars := utf8.RuneCountInString(visible)
		sc.logger.Infof("TTS usage (SSML): visibleChars=%d rawBytes=%d", chars, len(ssml))
		return chars
	}
	chars := utf8.RuneCountInString(text)
	sc.logger.Infof("TTS usage (text): chars=%d", chars)
	return chars
}
