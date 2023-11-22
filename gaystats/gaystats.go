package gaystats

// gay

import (
	dbot "github.com/SqueakyBeaver/botkisser"
)

// Not a great solution, but I can't think of a better way to
// let everything see the bot client thingy mabobber
var botClient *dbot.Bot

func SetupListeners(bot *dbot.Bot) {
	botClient = bot
	bot.Session.AddHandler(messageReactionAdd)
	bot.Session.AddHandler(messageReactionRemove)
}
