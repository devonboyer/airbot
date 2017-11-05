package airbot

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/botengine"
	"github.com/sirupsen/logrus"
)

const (
	showsBaseID     = "appwqWzX94IXnLEp5"
	showsTableID    = "Shows"
	showsFormulaFmt = "AND({Day of Week} = '%s', {Status} = 'Airing')"
)

// ShowsBase provides bot handlers for retrieving data from the Shows airtable base.
type ShowsBase struct {
	*airtable.BaseScopedClient
}

func NewShowsBase(client *airtable.Client) *ShowsBase {
	return &ShowsBase{
		client.WithBaseScope(showsBaseID),
	}
}

func (b *ShowsBase) TodayHandler() func(w io.Writer, ev *botengine.Event) {
	return func(w io.Writer, ev *botengine.Event) {
		logrus.WithField("pattern", "shows today").Info("handler called")

		ctx := context.Background()
		day := time.Now().Weekday()
		shows := &ShowList{}
		err := b.
			Table(showsTableID).
			List().
			FilterByFormula(fmt.Sprintf(showsFormulaFmt, day)).
			Do(ctx, shows)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			if len(shows.Records) > 0 {
				fmt.Fprintf(w, "Shows on today:\n%s", shows)
			} else {
				fmt.Fprintf(w, "No shows on today")
			}
		}
	}
}

func (b *ShowsBase) TomorrowHandler() func(w io.Writer, ev *botengine.Event) {
	return func(w io.Writer, ev *botengine.Event) {
		logrus.WithField("pattern", "shows tomorrow").Info("handler called")

		ctx := context.Background()
		day := time.Now().Add(24 * time.Hour).Weekday()
		shows := &ShowList{}
		err := b.
			Table(showsTableID).
			List().
			FilterByFormula(fmt.Sprintf(showsFormulaFmt, day)).
			Do(ctx, shows)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			if len(shows.Records) > 0 {
				fmt.Fprintf(w, "Shows on tomorrow:\n%s", shows)
			} else {
				fmt.Fprintf(w, "No shows on tomorrow")
			}
		}
	}
}
