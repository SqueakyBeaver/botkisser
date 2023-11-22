package commands

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/h2non/bimg"

	"github.com/SqueakyBeaver/botkisser/utils"
)

var TrumpPrideCommand = discordgo.ApplicationCommand{
	Name:        "trump-pride",
	Description: "Trump with a pride flag",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "string-option",
			Description: "String option",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "public",
			Description: "whether or not to send this publicly",
			Required:    false,
		},
	},
}

func TrumpPrideResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user
	options := i.ApplicationCommandData().Options

	flagUrl := utils.GetPrideFlag(options[0].StringValue())

	flagDataBody, err := http.Get(flagUrl)
	if err != nil {
		log.Panicf("Error Getting pride flag data in TrumpPrideResponse: %v", err)
	}
	flagData, err := io.ReadAll(flagDataBody.Body)
	if err != nil {
		log.Panicf("Error Getting pride flag data in TrumpPrideResponse: %v", err)
	}
	trumpImageBuf, err := bimg.Read("assets/trump_mug.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	flagImage := bimg.NewImage(flagData)

	// Resize until at least 1024 x 1024; I want to keep the aspect ratio or else it looks weird
	flagSize, err := flagImage.Size()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newData, err := flagImage.Enlarge(1024, 1024)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newData, err = bimg.NewImage(newData).Crop(1024, 1024, bimg.GravityNorth)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	flagImage = bimg.NewImage(newData)

	flagSize, err = flagImage.Size()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	trumpImage := bimg.WatermarkImage{
		Left:    flagSize.Width/2 - 512,
		Top:     flagSize.Height - 1024,
		Buf:     trumpImageBuf,
		Opacity: 100,
	}

	flagImage.WatermarkImage(trumpImage)

	bimg.Write("trump-mug-tmp.png", flagImage.Image())
	sendFile, err := os.Open("trump-mug-tmp.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong owo",
			Flags:   discordgo.MessageFlagsEphemeral,
			Files: []*discordgo.File{
				{
					ContentType: "image",
					Name:        fmt.Sprintf("Trump%s.png", options[0].StringValue()),
					Reader:      sendFile,
				},
			},
		},
	})
}
