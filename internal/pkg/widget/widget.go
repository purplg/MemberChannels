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

func (w *Widget) UserVoiceEvent(event *discordgo.VoiceStateUpdate) {
	w.log.Debugln("UserVoiceEvent")

	if prevChannel, ok := w.currentChannel[event.UserID]; ok {
		w.log.Debugf("Found previous channel: %d\n", prevChannel.userCount)
		w.UserLeftChannel(event.UserID, prevChannel)
	}

	if event.ChannelID == "" {
		return
	}

	// Create new
	if event.ChannelID == w.listenChannel.ID {
		w.log.Debugln("Creating new channel")
		w.UserNewChannel(event.UserID)
		return
	}

	// Join
	if existingChannel, ok := w.activeChannels[event.ChannelID]; ok {
		w.log.Debugln("Joining existing channel")
		w.UserJoinedChannel(event.UserID, existingChannel)
		w.log.Debugf("UserCount: %d\n", existingChannel.userCount)
	}
}

func (w *Widget) UserLeftChannel(userID string, uc *userChannel) {
	delete(w.currentChannel, userID)
	uc.userCount--
	if uc.userCount == 0 {
		w.log.Debugln("Empty. Deleting")
		w.session.ChannelDelete(uc.ID)
	}
}

func (w *Widget) UserNewChannel(userID string) {
	if userChannel, err := w.newUserChannel(userID); err != nil {
		w.log.WithError(err).Errorln("Error creating new user channel")
	} else {
		w.activeChannels[userChannel.ID] = userChannel
		w.session.GuildMemberMove(userChannel.GuildID, userID, userChannel.ID)
	}
}

func (w *Widget) UserJoinedChannel(userID string, uc *userChannel) {
	w.currentChannel[userID] = uc
	uc.userCount++
}

func (w *Widget) ChannelChanged(channel *discordgo.Channel) {
	if w.isListenChannel(channel) {
		w.guildDB.SetChannelName(channel.Name)
		w.log.WithField("channelName", channel.Name).Debugln("New listen channel name")
	} else if userChan, ok := w.activeChannels[channel.ID]; ok {
		w.guildDB.SetUserChannel(userChan.owner.ID, channel.ID, channel.Name)
		w.log.WithField("channelName", channel.Name).Debugln("New user channel name")
	}
}

func (w *Widget) IsManagedChannel(channel *discordgo.Channel) bool {
	return channel.ParentID == w.categoryChannel.ID
}

func (w *Widget) Close() {
	w.session.ChannelDelete(w.listenChannel.ID)
}

type Widget struct {
	session *discordgo.Session
	log     *logrus.Entry
	guildDB database.GuildDatabase

	categoryChannel *discordgo.Channel
	listenChannel   *discordgo.Channel

	currentChannel map[string]*userChannel // map[userID] -> userChannel
	activeChannels map[string]*userChannel // map[channelID] -> userChannel
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
func New(session *discordgo.Session, log *logrus.Entry, guildDB database.GuildDatabase, data *WidgetData) (*Widget, error) {
	w := &Widget{
		session:         session,
		log:             log,
		guildDB:         guildDB,
		categoryChannel: nil,
		listenChannel:   nil,
		currentChannel:  make(map[string]*userChannel),
		activeChannels:  make(map[string]*userChannel),
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
func (w *Widget) newUserChannel(userID string) (*userChannel, error) {
	// Look up the saved channel name for user
	channelName := w.guildDB.UserChannelName(userID)

	user, err := w.session.User(userID)
	if err != nil {
		w.log.WithError(err).Errorf("Could not resolve userID: %s\n", userID)
	}

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
				ID:    userID,
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
func (w *Widget) getCategory(categoryID, defaultName string) (*discordgo.Channel, error) {
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
func (w *Widget) getListenChannel(channelID, defaultName, parentID string) (*discordgo.Channel, error) {
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

func (w *Widget) isListenChannel(channel *discordgo.Channel) bool {
	return channel.ID == w.listenChannel.ID
}
