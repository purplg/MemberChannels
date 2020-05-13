package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) VoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {

	log := config.Log.WithField("GuildID", event.GuildID)

	guildDB := config.DB.AsGuild(event.GuildID)
	listeningChannelID := guildDB.ChannelID()

	widget := config.Widgets[guildDB.GuildID()]
	if widget == nil {
		log.Errorln("Could not find widget for guild")
		return
	}

	// When a user is connecting
	if event.ChannelID != "" {
		if event.ChannelID != listeningChannelID {
			return
		}

		user, err := session.User(event.UserID)
		if err != nil {
			log.WithError(err).Errorln("Could not lookup user on channel join")
			return
		}
		widget.NewChannel(user)
	}

	// Only care about channels under managed category
	var (
		channel *discordgo.Channel
		err     error
	)
	if channel, err = session.State.Channel(listeningChannelID); err != nil {
		if channel, err = session.Channel(listeningChannelID); err != nil {
			log.WithError(err).Errorln("Could not lookup category channel")
			return
		}
	}

	categoryID := guildDB.CategoryID()
	if channel.ParentID != categoryID {
		return
	}

	if err := widget.Sweep(); err != nil {
		log.WithError(err).Errorln("Deleting empty channels")
	}
}
