package util

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
)

// CountBillableCharacters estimates billable chars for Google TTS.
// - If input is SSML, pass isSSML=true; tags are ignored.
// - Counts Unicode runes in the final text content.
func CountBillableCharacters(input string, isSSML bool) int {
	text := input
	if isSSML {
		text = extractTextFromSSML(input)
	}
	// Optional normalization you actually apply before calling TTS:
	text = strings.TrimSpace(text)
	// If you collapse whitespace in your pipeline, do it here too:
	text = strings.Join(strings.Fields(text), " ")
	return utf8.RuneCountInString(text)
}

func extractTextFromSSML(ssml string) string {
	node, err := html.Parse(strings.NewReader(ssml))
	if err != nil {
		// Fallback: crude tag removal; still better than counting tags
		return stripTags(ssml)
	}
	var buf bytes.Buffer
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(node)
	return buf.String()
}

func stripTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}