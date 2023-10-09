package dbot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func New(logger *log.Logger, version string, config Config) *Bot {
	return &Bot{
		Logger:  logger,
		Config:  config,
		Version: version,
	}
}

type Bot struct {
	Logger  *log.Logger
	Session *discordgo.Session
	Version string
	Config  Config
}

func (b *Bot) SetupBot() {
	var err error
	b.Session, err = discordgo.New("Bot " + b.Config.Token)
	if err != nil {
		b.Logger.Fatal("Failed to setup b: ", err)
	}
}

func (b *Bot) OnReady(_ *discordgo.Ready) {
	b.Logger.Printf("Botkisser ready")
}
