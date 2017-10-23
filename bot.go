package airbot

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/bot"
)

type Bot struct {
	*bot.Bot
	*airtable.Client
	baseID, tableID string
}

func NewBot(secrets *Secrets, listener bot.Listener, sender bot.Sender) *Bot {
	bot := &Bot{
		bot.New(listener, sender), //
		airtable.New(secrets.Airtable.APIKey),
		"appwqWzX94IXnLEp5",
		"Shows",
	}
	bot.setupHandlers()
	return bot
}

func (b *Bot) setupHandlers() {
	b.Handle("shows today", func(s string) (string, error) {
		ctx := context.Background()
		shows := &ShowList{}
		formulaFmt := "AND({Day of Week} = '%s', {Status} = 'Airing')"
		err := b.Base(b.baseID).
			Table(b.tableID).
			List().
			FilterByFormula(fmt.Sprintf(formulaFmt, time.Now().Weekday())).
			Do(ctx, shows)
		if err != nil {
			return "", err
		}
		return processShows(shows), nil
	})
}

func processShows(shows *ShowList) string {
	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, "Shows on tonight:")
	for _, s := range shows.Records {
		fmt.Fprintln(buf, s.Fields.Name)
	}
	return buf.String()
}
