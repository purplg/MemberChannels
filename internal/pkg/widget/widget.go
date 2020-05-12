package widget

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
)

type Widget struct {
	session *discordgo.Session
	guildDB *database.GuildDB
	log     *logrus.Entry

	categoryID   string
	categoryName string
	channelID    string
	channelName  string
}

func New(session *discordgo.Session, log *logrus.Entry, guildDB *database.GuildDB, defaultCategoryName, defaultChannelName string) *Widget {
	w := &Widget{
		session: session,
		guildDB: guildDB,
		log:     log.WithField("Widget", guildDB.GuildID),

		categoryID:   guildDB.CategoryID(),
		categoryName: guildDB.CategoryName(),
		channelID:    guildDB.ChannelID(),
		channelName:  guildDB.ChannelName(),
	}

	if w.categoryName == "" {
		w.categoryName = defaultCategoryName
	}
	if w.channelName == "" {
		w.channelName = defaultChannelName
	}

	return w
}

func (w *Widget) Show() error {
	category, err := w.getCategory()
	if err != nil {
		return err
	}
	w.categoryID = category.ID
	channel, err := w.getChannel(w.categoryID)
	if err != nil {
		return err
	}
	w.channelID = channel.ID

	w.guildDB.SetCategoryID(w.categoryID)
	w.guildDB.SetCategoryName(w.categoryName)
	w.guildDB.SetChannelID(w.channelID)
	w.guildDB.SetChannelName(w.channelName)

	return nil
}

func (w *Widget) NewChannel(user *discordgo.User) {
	channelName := w.guildDB.UserChannelName(user.ID)
	if channelName == "" {
		channelName = fmt.Sprintf("%s's channel", user.Username)
	}

	if channel, err := w.session.GuildChannelCreateComplex(w.guildDB.GuildID, discordgo.GuildChannelCreateData{
		Name: channelName,
		Type: discordgo.ChannelTypeGuildVoice,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    user.ID,
				Type:  "member",
				Allow: 16,
				Deny:  0,
			},
		},
		ParentID: w.categoryID,
	}); err != nil {
		w.log.WithError(err).Errorln("Failed to create user channel")
	} else {
		w.guildDB.SetUserChannel(user.ID, channel.ID, channelName)
		w.session.GuildMemberMove(w.guildDB.GuildID, user.ID, channel.ID)
	}
}

func (w *Widget) Sweep() error {
	var (
		guild    *discordgo.Guild
		channels []*discordgo.Channel
		err      error
	)

	if guild, err = w.session.State.Guild(w.guildDB.GuildID); err != nil {
		if guild, err = w.session.Guild(w.guildDB.GuildID); err != nil {
			return err
		}
	}

	channels, err = w.session.GuildChannels(w.guildDB.GuildID)
	if err != nil {
		return err
	}

	categoryID := w.guildDB.CategoryID()
	listeningChannelID := w.guildDB.ChannelID()

	// Loop through all channels in guild
	for _, channel := range channels {
		if channel.ID == listeningChannelID {
			continue
		}
		// If this is a managed channel
		if channel.ParentID == categoryID {
			// Count the number of users in it
			userCount := 0
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == channel.ID {
					userCount++
				}
			}
			// If it's empty, remove it
			if userCount == 0 {
				w.session.ChannelDelete(channel.ID)
			}
		}
	}
	return nil
}

func (w *Widget) Close() {
	w.Sweep()
	w.session.ChannelDelete(w.channelID)
}
