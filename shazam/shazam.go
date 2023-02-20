package shazam

import (
	"fmt"
	"regexp"
	"strings"

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
	name = removeSpecialChars(name)
	name = removeDoubleSpaces(name)
	return name
}

func removeDoubleSpaces(name string) string {
	reg := regexp.MustCompile(`\s+`)
	name = reg.ReplaceAllString(name, " ")
	return strings.TrimSpace(name)
}

func removeShazamShareText(name string) string {
	if i := strings.Index(strings.ToLower(name), "shazam:"); i > 0 {
		return name[i+len("shazam:"):]
	}
	return name
}

func removeUrls(s string) string {
	reg := regexp.MustCompile(`(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(/\S*)?`)
	return reg.ReplaceAllString(s, "")
}

func removeSpecialChars(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z \-']+`)
	return reg.ReplaceAllString(s, "")
}
