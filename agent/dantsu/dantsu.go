package dantsu

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	fontOpenRe  = regexp.MustCompile(`(?i)<font[^>]*>`)
	fontCloseRe = regexp.MustCompile(`(?i)</font>`)
	boldOpenRe  = regexp.MustCompile(`(?i)<b>`)
	boldCloseRe = regexp.MustCompile(`(?i)</b>`)
	imgRe       = regexp.MustCompile(`(?i)<img>[^<]*</img>`)
)

// ToEscPos converts DantSu-style tagged text to ESC/POS bytes.
func ToEscPos(payload string) []byte {
	const esc = 0x1b
	const gs = 0x1d
	chunks := [][]byte{{esc, 0x40}}

	lines := strings.Split(payload, "\n")
	for _, rawLine := range lines {
		align := byte(0)
		text := rawLine
		switch {
		case strings.HasPrefix(text, "[C]"):
			align = 1
			text = text[3:]
		case strings.HasPrefix(text, "[L]"):
			text = text[3:]
		}

		boldOn := strings.Contains(strings.ToLower(text), "<b>")
		text = fontOpenRe.ReplaceAllString(text, "")
		text = fontCloseRe.ReplaceAllString(text, "")
		text = boldOpenRe.ReplaceAllString(text, "")
		text = boldCloseRe.ReplaceAllString(text, "")
		text = imgRe.ReplaceAllString(text, "")

		chunks = append(chunks, []byte{esc, 0x61, align})
		if boldOn {
			chunks = append(chunks, []byte{esc, 0x45, 1})
		}
		chunks = append(chunks, []byte(text+"\n"))
		if boldOn {
			chunks = append(chunks, []byte{esc, 0x45, 0})
		}
	}

	chunks = append(chunks, []byte{gs, 0x56, 0x42, 0x00})
	return bytes.Join(chunks, nil)
}
