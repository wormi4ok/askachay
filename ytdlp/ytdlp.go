package ytdlp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/wormi4ok/askachay/internal"
)

// ytdlp implements both Searcher and Downloader interfaces.
// It uses yt-dlp cli tool to download files to a temporary directory
// reading cookies with authentication information from file.
type ytDlp struct {
	cookiesFile string
	tmpdir      string
}

func (y *ytDlp) SearchMusic(input string) (name, url string, err error) {
	if !isUrl(input) {
		return "", "", fmt.Errorf("only valid URLs are supported: %w", internal.UnrecognizedInputErr)
	}

	return getTitle(input), input, err
}

type VideoInfo struct {
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
}

func getTitle(VideoURL string) string {
	ytUrl, _ := url.Parse("https://www.youtube.com/oembed")
	params := ytUrl.Query()
	params.Add("format", "json")
	params.Add("url", VideoURL)
	ytUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", ytUrl.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Response is %d", resp.StatusCode))
	}

	var info VideoInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		panic(err)
	}

	return normalizedTitle(info)
}

func normalizedTitle(videoInfo VideoInfo) string {
	cleanBracketsRe := regexp.MustCompile(`\([^()]*\)`)
	reSquareBrackets := regexp.MustCompile(`\[.*?]`)
	authorName := strings.TrimSpace(videoInfo.AuthorName)
	title := strings.TrimSpace(videoInfo.Title)

	// remove "lyrics" from the title
	if i := strings.Index(strings.ToLower(title), "lyrics"); i > 0 {
		title = strings.Join([]string{title[:i], title[i+len("lyrics"):]}, "")
	}

	// Put author Name before the title, if reversed
	if strings.Contains(title, authorName) && !strings.HasPrefix(title, authorName) {
		title = strings.ReplaceAll(title, authorName, "")
		title = strings.Trim(title, "-– ")
		title = fmt.Sprintf("%s - %s", authorName, title)
	}

	// Capitalize all words
	title = cases.Title(language.English).String(title)
	// Remove brackets
	title = cleanBracketsRe.ReplaceAllString(title, "")
	title = reSquareBrackets.ReplaceAllString(title, "")

	return strings.TrimSpace(title)
}

func isUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func NewYtDlp(cookiesPath string) *ytDlp {
	return &ytDlp{
		cookiesFile: cookiesPath,
		tmpdir:      os.TempDir(),
	}
}

func (y *ytDlp) SetTmpDir(path string) error {
	fInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fInfo.IsDir() {
		return errors.New("directory does not exist")
	}

	y.tmpdir = path
	return err
}

func (y *ytDlp) Download(url string) (io.ReadSeekCloser, error) {
	tmpFile := filepath.Join(y.tmpdir, "output.mp3")
	_, err := os.Stat(tmpFile)
	if err == nil {
		err = os.Remove(tmpFile)
		if err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"--no-cache-dir",
		"--cookies", y.cookiesFile,
		"--abort-on-error",
		"--newline",
		"--restrict-filenames",
		"--format", "140",
		"--extract-audio",
		"--audio-format", "mp3",
		"-o", tmpFile,
		url,
	)

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	log.Printf("execuriting command: %s", cmd.String())
	err = cmd.Run()
	errMessage := ""
	if err != nil && strings.Contains(err.Error(), "exit status 1") {
		stderrLineScanner := bufio.NewScanner(&stderrBuf)

		for stderrLineScanner.Scan() {
			r := regexp.MustCompile("(?i)error: (.*)")
			line := stderrLineScanner.Text()
			if i := r.FindStringSubmatchIndex(line); i != nil {
				errMessage += line[i[0]:]
			}
		}
		if errMessage != "" {
			return nil, fmt.Errorf("extracted error: %s", errMessage)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("unprocessed error: %w", err)
	}

	f, err := os.Open(tmpFile)
	if err != nil {
		return nil, err
	}

	return f, err
}
