package util

import (
	"regexp"
	"strings"
)

// BuildSSMLFromStory converts plain AI-generated story text into SSML.
// Heuristics applied:
// 1) Split into paragraphs by blank lines → wrap each in <p> ... </p>
// 2) Split into sentences (. ! ?) → insert short <break time="700ms"/> between sentences
// 3) Add longer pause between paragraphs (<break time="2s"/>)
// 4) Emphasize impactful keywords (e.g., brave, danger, mysterious, hero, magic)
// 5) XML-escape special characters to keep SSML valid
func BuildSSMLFromStory(story string) string {
	story = strings.TrimSpace(story)
	if story == "" {
		return "<speak></speak>"
	}

	// Local helper: XML escape
	escapeXML := func(s string) string {
		s = strings.ReplaceAll(s, "&", "&amp;")
		s = strings.ReplaceAll(s, "<", "&lt;")
		s = strings.ReplaceAll(s, ">", "&gt;")
		s = strings.ReplaceAll(s, "\"", "&quot;")
		s = strings.ReplaceAll(s, "'", "&apos;")
		return s
	}

	// Normalize whitespace
	normalize := func(s string) string {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		s = strings.ReplaceAll(s, "\r", "\n")
		// Collapse spaces within lines
		spaceRe := regexp.MustCompile(`\s+`)
		s = spaceRe.ReplaceAllString(s, " ")
		// Restore paragraph breaks (two or more newlines)
		paraRe := regexp.MustCompile(`\n{2,}`)
		s = paraRe.ReplaceAllString(s, "\n\n")
		return strings.TrimSpace(s)
	}

	story = normalize(story)

	// Split paragraphs by blank lines
	paragraphs := strings.Split(story, "\n\n")

	// Prepare keyword emphasis (case-insensitive word boundaries)
	keywords := []string{
		"brave", "hero", "danger", "mysterious", "magic", "magical", "whisper", "silence",
		"dark", "shadow", "curse", "destiny", "legend", "secret", "forbidden",
		"victory", "fear", "hope",
	}
	// Build a regex that matches any keyword as a whole word, case-insensitive
	// Example: (?i)\b(brave|hero|danger)\b
	kwPattern := "(?i)\\b(" + strings.Join(keywords, "|") + ")\\b"
	kwRe := regexp.MustCompile(kwPattern)

	// Local helper: emphasize keywords in a sentence
	emphasizeKeywords := func(s string) string {
		return kwRe.ReplaceAllStringFunc(s, func(m string) string {
			return "<emphasis level=\"moderate\">" + m + "</emphasis>"
		})
	}

	// Sentence splitter: split on ., !, ? while keeping punctuation
	// We first separate sentence end marks with a delimiter, then split
	endMarkRe := regexp.MustCompile(`([.!?])`)

	var b strings.Builder
	b.WriteString("<speak>")

	for pi, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// Escape paragraph text first so tags we add are the only tags present
		p = escapeXML(p)

		// Split sentences
		// Insert a pipe after sentence terminators, then split on the pipe
		marked := endMarkRe.ReplaceAllString(p, "${1}|")
		rawSentences := strings.Split(marked, "|")

		// Write paragraph
		b.WriteString("<p>")

		firstSentence := true
		for _, s := range rawSentences {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}

			// Re-normalize spaces in the sentence
			s = strings.TrimSpace(s)

			// Emphasize keywords (after escaping)
			s = emphasizeKeywords(s)

			if !firstSentence {
				// Short pause between sentences
				b.WriteString(`<break time="700ms"/>`)
			}
			firstSentence = false

			b.WriteString(s)
		}

		b.WriteString("</p>")

		// Longer pause between paragraphs
		if pi < len(paragraphs)-1 {
			b.WriteString(`<break time="2s"/>`)
		}
	}

	b.WriteString("</speak>")
	return b.String()
}
