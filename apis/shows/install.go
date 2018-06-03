package shows

import (
	"github.com/devonboyer/airbot"
	"github.com/devonboyer/airbot/airtable"
	"github.com/devonboyer/airbot/botengine"
)

func Install(bot *botengine.Bot, client *airtable.Client) {
	shows := airbot.NewShowsBase(client)
	bot.HandleFunc("shows today", shows.TodayHandler())
	bot.HandleFunc("shows tomorrow", shows.TomorrowHandler())
}
