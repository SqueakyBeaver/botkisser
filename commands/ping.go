package commands

import (
	"github.com/bwmarrin/discordgo"
)

var PingCommand = discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "pong .o.",
}

func PingCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong owo",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
