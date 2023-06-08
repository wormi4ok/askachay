package shazam

import (
	"os"
	"testing"

	"github.com/wormi4ok/askachay/youtube"
)

func Test_extractName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"shazam share",
			"Мое открытие на Shazam: Alice Cooper - Poison. https://www.shazam.com/track/55255280/poison?referrer=share",
			"Alice Cooper - Poison",
		},
		{
			"shazam share 2",
			"Мое открытие на Shazam: Johnny Cash - God's Gonna Cut You Down. https://www.shazam.com/track/44335694/gods-gonna-cut-you-down?referrer=share",
			"Johnny Cash - God's Gonna Cut You Down",
		},
		{
			"shazam share 3",
			"Baby Said от Måneskin https://www.shazam.com/track/645856503/baby-said?referrer=share",
			"Baby Said - Maneskin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractName(tt.input); got != tt.want {
				got, _ = convertToAscii(got)
				t.Errorf("extractName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShazam_SearchMusic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantUrl  string
		wantErr  bool
	}{
		{
			name:     "Happy path",
			input:    "Мое открытие на Shazam: Alice Cooper - Poison. https://www.shazam.com/track/55255280/poison?referrer=share",
			wantName: "Alice Cooper Poison",
			wantUrl:  "https://www.youtube.com/watch?v=RViRl2vjEmc",
			wantErr:  false,
		},
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		t.Skip("YOUTUBE_API_KEY is not set. Skipping...")
	}
	yt, err := youtube.NewClient(apiKey)
	if err != nil {
		t.Fatal(err)
	}
	s := NewShazam(yt)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotUrl, err := s.SearchMusic(tt.input)
			if gotName != tt.wantName {
				t.Errorf("SearchMusic() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotUrl != tt.wantUrl {
				t.Errorf("SearchMusic() gotUrl = %v, want %v", gotUrl, tt.wantUrl)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchMusic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
