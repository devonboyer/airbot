package messenger

import "github.com/abhinavdahiya/go-messenger-bot"

type Bot struct {
	*mbotapi.BotAPI
}

func NewBot() *Bot {
	//bot := mbotapi.NewBotAPI("ACCESS_TOKEN", "VERIFY_TOKEN")

}
