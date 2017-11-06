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

func (b *ShowsBase) TodayHandler() func(io.Writer, *botengine.Message) {
	return func(w io.Writer, _ *botengine.Message) {
		logrus.WithField("pattern", "shows today").Info("handler called")

		ctx := context.Background()
		day := time.Now().Weekday().String()
		shows, err := b.GetShows(ctx, day)
		if err != nil {
			handleError(err, w)
		} else {
			if len(shows.Records) > 0 {
				fmt.Fprintf(w, "Shows on today:\n%s", shows)
			} else {
				fmt.Fprintf(w, "No shows on today")
			}
		}
	}
}

func (b *ShowsBase) TomorrowHandler() func(io.Writer, *botengine.Message) {
	return func(w io.Writer, _ *botengine.Message) {
		logrus.WithField("pattern", "shows tomorrow").Info("handler called")

		ctx := context.Background()
		day := time.Now().Add(24 * time.Hour).Weekday().String()
		shows, err := b.GetShows(ctx, day)
		if err != nil {
			handleError(err, w)
		} else {
			if len(shows.Records) > 0 {
				fmt.Fprintf(w, "Shows on tomorrow:\n%s", shows)
			} else {
				fmt.Fprintf(w, "No shows on tomorrow")
			}
		}
	}
}

func (b *ShowsBase) DayOfWeekHandler() func(io.Writer, *botengine.Message) {
	return func(w io.Writer, msg *botengine.Message) {
		if len(msg.Args()) < 1 {
			fmt.Fprint(w, "I didn't understand that. Maybe try again?")
			return
		}

		ctx := context.Background()
		day := msg.Args()[1]
		shows, err := b.GetShows(ctx, day)
		if err != nil {
			handleError(err, w)
		} else {
			if len(shows.Records) > 0 {
				fmt.Fprintf(w, "Shows on %s:\n%s", day, shows)
			} else {
				fmt.Fprintf(w, "No shows on %s", day)
			}
		}
	}
}

func (b *ShowsBase) GetShows(ctx context.Context, day string) (*ShowList, error) {
	shows := &ShowList{}
	err := b.
		Table(showsTableID).
		List().
		FilterByFormula(fmt.Sprintf(showsFormulaFmt, day)).
		Do(ctx, shows)
	if err != nil {
		return nil, err
	}
	return shows, nil
}

func handleError(err error, w io.Writer) {
	fmt.Fprint(w, err)
}
