package airbot

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/botengine"
)

type Bot struct {
	*botengine.Engine
	*airtable.Client
	baseID, tableID string
}

func NewBot(secrets *Secrets, source botengine.Source, sink botengine.Sink) *Bot {
	bot := &Bot{
		botengine.New(source, sink),
		airtable.New(secrets.Airtable.APIKey),
		"appwqWzX94IXnLEp5",
		"Shows",
	}
	bot.setupHandlers()
	return bot
}

func (b *Bot) setupHandlers() {
	b.Handle("shows today", func(usr botengine.User) (string, error) {
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
