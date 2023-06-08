package shazam

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/wormi4ok/askachay/internal"
	"github.com/wormi4ok/askachay/youtube"
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
	asciiName, err := convertToAscii(name)
	if err != nil {
		asciiName = name
	}
	name = convertSeparator(asciiName)
	name = removeSpecialChars(name)

	return name
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

func removeSpecialChars(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z \-']+`)
	return reg.ReplaceAllString(s, "")
}

func convertSeparator(s string) string {
	reg := regexp.MustCompile(`( от | by )`)
	return reg.ReplaceAllString(s, " - ")
}

func convertToAscii(s string) (string, error) {
	result, _, err := transform.String(transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn))), s)
	if err != nil {
		return "", err
	}
	return result, nil
}
