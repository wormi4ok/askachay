package nas

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
)

type WebDavClient struct {
	*gowebdav.Client
}

func NewWebDavClient(host, user, pass string) *WebDavClient {
	return &WebDavClient{gowebdav.NewClient(host, user, pass)}
}

func (w *WebDavClient) Upload(f io.Reader, path, name string) error {
	target := filepath.Join(path, name)
	if err := w.WriteStream(target, f, 0660); err != nil {
		return fmt.Errorf("failed to upload a file: %w", err)
	}

	return nil
}

func (w *WebDavClient) GetShareList(path string) error {
	err := w.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to the target server via WebDAV: %w", err)
	}
	fileList, err := w.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to fetch a list of files: %w", err)
	}
	for _, fileInfo := range fileList {
		println(fileInfo.Name())
	}
	return nil
}
