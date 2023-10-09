package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/SqueakyBeaver/botkisser"
	"github.com/SqueakyBeaver/botkisser/commands"
)

var (
	shouldSyncCommands *bool
	version            = "dev"
)

func init() {
	shouldSyncCommands = flag.Bool("sync-commands", false, "Whether to sync commands to discord")
	flag.Parse()
}

func main() {
	cfg, err := dbot.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Starting bot version: %s", version)
	logger.Printf("Syncing commands? %t", *shouldSyncCommands)

	b := dbot.New(logger, version, *cfg)

	b.SetupBot()

	b.Session.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	commands.SetupCommands(b)

	err = b.Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")

	if *shouldSyncCommands {
		commands.SyncCommands(b, cfg, logger)
	}
	defer b.Session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	log.Println("Gracefully shutting down.")
}
