package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"purplg.com/memberchannels/internal/pkg/database"
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

func (config *Events) onUserConnect(session *discordgo.Session, guildDB database.GuildDatabase, user *discordgo.User) error {
	userChannelName := guildDB.UserChannelName(user.ID)
	if userChannelName == "" {
		userChannelName = fmt.Sprintf("%s's channel", user.Username)
	}
	listeningChannelID := guildDB.ChannelID()
	listeningChannelName := guildDB.ChannelName()

	// Reconfigure listening channel to a user channel
	channel, err := session.ChannelEditComplex(listeningChannelID, &discordgo.ChannelEdit{
		Name: userChannelName,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    user.ID,
				Type:  "member",
				Allow: 16,
				Deny:  0,
			},
		},
	})
	if err != nil {
		return err
	}
	guildDB.SetUserChannel(user.ID, channel.ID, userChannelName)

	// Create new listening channel
	channel, err = session.GuildChannelCreateComplex(guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:     listeningChannelName,
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: guildDB.CategoryID(),
	})
	if err != nil {
		return err
	}

	guildDB.SetChannelID(channel.ID)
	return nil
}
