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

func New(session *discordgo.Session, log *logrus.Entry, guildDB database.GuildDatabase) *Widget {
	return &Widget{
		session:         session,
		log:             log,
		guildDB:         guildDB,
		categoryChannel: nil,
		listenChannel:   nil,
		currentChannel:  make(map[string]*userChannel),
		activeChannels:  make(map[string]*userChannel),
	}
}

func (w *Widget) Spawn(data *WidgetData) error {
	var err error

	// Resolve existing categoryChannel or create a new one
	if w.categoryChannel, err = w.session.Channel(data.CategoryID); err != nil {
		w.categoryChannel, err = w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), categoryChannelData(data.CategoryName))
		if err != nil {
			return err
		}
	}

	if w.listenChannel, err = w.session.Channel(data.ListenChannelID); err != nil {
		w.listenChannel, err = w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), listenChannelData(data.ListenChannelName, w.categoryChannel.ID))
		if err != nil {
			return err
		}
	}

	w.guildDB.SetCategoryID(w.categoryChannel.ID)
	w.guildDB.SetCategoryName(w.categoryChannel.Name)
	w.guildDB.SetChannelID(w.listenChannel.ID)
	w.guildDB.SetChannelName(w.listenChannel.Name)

	return nil
}

func (w *Widget) UserJoined(event *discordgo.VoiceStateUpdate) {
	if existingChannel, ok := w.activeChannels[event.ChannelID]; ok {
		w.log.Debugln("Joining existing channel")
		w.currentChannel[event.UserID] = existingChannel
		existingChannel.userCount++
		w.log.Debugf("UserCount: %d\n", existingChannel.userCount)
	}
}

func (w *Widget) ChannelChangedEvent(channel *discordgo.Channel) {
	if w.IsListenChannel(channel.ID) {
		w.guildDB.SetChannelName(channel.Name)
		w.log.WithField("channelName", channel.Name).Debugln("New listen channel name")
	} else if userChan, ok := w.activeChannels[channel.ID]; ok {
		w.guildDB.SetUserChannel(userChan.owner.ID, channel.ID, channel.Name)
		w.log.WithField("channelName", channel.Name).Debugln("New user channel name")
	}
}

func (w *Widget) UserLeft(userID string) {
	prevChannel, ok := w.currentChannel[userID]
	if !ok {
		return
	}

	w.log.Debugf("Found previous channel: %d\n", prevChannel.userCount)
	delete(w.currentChannel, userID)
	prevChannel.userCount--
	if prevChannel.userCount == 0 {
		w.log.Debugln("Empty. Deleting")
		w.session.ChannelDelete(prevChannel.ID)
	}
}

func (w *Widget) Close() {
	w.session.ChannelDelete(w.listenChannel.ID)
}

func (w *Widget) UserRequestChannel(userID string) {
	userChannel, err := w.newUserChannel(userID)
	if err != nil {
		w.log.WithError(err).Errorln("Error creating new user channel")
		return
	}

	w.activeChannels[userChannel.ID] = userChannel
	w.session.GuildMemberMove(userChannel.GuildID, userID, userChannel.ID)
}

// Create a new channel for user
func (w *Widget) newUserChannel(userID string) (*userChannel, error) {
	// Look up the saved channel name for user
	channelName := w.guildDB.UserChannelName(userID)

	user, err := w.session.User(userID)
	if err != nil {
		w.log.WithError(err).Errorf("Could not resolve userID: %s\n", userID)
		return nil, err
	}

	// Or generate one if none found
	if channelName == "" {
		channelName = fmt.Sprintf("%s's channel", user.Username)
	}

	// Send API request to create the voice channel
	channel, err := w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), userChannelData(channelName, userID, w.categoryChannel.ID))
	if err != nil {
		return nil, err
	}

	return &userChannel{
		Channel:   channel,
		owner:     user,
		userCount: 0,
	}, nil
}

// Returns true if the given `channel` is managed by this widget
func (w *Widget) IsManaged(channel *discordgo.Channel) bool {
	return channel.ParentID == w.categoryChannel.ID || channel.ID == w.categoryChannel.ID
}

// Returns true if the given `channelID` is the listen channel for this widget
func (w *Widget) IsListenChannel(channelID string) bool {
	return channelID == w.listenChannel.ID
}

func userChannelData(channelName, userID, parentID string) discordgo.GuildChannelCreateData {
	return discordgo.GuildChannelCreateData{
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
		ParentID: parentID,
	}
}

func categoryChannelData(channelName string) discordgo.GuildChannelCreateData {
	return discordgo.GuildChannelCreateData{
		Name: channelName,
		Type: discordgo.ChannelTypeGuildCategory,
	}
}

func listenChannelData(channelName, parentID string) discordgo.GuildChannelCreateData {
	return discordgo.GuildChannelCreateData{
		Name:      channelName,
		Type:      discordgo.ChannelTypeGuildVoice,
		UserLimit: 1,
		ParentID:  parentID,
	}
}
