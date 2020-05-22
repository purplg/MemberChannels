package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"purplg.com/memberchannels/internal/pkg/database"
	"purplg.com/memberchannels/internal/pkg/environment"
	"purplg.com/memberchannels/internal/pkg/events"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)

func startDiscordSession(token string, evnts *events.Events) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	err = session.Open()
	if err != nil {
		return nil, err
	}

	session.AddHandler(evnts.GuildCreated)
	session.AddHandler(evnts.VoiceStateUpdate)
	session.AddHandler(evnts.ChannelUpdate)

	return session, nil
}

func main() {
	logger := logrus.New()
	log := logrus.NewEntry(logger)

	vars, err := environment.New()
	if err != nil {
		log.WithError(err).Fatal()
	}

	if vars.Verbose {
		log.Logger.SetLevel(logrus.DebugLevel)
	}

	data, err := database.Open(vars.DBFile, log)
	if err != nil {
		log.WithError(err).Fatalln()
	}

	config := events.New(vars, log, data)
	defer config.Close()

	session, err := startDiscordSession(vars.DiscordAPIToken, config)
	if err != nil {
		log.WithError(err).Errorln()
		return
	}
	defer session.Close()

	fmt.Printf("Version: %s\tBuilt at: %s\n", BuildVersion, BuildTime)
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
