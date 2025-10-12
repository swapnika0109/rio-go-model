package util

import (
	"regexp"
	"strings"
)

// TeluguSSMLBuilder creates SSML for Telugu text with proper emotion, pauses, and stress
type TeluguSSMLBuilder struct {
	// Telugu vowel signs (vothulu) for stress detection
	vothulu map[rune]bool
	// Emotional keywords in Telugu
	emotionalKeywords map[string]string
	// Punctuation to emotion mapping
	punctuationEmotions map[string]string
}

// NewTeluguSSMLBuilder creates a new Telugu SSML builder
func NewTeluguSSMLBuilder() *TeluguSSMLBuilder {
	return &TeluguSSMLBuilder{
		vothulu: map[rune]bool{
			// Telugu vowel signs (vothulu) - all 16 vowel signs
			'ా': true, 'ి': true, 'ీ': true, 'ు': true, 'ూ': true, 'ె': true, 'ే': true, 'ొ': true, 'ో': true,
			'ౌ': true, 'అ': true, 'ఆ': true, 'ఇ': true, 'ఈ': true, 'ఉ': true, 'ఊ': true, 'ఎ': true, 'ఏ': true,
			'ఒ': true, 'ఓ': true, 'ఔ': true, 'ం': true, 'ః': true,

			// Additional Telugu vowel signs and diacritical marks
			'ృ': true, 'ౄ': true, // R-vowels (ఋ, ౠ)
			'ౢ': true, 'ౣ': true, // L-vowels (ఌ, ౡ)

			// Telugu diacritical marks and modifiers
			'్': true,                       // Virama (halant) - removes inherent vowel
			'ఁ': true,                       // Candrabindu (anusvara)
			'ౘ': true, 'ౙ': true, 'ౚ': true, // Additional Telugu characters

			// Telugu consonants (vyanjanalu) - all 36 consonants
			'క': true, 'ఖ': true, 'గ': true, 'ఘ': true, 'ఙ': true,
			'చ': true, 'ఛ': true, 'జ': true, 'ఝ': true, 'ఞ': true,
			'ట': true, 'ఠ': true, 'డ': true, 'ఢ': true, 'ణ': true,
			'త': true, 'థ': true, 'ద': true, 'ధ': true, 'న': true,
			'ప': true, 'ఫ': true, 'బ': true, 'భ': true, 'మ': true,
			'య': true, 'ర': true, 'ల': true, 'వ': true, 'శ': true,
			'ష': true, 'స': true, 'హ': true, 'ళ': true, 'ఱ': true,

			// Telugu numbers (can have stress too)
			'౦': true, '౧': true, '౨': true, '౩': true, '౪': true,
			'౫': true, '౬': true, '౭': true, '౮': true, '౯': true,
		},
		emotionalKeywords: map[string]string{
			// Positive emotions
			"ఆనందం": "strong", "సంతోషం": "strong", "ఆశ్చర్యం": "moderate", "ఉత్సాహం": "strong",
			"ప్రేమ": "moderate", "ఆశ": "moderate", "విజయం": "strong", "గర్వం": "moderate",
			"ఆశీర్వాదం": "moderate", "ఆనందించు": "moderate", "ఆశ్చర్యపడు": "moderate",

			// Negative emotions
			"భయం": "moderate", "దుఃఖం": "moderate", "కోపం": "strong", "అసహనం": "moderate",
			"భయపడు": "moderate", "ఏడ్చు": "moderate", "కోపగించు": "strong", "అసహనపడు": "moderate",
			"భయంకర": "strong", "భయానక": "strong", "భయపెట్టు": "moderate",

			// Action/Intensity words
			"వేగంగా": "moderate", "నెమ్మదిగా": "moderate", "ధైర్యంగా": "moderate", "భయంకరంగా": "strong",
			"ఆశ్చర్యకరంగా": "moderate", "అద్భుతంగా": "strong", "అత్యంత": "moderate",

			// Story elements
			"రాజు": "moderate", "రాణి": "moderate", "యువరాజు": "moderate", "యువరాణి": "moderate",
			"మంత్రి": "moderate", "సేనాధిపతి": "moderate", "యోధుడు": "moderate", "యోధురాలు": "moderate",
			"మంత్రుడు": "moderate", "సాధువు": "moderate", "సాధ్వి": "moderate",
		},
		punctuationEmotions: map[string]string{
			"!":   "excited",      // Exclamation - excitement, surprise
			"?":   "questioning",  // Question - curiosity, doubt
			"...": "suspense",     // Ellipsis - suspense, pause
			"!!":  "very_excited", // Double exclamation - very excited
			"?!":  "surprised",    // Question-exclamation - surprise
		},
	}
}

// BuildTeluguSSML converts Telugu story text into SSML with proper emotion, pauses, and stress
func (t *TeluguSSMLBuilder) BuildTeluguSSML(story string) string {
	story = strings.TrimSpace(story)
	if story == "" {
		return "<speak></speak>"
	}

	// XML escape function - but preserve quotes for SSML
	escapeXML := func(s string) string {
		s = strings.ReplaceAll(s, "&", "&amp;")
		s = strings.ReplaceAll(s, "<", "&lt;")
		s = strings.ReplaceAll(s, ">", "&gt;")
		// Don't escape quotes - they're needed for SSML
		// s = strings.ReplaceAll(s, "\"", "&quot;")
		// s = strings.ReplaceAll(s, "'", "&apos;")
		return s
	}

	// Normalize whitespace
	normalize := func(s string) string {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		s = strings.ReplaceAll(s, "\r", "\n")
		spaceRe := regexp.MustCompile(`\s+`)
		s = spaceRe.ReplaceAllString(s, " ")
		paraRe := regexp.MustCompile(`\n{2,}`)
		s = paraRe.ReplaceAllString(s, "\n\n")
		return strings.TrimSpace(s)
	}

	story = normalize(story)
	paragraphs := strings.Split(story, "\n\n")

	var b strings.Builder
	b.WriteString("<speak>")

	for pi, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		p = escapeXML(p)
		b.WriteString("<p>")

		// Process sentences with emotion detection
		processedText := t.processTextWithEmotions(p)
		b.WriteString(processedText)

		b.WriteString("</p>")

		// Longer pause between paragraphs
		if pi < len(paragraphs)-1 {
			b.WriteString(`<break time="1s"/>`)
		}
	}

	b.WriteString("</speak>")
	return b.String()
}

// processTextWithEmotions processes text to add emotions, pauses, and stress
func (t *TeluguSSMLBuilder) processTextWithEmotions(text string) string {
	// Process in order: emotions first, then pauses, then stress
	// This prevents overlapping tags

	// Step 1: Add emotional emphasis to keywords (before punctuation processing)
	text = t.addEmotionalEmphasis(text)

	// Step 2: Add stress to words with vothulu (before punctuation processing)
	text = t.addVothuluStress(text)

	// Step 3: Add breathing and emotional expressions
	text = t.addBreathingAndEmotions(text)

	// Step 4: Add sentence-level excitement for exclamation sentences
	text = t.addSentenceExcitement(text)

	// Step 5: Handle punctuation-based emotions
	text = t.addPunctuationEmotions(text)

	// Step 6: Add pauses for punctuation (last to avoid conflicts)
	text = t.addPunctuationPauses(text)

	return text
}

// addBreathingAndEmotions adds breathing, smiling, and emotional expressions
func (t *TeluguSSMLBuilder) addBreathingAndEmotions(text string) string {
	// Add breathing before emotional words
	breathingWords := map[string]string{
		"ఆనందం":    `<break time="200ms"/><prosody rate="slow" pitch="+5%" volume="soft">*inhales deeply*</prosody><break time="100ms"/>`,
		"సంతోషం":   `<break time="200ms"/><prosody rate="slow" pitch="+5%" volume="soft">*inhales happily*</prosody><break time="100ms"/>`,
		"ఉత్సాహం":  `<break time="200ms"/><prosody rate="fast" pitch="+10%" volume="medium">*breathes excitedly*</prosody><break time="100ms"/>`,
		"ఆశ్చర్యం": `<break time="200ms"/><prosody rate="slow" pitch="+15%" volume="medium">*gasps*</prosody><break time="100ms"/>`,
		"భయం":      `<break time="200ms"/><prosody rate="slow" pitch="-10%" volume="soft">*shudders*</prosody><break time="100ms"/>`,
		"కోపం":     `<break time="200ms"/><prosody rate="fast" pitch="+5%" volume="loud">*huffs*</prosody><break time="100ms"/>`,
		"ప్రేమ":    `<break time="200ms"/><prosody rate="slow" pitch="-5%" volume="soft">*sighs lovingly*</prosody><break time="100ms"/>`,
		"విచారం":   `<break time="200ms"/><prosody rate="slow" pitch="-10%" volume="soft">*sighs sadly*</prosody><break time="100ms"/>`,
	}

	for word, breathing := range breathingWords {
		pattern := `\b` + regexp.QuoteMeta(word) + `\b`
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, breathing+word)
	}

	// Add smiling expressions for positive words
	smilingWords := map[string]string{
		"నవ్వు":      `<prosody rate="medium" pitch="+10%">`,
		"ఆనంద":       `<prosody rate="medium" pitch="+10%">`,
		"సంతోష":      `<prosody rate="medium" pitch="+10%">`,
		"ఉత్సాహ":     `<prosody rate="medium" pitch="+12%" volume="loud">`,
		"అద్భుత":     `<prosody rate="slow" pitch="+15%" volume="loud">`,
		"అందమైన":     `<prosody rate="slow" pitch="+5%">`,
		"అందంగా":     `<prosody rate="slow" pitch="+5%">`,
		"చక్కని":     `<prosody rate="slow" pitch="+5%">`,
		"బాగున్నాయి": `<prosody rate="medium" pitch="+10%">`,
		"చాలా":       `<prosody rate="medium" pitch="+5%">`,
		"Wow":        `<prosody rate="slow" pitch="+15%" volume="loud">`,
		"ఎంత":        `<prosody rate="medium" pitch="+15%">`,
	}

	for word, expression := range smilingWords {
		pattern := `\b` + regexp.QuoteMeta(word) + `\b`
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, expression+word+`</prosody>`)
	}

	// Add emotional expressions for story elements
	storyExpressions := map[string]string{
		"ఒక ఉదయం":   `<prosody rate="medium" pitch="+10%">`,
		"ఒక రోజు":   `<prosody rate="medium" pitch="+5%">`,
		"ఒక రాత్రి": `<prosody rate="slow" pitch="-5%">`,
		"అప్పుడే":   `<prosody rate="medium" pitch="+12%" volume="loud">`,
		"అక్కడ":     `<prosody rate="medium" pitch="+5%">`,
		"ఇక్కడ":     `<prosody rate="medium" pitch="+5%">`,
		"అదే":       `<prosody rate="slow" pitch="+10%">`,
		"ఇదే":       `<prosody rate="slow" pitch="+10%">`,
	}

	for phrase, expression := range storyExpressions {
		pattern := `\b` + regexp.QuoteMeta(phrase) + `\b`
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, expression+phrase+`</prosody>`)
	}

	return text
}

// addSentenceExcitement adds excitement to entire sentences ending with !
func (t *TeluguSSMLBuilder) addSentenceExcitement(text string) string {
	// Find sentences ending with ! and wrap them with excitement
	// This regex finds text that ends with ! (including the !)
	exclamationPattern := regexp.MustCompile(`([^.!?]*!)`)

	text = exclamationPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Remove the ! temporarily to process the sentence
		sentence := strings.TrimSpace(strings.TrimSuffix(match, "!"))
		if sentence == "" {
			return match
		}

		// Check if sentence has comma - if so, stress the part before comma more
		if strings.Contains(sentence, ",") {
			// Split by comma and stress each part
			parts := strings.Split(sentence, ",")
			var processedParts []string

			for i, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}

				// Last part before ! gets extra excitement
				if i == len(parts)-1 {
					processedParts = append(processedParts, `<prosody rate="medium" pitch="+20%" volume="loud">`+part+`</prosody>`)
				} else {
					// Parts before comma get high excitement
					processedParts = append(processedParts, `<prosody rate="medium" pitch="+15%" volume="loud">`+part+`</prosody>`)
				}
			}

			// Join with excited comma
			excitedSentence := strings.Join(processedParts, `<prosody rate="slow" pitch="+10%" volume="medium">,</prosody> `)
			return excitedSentence + `<prosody rate="slow" pitch="+25%" volume="loud">!</prosody>`
		} else {
			// No comma - stress the whole sentence
			return `<prosody rate="slow" pitch="+20%" volume="loud">` + sentence + `</prosody><prosody rate="slow" pitch="+25%" volume="loud">!</prosody>`
		}
	})

	return text
}

// addPunctuationEmotions adds emotion based on punctuation
func (t *TeluguSSMLBuilder) addPunctuationEmotions(text string) string {
	// Handle double exclamation
	text = regexp.MustCompile(`!!+`).ReplaceAllStringFunc(text, func(match string) string {
		return `<prosody rate="fast" pitch="high">` + match + `</prosody>`
	})

	// Handle single exclamation
	text = regexp.MustCompile(`!`).ReplaceAllStringFunc(text, func(match string) string {
		return `<prosody rate="medium" pitch="high">` + match + `</prosody>`
	})

	// Handle question marks
	text = regexp.MustCompile(`\?+`).ReplaceAllStringFunc(text, func(match string) string {
		return `<prosody rate="medium" pitch="rising">` + match + `</prosody>`
	})

	// Ellipsis is handled in addPunctuationPauses with proper pause timing

	return text
}

// addPunctuationPauses adds appropriate pauses for punctuation
func (t *TeluguSSMLBuilder) addPunctuationPauses(text string) string {
	// Process all punctuation in one pass to avoid conflicts

	// First handle ellipsis - replace with just a pause
	text = regexp.MustCompile(`\.{3,}`).ReplaceAllString(text, `<break time="400ms"/>`)

	// Handle single periods
	text = regexp.MustCompile(`\.`).ReplaceAllString(text, `.<break time="700ms"/>`)

	// Handle other punctuation
	// No pause after comma - removed for natural flow
	text = regexp.MustCompile(`;`).ReplaceAllString(text, `;<break time="300ms"/>`)
	text = regexp.MustCompile(`:`).ReplaceAllString(text, `:<break time="200ms"/>`)

	return text
}

// addVothuluStress adds stress to words containing vothulu
func (t *TeluguSSMLBuilder) addVothuluStress(text string) string {
	// Use regex to find words that are not already inside SSML tags
	wordPattern := regexp.MustCompile(`\b(\S+)\b`)

	return wordPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Check if this word is already inside SSML tags
		if strings.Contains(match, "<") || strings.Contains(match, ">") {
			return match
		}

		// Check if word has vothulu and add maximum emphasis with multiple effects
		if t.HasVothulu(match) {
			return `<prosody rate="x-slow" pitch="+20%" volume="loud"><emphasis level="strong">` + match + `</emphasis></prosody>`
		}

		return match
	})
}

// HasVothulu checks if a word contains Telugu vowel signs
func (t *TeluguSSMLBuilder) HasVothulu(word string) bool {
	for _, char := range word {
		if t.vothulu[char] {
			return true
		}
	}
	return false
}

// addEmotionalEmphasis adds emphasis to emotional keywords
func (t *TeluguSSMLBuilder) addEmotionalEmphasis(text string) string {
	for keyword, level := range t.emotionalKeywords {
		// Create word boundary regex for the keyword
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		re := regexp.MustCompile(pattern)

		text = re.ReplaceAllStringFunc(text, func(match string) string {
			// Check if this word is already inside SSML tags
			if strings.Contains(match, "<") || strings.Contains(match, ">") {
				return match
			}
			// Add extra stress for emotional keywords
			if level == "strong" {
				return `<prosody rate="slow" pitch="+15%" volume="loud"><emphasis level="` + level + `">` + match + `</emphasis></prosody>`
			}
			return `<emphasis level="` + level + `">` + match + `</emphasis>`
		})
	}

	return text
}

// BuildTeluguSSMLFromStory is a convenience function that creates a builder and processes the story
func BuildTeluguSSMLFromStory(story string) string {
	builder := NewTeluguSSMLBuilder()
	return builder.BuildTeluguSSML(story)
}

// AddCustomEmotionalKeyword adds a custom emotional keyword to the builder
func (t *TeluguSSMLBuilder) AddCustomEmotionalKeyword(keyword, emphasisLevel string) {
	t.emotionalKeywords[keyword] = emphasisLevel
}

// AddCustomVothulu adds a custom vothulu character to the builder
func (t *TeluguSSMLBuilder) AddCustomVothulu(vothulu rune) {
	t.vothulu[vothulu] = true
}

// GetSupportedEmotions returns the list of supported emotional emphasis levels
func (t *TeluguSSMLBuilder) GetSupportedEmotions() []string {
	return []string{"strong", "moderate", "reduced"}
}

// GetSupportedPunctuationEmotions returns the list of supported punctuation-based emotions
func (t *TeluguSSMLBuilder) GetSupportedPunctuationEmotions() []string {
	return []string{"excited", "questioning", "suspense", "very_excited", "surprised"}
}
