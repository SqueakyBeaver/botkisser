package gaystats

import (
	"log"
	"slices"

	"github.com/SqueakyBeaver/botkisser/db"
	"github.com/bwmarrin/discordgo"
)

func messageReactionAdd(_ *discordgo.Session, event *discordgo.MessageReactionAdd) {
	if event.GuildID == "" {
		return
	}
	if event.Member.User.Bot {
		return
	}

	guildData, err := botClient.Database.GetGuildSettings(event.GuildID)
	if err != nil {
		botClient.Logger.Panicf("Error getting guild settings: %v", err)
	}

	if slices.Index(guildData.GayConf.BlocklistedChannels, event.ChannelID) != -1 {
		log.Printf("Channel %v blocklisted in guild %v for gay reactions", event.ChannelID, event.GuildID)
		return
	}

	if slices.Index(guildData.GayConf.BlocklistedUsers, event.UserID) != -1 {
		log.Printf("User %v blocklisted in guild %v for gay reactions", event.UserID, event.GuildID)
		return
	}

	emojiID := event.Emoji.ID

	if event.Emoji.ID == "" {
		if event.Emoji.Name != "" {
			emojiID = event.Emoji.Name
		}
	}
	memberData := botClient.Database.GetMemberStats(event.GuildID, event.UserID)

	switch emojiID {
	case guildData.GayConf.GayEmote:
		memberData.Gays++
	case guildData.GayConf.SuperGayEmote:
		memberData.SuperGays++
	}

	guildData.Members[event.UserID] = memberData

	botClient.Database.UpdateGuildSettings(event.GuildID, &db.Guild{Members: guildData.Members})
}
