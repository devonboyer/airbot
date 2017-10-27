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

const showsFormulaFmt = "AND({Day of Week} = '%s', {Status} = 'Airing')"

type showsController struct {
	client          *airtable.Client
	baseID, tableID string
}

func newShowsController(client *airtable.Client) *showsController {
	return &showsController{
		client:  client,
		baseID:  "appwqWzX94IXnLEp5",
		tableID: "Shows",
	}
}

func (c *showsController) todayHandler() func(w io.Writer, ev *botengine.Event) {
	return func(w io.Writer, ev *botengine.Event) {
		logrus.WithField("pattern", "shows today").Info("handler called")

		ctx := context.Background()
		day := time.Now().Weekday()
		shows, err := c.listShows(ctx, fmt.Sprintf(showsFormulaFmt, day))
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

func (c *showsController) tomorrowHandler() func(w io.Writer, ev *botengine.Event) {
	return func(w io.Writer, ev *botengine.Event) {
		logrus.WithField("pattern", "shows tomorrow").Info("handler called")

		ctx := context.Background()
		day := time.Now().Add(24 * time.Hour).Weekday()
		shows, err := c.listShows(ctx, fmt.Sprintf(showsFormulaFmt, day))
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

func (c *showsController) listShows(ctx context.Context, filterByForumla string) (*ShowList, error) {
	shows := &ShowList{}
	err := c.client.
		Base(c.baseID).
		Table(c.tableID).
		List().
		FilterByFormula(filterByForumla).
		Do(ctx, shows)
	if err != nil {
		return nil, err
	}
	return shows, nil
}
