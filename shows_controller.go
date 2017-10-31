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
	*airtable.TableScopedClient
}

func newShowsController(client *airtable.Client) *showsController {
	return &showsController{client.WithTableScope("appwqWzX94IXnLEp5", "Shows")}
}

func (c *showsController) todayHandler() func(w io.Writer, ev *botengine.Event) {
	return func(w io.Writer, ev *botengine.Event) {
		logrus.WithField("pattern", "shows today").Info("handler called")

		ctx := context.Background()
		day := time.Now().Weekday()
		shows := &ShowList{}
		err := c.
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

func (c *showsController) tomorrowHandler() func(w io.Writer, ev *botengine.Event) {
	return func(w io.Writer, ev *botengine.Event) {
		logrus.WithField("pattern", "shows tomorrow").Info("handler called")

		ctx := context.Background()
		day := time.Now().Add(24 * time.Hour).Weekday()
		shows := &ShowList{}
		err := c.
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
