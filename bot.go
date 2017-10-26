package airbot

import (
	"context"
	"fmt"
	"io"
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
	// Responds with a list of shows that are airing tonight
	b.Handle("shows today", func(w io.Writer, ev *botengine.Event) {
		ctx := context.Background()
		shows := &ShowList{}
		formulaFmt := "AND({Day of Week} = '%s', {Status} = 'Airing')"
		err := b.Base(b.baseID).
			Table(b.tableID).
			List().
			FilterByFormula(fmt.Sprintf(formulaFmt, time.Now().Weekday())).
			Do(ctx, shows)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			fmt.Fprintln(w, "Shows on tonight:")
			for _, s := range shows.Records {
				fmt.Fprintln(w, s.Fields.Name)
			}
		}
	})
}
