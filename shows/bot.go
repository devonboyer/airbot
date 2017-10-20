package shows

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devonboyer/airbot/airtable"
)

type Bot struct {
	client  *airtable.Client
	baseID  string
	tableID string
}

func New(apiKey string) *Bot {
	return &Bot{
		client:  airtable.New(apiKey),
		baseID:  "appwqWzX94IXnLEp5",
		tableID: "Shows",
	}
}

func (b *Bot) TodayHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()

		// Get shows
		shows := &ShowList{}
		err := b.client.Base(b.baseID).
			Table(b.tableID).
			List().
			FilterByFormula(fmt.Sprintf("{Day of Week} = '%s'", today())).
			Do(ctx, shows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Process shows.
		_ = b.processList(shows)
	}
}

func (b *Bot) processList(shows *ShowList) string {
	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, "Shows on tonight:")
	for _, s := range shows.Records {
		fmt.Fprintln(buf, s.Fields.Name)
	}
	return buf.String()
}

func today() time.Weekday {
	return time.Now().Weekday()
}
