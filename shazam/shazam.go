package shazam

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/wormi4ok/askachay/internal"
	"github.com/wormi4ok/askachay/youtube"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Shazam struct {
	yt *youtube.Client
}

func NewShazam(youtubeService *youtube.Client) *Shazam {
	return &Shazam{youtubeService}
}

func (s *Shazam) SearchMusic(input string) (name, url string, err error) {
	if !strings.Contains(strings.ToLower(input), "shazam") {
		return "", "", fmt.Errorf("only Shazam links are supported: %w", internal.UnrecognizedInputErr)
	}
	name = extractName(input)

	return s.yt.SearchMusic(name)
}

func extractName(name string) string {
	name = removeShazamShareText(name)
	name = removeUrls(name)
	name = removeDoubleSpaces(name)
	name = swapTitleArtist(name)
	name = convertSeparator(name)
	name = trimPunctuation(name)
	name = removeDiacritics(name)

	return name
}

// If original used Russian ' от ' treat left as title and right as artist: swap to "Artist - Title"
// Swap only when both sides contain Cyrillic characters (to avoid swapping latin inputs like 'Baby Said от Måneskin').
func swapTitleArtist(name string) string {
	if strings.Contains(strings.ToLower(name), " от ") {
		parts := strings.SplitN(name, " от ", 2)
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			name = right + " - " + left
		}
	}
	return name
}

func trimPunctuation(s string) string {
	parts := strings.Split(s, " - ")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, " .,!?:;\"'«»()[]")
		parts[i] = p
	}
	return strings.Join(parts, " - ")
}

func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)
	return s
}

func removeDoubleSpaces(s string) string {
	reg := regexp.MustCompile(`\s+`)
	s = reg.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func removeShazamShareText(s string) string {
	if i := strings.Index(strings.ToLower(s), "shazam:"); i > 0 {
		return s[i+len("shazam:"):]
	}
	return s
}

func removeUrls(s string) string {
	reg := regexp.MustCompile(`(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(/\S*)?`)
	return reg.ReplaceAllString(s, "")
}

func convertSeparator(s string) string {
	reg := regexp.MustCompile(`( от | by )`)
	return reg.ReplaceAllString(s, " - ")
}
