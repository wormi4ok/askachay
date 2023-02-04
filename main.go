package main

import (
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kong"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"github.com/wormi4ok/askachay/config"
	"github.com/wormi4ok/askachay/internal"
	"github.com/wormi4ok/askachay/nas"
	"github.com/wormi4ok/askachay/shazam"
	"github.com/wormi4ok/askachay/telegram"
	"github.com/wormi4ok/askachay/youtube"
	"github.com/wormi4ok/askachay/ytdlp"
)

var CLI struct {
	URL        string `xor:"mode" help:"Source URL"`
	ConfigPath string `env:"CONFIG_PATH" default:"users.json"`
	Webdav     struct {
		Host       string `required:"" env:"HOST"`
		User       string `required:"" env:"USER"`
		Pass       string `required:"" env:"PASS"`
		UploadPath string `env:"UPLOAD_PATH"`
	} `embed:"" prefix:"webdav-" envprefix:"WEBDAV_"`
	Bot              bool   `env:"BOT_MODE" xor:"mode" required:""`
	TelegramApiToken string `env:"TELEGRAM_API_TOKEN"`
	Debug            bool   `env:"DEBUG"`

	YoutubeApiKey string   `env:"YOUTUBE_API_KEY"`
	LoadLanguage  []string `env:"LOAD_LANGUAGE"`
}

func main() {
	_ = kong.Parse(&CLI)

	cfg, err := config.FromFile(CLI.ConfigPath)
	if err != nil {
		log.Error(err)
	}

	yt, err := youtube.NewClient(CLI.YoutubeApiKey)
	if err != nil {
		log.Error(err)
	}
	upl := nas.NewWebDavClient(CLI.Webdav.Host, CLI.Webdav.User, CLI.Webdav.Pass)
	dl := ytdlp.NewYtDlp()
	sh := shazam.NewShazam(yt)

	bundle := configureLanguages(CLI.LoadLanguage)

	app := internal.NewApp(internal.NewMultiSearcher(dl, sh, yt), dl, upl)

	if CLI.Bot {
		bot, err := telegram.NewBot(CLI.TelegramApiToken, cfg.Users, CLI.Debug)
		if err != nil {
			log.Fatal(err)
		}
		err = bot.ServeUpdates(app, bundle)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if CLI.Webdav.UploadPath == "" {
		log.Fatal("missing flag --webdav-upload-path")
	}
	err = app.FetchMusic(CLI.URL, CLI.Webdav.UploadPath)
	if err != nil {
		log.Fatal(err)
	}
}

func configureLanguages(paths []string) *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, path := range paths {
		if _, err := bundle.LoadMessageFile(path); err != nil {
			log.Fatal(err)
		}
	}

	return bundle
}
