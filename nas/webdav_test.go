package nas

import (
	"os"
	"slices"
	"testing"
)

func TestGetShareList(t *testing.T) {
	host := os.Getenv("WEBDAV_HOST")
	user := os.Getenv("WEBDAV_USER")
	pass := os.Getenv("WEBDAV_PASS")
	uploadPath := os.Getenv("WEBDAV_UPLOAD_PATH")
	if anyIsEmpty(host, user, pass, uploadPath) {
		t.Skip("WEBDAV_* env vars are not set. Skipping...")
	}

	client := NewWebDavClient(host, user, pass)
	if err := client.GetShareList(uploadPath); err != nil {
		t.Error(err)
	}
}

func anyIsEmpty(ss ...string) bool {
	return slices.Contains(ss, "")
}
