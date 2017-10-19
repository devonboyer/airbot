package shows

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devonboyer/airbot/airtable"
)

var baseID, tableID string

func init() {
	baseID = "appwqWzX94IXnLEp5"
	tableID = "Shows"
}

func TodayHandler(ctx context.Context, client *airtable.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		day := time.Now().Weekday()
		// Get shows
		shows := &ShowList{}
		err := client.Base(baseID).
			Table(tableID).
			List().
			FilterByFormula(fmt.Sprintf("{Day of Week} = '%s'", day)).
			Do(ctx, shows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Send message.
	}
}
