package commands

import (
	"log"

	dbot "github.com/SqueakyBeaver/botkisser"
	"github.com/bwmarrin/discordgo"
)

/*
  - type ApplicationCommand struct {
    ID                string                 `json:"id,omitempty"`
    ApplicationID     string                 `json:"application_id,omitempty"`
    GuildID           string                 `json:"guild_id,omitempty"`
    Version           string                 `json:"version,omitempty"`
    Type              ApplicationCommandType `json:"type,omitempty"`
    Name              string                 `json:"name"`
    NameLocalizations *map[Locale]string     `json:"name_localizations,omitempty"`
    DefaultMemberPermissions *int64 `json:"default_member_permissions,string,omitempty"`
    DMPermission             *bool  `json:"dm_permission,omitempty"`
    NSFW                     *bool  `json:"nsfw,omitempty"`

    Description              string                      `json:"description,omitempty"`
    DescriptionLocalizations *map[Locale]string          `json:"description_localizations,omitempty"`
    Options                  []*ApplicationCommandOption `json:"options"`
    }
*/
var commands = []*discordgo.ApplicationCommand{
	&PingCommand,
	&TrumpPrideCommand,
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":        PingCommandResponse,
	"trump-pride": TrumpPrideResponse,
}

func SetupCommands(bot *dbot.Bot) {
	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func SyncCommands(bot *dbot.Bot, cfg *dbot.Config, log *log.Logger) {
	if cfg.DevMode {
		for _, v := range commands {
			_, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, cfg.DevGuildID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
		}
	} else {
		for _, v := range commands {
			_, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, "", v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
		}
	}
}
