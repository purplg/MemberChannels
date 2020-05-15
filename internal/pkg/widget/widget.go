package widget

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
)

// A Widget consists of a category channel and a designated Voice channel (listen channel)
// The listen channel waits for a user to join.
// When a user joins, a userChannel is created which:
// - Is a voice channel parented to the category channel
// - Named after the user (default) or by the user
// - Gives the user addition permissions the userChannel
// Then the user is automatically moved into this new channel

type Widget interface {
	UserJoinedManagedChannel(user *discordgo.User, channel *discordgo.Channel)
	UserDisconnected(userID string)
	ManagedChannelChanged(channel *discordgo.Channel)

	IsManagedChannel(channel *discordgo.Channel) bool

	Close()
}

func (w *widget) UserJoinedManagedChannel(user *discordgo.User, channel *discordgo.Channel) {
	log := w.log.WithFields(logrus.Fields{
		"user":    user.Username,
		"channel": channel.Name,
	})
	log.Debugln("UserJoinedManagedChannel")

	if _, ok := w.currentChannel[user.ID]; ok {
		w.UserDisconnected(user.ID)
	}

	if w.isListenChannel(channel) {
		if userChan, err := w.newUserChannel(user); err != nil {
			log.WithError(err).Errorln("Error creating User Channel")
		} else {
			w.userChannels[userChan.ID] = userChan
			w.session.GuildMemberMove(w.guildDB.GuildID(), user.ID, userChan.ID)
		}
	} else {
		if userChan, ok := w.userChannels[channel.ID]; ok {
			w.currentChannel[user.ID] = userChan
			userChan.userCount++
		}
	}
}

func (w *widget) UserDisconnected(userID string) {
	log := w.log.WithFields(logrus.Fields{
		"userID": userID,
	})
	log.Debugln("UserDisconnected")

	if userChan, ok := w.currentChannel[userID]; ok {
		userChan.userCount--
		log.Debugln("User Channel found")
		if userChan.userCount == 0 {
			log.Debugln("User Channel deleted")
			w.session.ChannelDelete(userChan.ID)
		}
	} else {
		log.Debugln("User Channel doesn't exist")
	}
}

func (w *widget) ManagedChannelChanged(channel *discordgo.Channel) {
	if w.isListenChannel(channel) {
		w.guildDB.SetChannelName(channel.Name)
		w.log.WithField("channelName", channel.Name).Debugln("New listen channel name")
	} else if userChan, ok := w.userChannels[channel.ID]; ok {
		w.guildDB.SetUserChannel(userChan.owner.ID, channel.ID, channel.Name)
		w.log.WithField("channelName", channel.Name).Debugln("New user channel name")
	}
}

func (w *widget) IsManagedChannel(channel *discordgo.Channel) bool {
	return channel.ParentID == w.categoryChannel.ID
}

func (w *widget) Close() {
	w.session.ChannelDelete(w.listenChannel.ID)
}

type widget struct {
	session *discordgo.Session
	log     *logrus.Entry
	guildDB database.GuildDatabase

	categoryChannel *discordgo.Channel
	listenChannel   *discordgo.Channel

	currentChannel map[string]*userChannel // map[userID] -> userChannel
	userChannels   map[string]*userChannel // map[channelID] -> userChannel
}

type userChannel struct {
	*discordgo.Channel
	owner     *discordgo.User
	userCount uint8
}

// Only used to initialize a new Widget
type WidgetData struct {
	CategoryID        string
	CategoryName      string
	ListenChannelID   string
	ListenChannelName string
}

// Just initialize values to prepare the widget
func New(session *discordgo.Session, log *logrus.Entry, guildDB database.GuildDatabase, data *WidgetData) (Widget, error) {
	w := &widget{
		session:         session,
		log:             log,
		guildDB:         guildDB,
		categoryChannel: nil,
		listenChannel:   nil,
		currentChannel:  make(map[string]*userChannel),
		userChannels:    make(map[string]*userChannel),
	}

	category, err := w.getCategory(data.CategoryID, data.CategoryName)
	if err != nil {
		return nil, err
	}
	w.categoryChannel = category

	listenChannel, err := w.getListenChannel(data.ListenChannelID, data.ListenChannelName, w.categoryChannel.ID)
	if err != nil {
		return nil, err
	}
	w.listenChannel = listenChannel

	w.guildDB.SetCategoryID(w.categoryChannel.ID)
	w.guildDB.SetCategoryName(w.categoryChannel.Name)
	w.guildDB.SetChannelID(w.listenChannel.ID)
	w.guildDB.SetChannelName(w.listenChannel.Name)

	return w, nil
}

// Create a new channel for user
func (w *widget) newUserChannel(user *discordgo.User) (*userChannel, error) {
	// Look up the saved channel name for user
	channelName := w.guildDB.UserChannelName(user.ID)

	// Or generate one if none found
	if channelName == "" {
		channelName = fmt.Sprintf("%s's channel", user.Username)
	}

	// Send API request to create the voice channel
	channel, err := w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
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
		ParentID: w.categoryChannel.ID,
	})
	if err != nil {
		return nil, err
	}

	return &userChannel{
		Channel:   channel,
		owner:     user,
		userCount: 0,
	}, nil
}

// Attempt to find existing category channel or create a new one
func (w *widget) getCategory(categoryID, defaultName string) (*discordgo.Channel, error) {
	if category, err := w.session.Channel(categoryID); err == nil {
		return category, nil
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:      defaultName,
		Type:      discordgo.ChannelTypeGuildCategory,
		UserLimit: 1,
	})
}

// Attempt to find existing listen channel or create a new one
func (w *widget) getListenChannel(channelID, defaultName, parentID string) (*discordgo.Channel, error) {
	if channel, err := w.session.Channel(channelID); err == nil {
		return channel, nil
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:      defaultName,
		Type:      discordgo.ChannelTypeGuildVoice,
		UserLimit: 1,
		ParentID:  parentID,
	})
}

func (w *widget) isListenChannel(channel *discordgo.Channel) bool {
	return channel.ID == w.listenChannel.ID
}
