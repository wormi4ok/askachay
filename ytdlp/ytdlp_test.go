package ytdlp

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os/exec"
	"testing"
)

func TestYtDlp_Download(t *testing.T) {
	_, err := exec.LookPath("yt-dlp")
	if err != nil {
		t.Skip("Missing yt-dlp binary. Skipping...")
	}
	type args struct {
		url string
	}
	tests := []struct {
		name         string
		args         args
		wantChecksum string
		wantErr      bool
	}{
		{
			name: "Happy path",
			args: args{
				url: "https://www.youtube.com/watch?v=NEDG31pF2sw",
			},
			wantChecksum: "d37f2e7134917144c66558ba00e55bff17e402c8930fdd62f724b2601340f722",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &ytDlp{
				tmpdir:      t.TempDir(),
				cookiesFile: "~/cookies.txt",
			}
			got, err := y.Download(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}

			h := sha256.New()
			if _, err = io.Copy(h, got); err != nil {
				t.Fatalf("Failed to calculate checksum: %s", err)
			}

			gotChecksum := fmt.Sprintf("%x", h.Sum(nil))
			if gotChecksum != tt.wantChecksum {
				t.Errorf("SHA256 checksum didn't match.\nWant = %s\nGot  = %s", tt.wantChecksum, gotChecksum)
			}
		})
	}
}

func Test_getTitle(t *testing.T) {
	type args struct {
		VideoURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy path",
			args: args{
				VideoURL: "https://www.youtube.com/watch?v=NEDG31pF2sw",
			},
			want: "Linkin Park - Lost In The Echo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTitle(tt.args.VideoURL); got != tt.want {
				t.Errorf("getTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizedTitle(t *testing.T) {
	tests := []struct {
		name  string
		input VideoInfo
		want  string
	}{
		{
			name: "The Title is already Perfect",
			input: VideoInfo{
				Title:      "Linkin Park - Numb",
				AuthorName: "Linkin Park",
			},
			want: "Linkin Park - Numb",
		},
		{
			name: "Put Author in front",
			input: VideoInfo{
				Title:      "LOST IN THE ECHO (Trailer) - Linkin Park",
				AuthorName: "Linkin Park",
			},
			want: "Linkin Park - Lost In The Echo",
		},
		{
			name: "Remove 'Lyrics'",
			input: VideoInfo{
				Title:      "Fall Out Boy I Don't Care Lyrics",
				AuthorName: "random uploader",
			},
			want: "Fall Out Boy I Don't Care",
		},
		{
			name: "Remove text in square brackets",
			input: VideoInfo{
				Title:      "Green Day - Holiday [Official Music Video]",
				AuthorName: "Green Day",
			},
			want: "Green Day - Holiday",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizedTitle(tt.input); got != tt.want {
				t.Errorf("normalizedTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
