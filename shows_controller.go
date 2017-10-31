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

var showsScope = Scope{
	BaseID:  "appwqWzX94IXnLEp5",
	TableID: "Shows",
}

const showsFormulaFmt = "AND({Day of Week} = '%s', {Status} = 'Airing')"

type showsController struct {
	*ScopeController
}

func newShowsController(client *airtable.Client) *showsController {
	return &showsController{ScopeController: NewScopeController(client, showsScope)}
}

func (c *showsController) todayHandler() func(w io.Writer, ev *botengine.Event) {
	return func(w io.Writer, ev *botengine.Event) {
		logrus.WithField("pattern", "shows today").Info("handler called")

		ctx := context.Background()
		day := time.Now().Weekday()
		shows := &ShowList{}
		err := c.List(ctx, fmt.Sprintf(showsFormulaFmt, day), shows)
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
		shows := &ShowList{}
		err := c.List(ctx, fmt.Sprintf(showsFormulaFmt, day), shows)
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
