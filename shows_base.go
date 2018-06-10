package airbot

import (
	"context"
	"fmt"
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
	*airtable.TableClient
}

func NewShowsBase(client *airtable.Client) *ShowsBase {
	return &ShowsBase{client.Base(showsBaseID).Table(showsTableID)}
}

func (b *ShowsBase) TodayHandler() func(botengine.ResponseWriter, *botengine.Message) {
	return func(w botengine.ResponseWriter, _ *botengine.Message) {
		logrus.WithField("pattern", "shows today").Info("handler called")

		ctx := context.Background()
		day := time.Now().Weekday().String()
		shows, err := b.GetShows(ctx, day)
		if err != nil {
			botengine.Error(w, err)
		} else {
			handleShows(w, shows, "today")
		}
	}
}

func (b *ShowsBase) TomorrowHandler() func(botengine.ResponseWriter, *botengine.Message) {
	return func(w botengine.ResponseWriter, _ *botengine.Message) {
		logrus.WithField("pattern", "shows tomorrow").Info("handler called")

		ctx := context.Background()
		day := time.Now().Add(24 * time.Hour).Weekday().String()
		shows, err := b.GetShows(ctx, day)
		if err != nil {
			botengine.Error(w, err)
		} else {
			handleShows(w, shows, "tomorrow")
		}
	}
}

func (b *ShowsBase) GetShows(ctx context.Context, day string) (*ShowList, error) {
	shows := &ShowList{}
	opts := &airtable.ListOptions{FilterByFormula: fmt.Sprintf(showsFormulaFmt, day)}
	if err := b.List(ctx, opts, shows); err != nil {
		return nil, err
	}
	return shows, nil
}

func handleShows(w botengine.ResponseWriter, shows *ShowList, day string) {
	if len(shows.Records) > 0 {
		fmt.Fprintf(w, "Shows on %s:\n%s", day, shows)
	} else {
		fmt.Fprintf(w, "No shows on %s", day)
	}
	w.SetStatus(botengine.StatusOk)
}
