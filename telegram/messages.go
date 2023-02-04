package telegram

import "github.com/nicksnyder/go-i18n/v2/i18n"

const helpText = "You can send me a YouTube link or share your Shazam finding and the song will appear on the NAS in your Music folder"
const privateGreetingText = `Hey there!
Simply send me a YouTube link or share your Shazam finding and the song will appear on the NAS in your Music folder`
const successText = `Done! File is uploaded to your music folder!`
const greetingText = `Hello! How can I help?`
const confirmationText = "Got it! Searching!"
const unrecognizedInputText = "I didn't understand that. Are you sure this is a song name?"
const unknownCommandText = "Sorry, I didn't understand that"
const bailOutText = "Sorry, this is a private bot. Please contact the author to enable it fo you."
const errorText = `I couldn't download the song you want...
Please send the author this text, he will fix it
` + "```" + `
%s
` + "```"

type messages struct {
	Help              string
	PrivateGreeting   string
	Success           string
	Greeting          string
	Confirmation      string
	UnrecognizedInput string
	UnknownCommand    string
	BailOut           string
	Error             string
}

func localizedMessages(l *i18n.Localizer) messages {
	return messages{
		Help:              l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "Help", Other: helpText}}),
		PrivateGreeting:   l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "PrivateGreeting", Other: privateGreetingText}}),
		Success:           l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "Success", Other: successText}}),
		Greeting:          l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "Greeting", Other: greetingText}}),
		Confirmation:      l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "Confirmation", Other: confirmationText}}),
		UnrecognizedInput: l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "UnrecognizedInput", Other: unrecognizedInputText}}),
		UnknownCommand:    l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "UnknownCommand", Other: unknownCommandText}}),
		BailOut:           l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "BailOut", Other: bailOutText}}),
		Error:             l.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "Error", Other: errorText}}),
	}
}
