package airbot

import (
	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/botengine"
)

type Bot struct {
	*botengine.Engine
}

func NewBot(client *airtable.Client, source botengine.Source, sink botengine.Sink) *Bot {
	bot := &Bot{
		botengine.New(source, sink, botengine.DefaultSettings),
	}
	bot.setupHandlers(client)
	return bot
}

func (b *Bot) setupHandlers(client *airtable.Client) {
	shows := NewShowsBase(client)
	b.HandleFunc("shows today", shows.TodayHandler())
	b.HandleFunc("shows tomorrow", shows.TomorrowHandler())
}
