package internal

import (
	"fmt"
	"io"
)

type Downloader interface {
	Download(url string) (io.ReadCloser, error)
}

type Uploader interface {
	Upload(file io.Reader, path, name string) error
}

type App struct {
	s  Searcher
	dl Downloader
	up Uploader
}

func NewApp(s Searcher, dl Downloader, up Uploader) *App {
	return &App{s: s, dl: dl, up: up}
}

func (app *App) FetchMusic(input string, dest string) error {
	name, url, err := app.s.SearchMusic(input)
	if err != nil {
		return err
	}

	f, err := app.dl.Download(url)
	if err != nil {
		return err
	}

	err = app.up.Upload(f, dest, fmt.Sprintf("%s.mp3", name))
	if err != nil {
		return err
	}
	_ = f.Close()
	return nil
}
