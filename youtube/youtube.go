package youtube

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	td "github.com/ChannelMeter/iso8601duration"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/wormi4ok/askachay/internal"
)

const minVideoDuration = 2*time.Minute + 30*time.Second
const maxVideoDuration = 10 * time.Minute

type Client struct {
	base *youtube.Service
}

func NewClient(apiKey string) (*Client, error) {
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	return &Client{base: youtubeService}, err
}

func (c *Client) SearchMusic(input string) (name, url string, err error) {
	if !isSongName(input) {
		return "", "", fmt.Errorf("not a valid song name: %w", internal.UnrecognizedInputErr)
	}

	searchCall := c.base.Search.List([]string{"id", "snippet"}).
		Q(input + " lyrics").
		Type("video").
		MaxResults(5)
	response, err := searchCall.Do()
	if err != nil {
		return "", "", fmt.Errorf("failed to get search results: %w", err)
	}
	if len(response.Items) == 0 {
		return name, "", fmt.Errorf("nothing found by search request: %s", name)
	}

	for _, item := range response.Items {
		res, err := c.base.Videos.List([]string{"snippet", "contentDetails"}).
			Id(item.Id.VideoId).Do()
		if err != nil {
			fmt.Printf("failed to fetch video details: %c", err)
			continue
		}
		duration := res.Items[0].ContentDetails.Duration
		d, err := td.FromString(duration)
		if err != nil {
			fmt.Printf("failed to parse video duration: %c", err)
			continue
		}

		if d.ToDuration() < minVideoDuration || d.ToDuration() > maxVideoDuration {
			continue
		}
		return input, fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId), nil
	}

	return input, "", errors.New("no suitable videos found")
}

func isSongName(input string) bool {
	if !strings.Contains(input, " - ") {
		return false
	}
	if strings.Contains(input, "\n") {
		return false
	}
	return true
}
