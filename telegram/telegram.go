package telegram

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/wormi4ok/askachay/config"
	"github.com/wormi4ok/askachay/internal"
)

type Bot struct {
	tgAPI *tgbotapi.BotAPI
	users []config.User
}

func NewBot(apiToken string, users []config.User, debug bool) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, fmt.Errorf("faild to initialize Telegram API client: %w ", err)
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{tgAPI: bot, users: users}, nil
}

func (bot *Bot) ServeUpdates(app *internal.App, bundle *i18n.Bundle) error {
	var lang string

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.tgAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		dest := bot.getUploadPath(update.Message.From.UserName)
		if dest != "" {
			lang = "ru"
		} else {
			lang = update.Message.From.LanguageCode
		}

		m := localizedMessages(i18n.NewLocalizer(bundle, lang))

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			err := bot.handleCommand(update, m)
			if err == nil {
				continue
			}
			if errors.Is(err, UnsupportedCommandErr) {
				if err := bot.respond(msg, m.UnknownCommand); err != nil {
					return err
				}
				continue
			}
			return err
		}

		if dest == "" {
			if err := bot.respond(msg, m.BailOut); err != nil {
				return err
			}
			continue
		}

		msg.ReplyToMessageID = update.Message.MessageID
		if err := bot.respond(msg, m.Confirmation); err != nil {
			return err
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
		err := app.FetchMusic(update.Message.Text, dest)
		if errors.Is(err, internal.UnrecognizedInputErr) {
			if err := bot.respond(msg, m.UnrecognizedInput); err != nil {
				return err
			}
			continue
		}

		if err != nil {
			log.Printf("Error: %s", err)
			if err := bot.respond(msg, fmt.Sprintf(m.Error, err)); err != nil {
				return err
			}
			continue
		}

		msg.ReplyToMessageID = update.Message.MessageID
		err = bot.respond(msg, m.Success)
		if err == nil {
			continue
		}

		return err
	}

	return nil
}

func (bot *Bot) respond(msg tgbotapi.MessageConfig, text string) error {
	msg.Text = text
	_, err := bot.tgAPI.Send(msg)
	return err
}

var UnsupportedCommandErr = errors.New("unsupported command")

func (bot *Bot) handleCommand(update tgbotapi.Update, m messages) error {
	cmd := update.Message.Command()

	if cmd == "help" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, m.Help)
		_, err := bot.tgAPI.Send(msg)
		return err
	}

	if cmd != "start" {
		return fmt.Errorf("%w: %s", UnsupportedCommandErr, cmd)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if bot.getUploadPath(update.Message.From.UserName) != "" {
		msg.Text = m.PrivateGreeting
	} else {
		msg.Text = m.Greeting
	}

	_, err := bot.tgAPI.Send(msg)
	return err
}

func (bot *Bot) getUploadPath(userName string) string {
	for _, user := range bot.users {
		if userName == user.Username {
			return user.UploadPath
		}
	}
	return ""
}
