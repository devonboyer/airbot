package airbot

import (
	"testing"

	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/botengine"
	"github.com/stretchr/testify/require"
)

func Test_ShowsBase(t *testing.T) {
	client := airtable.New(secrets.Airtable.APIKey)
	base := NewShowsBase(client)

	tests := []struct {
		name    string
		handler botengine.Handler
	}{
		{
			"today",
			botengine.HandlerFunc(base.TodayHandler()),
		},
		{
			"tomorrow",
			botengine.HandlerFunc(base.TomorrowHandler()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := botengine.NewRecorder()
			tt.handler.Handle(rr, nil)
			require.NotEmpty(t, botengine.StatusOk, rr.Status)
			require.NotEmpty(t, rr.Body.String())
		})
	}
}
