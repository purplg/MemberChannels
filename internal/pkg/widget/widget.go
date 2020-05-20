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
	GuildDB database.GuildDatabase

	categoryChannel *discordgo.Channel
	listenChannel   *discordgo.Channel

	currentChannel map[string]*userChannel // map[userID] -> userChannel
	activeChannels map[string]*userChannel // map[channelID] -> userChannel
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
		GuildDB:         guildDB,
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
		w.categoryChannel, err = w.session.GuildChannelCreateComplex(w.GuildDB.GuildID(), categoryChannelData(data.CategoryName))
		if err != nil {
			return err
		}
	}

	if w.listenChannel, err = w.session.Channel(data.ListenChannelID); err != nil {
		w.listenChannel, err = w.session.GuildChannelCreateComplex(w.GuildDB.GuildID(), listenChannelData(data.ListenChannelName, w.categoryChannel.ID))
		if err != nil {
			return err
		}
	}

	w.GuildDB.SetCategoryID(w.categoryChannel.ID)
	w.GuildDB.SetCategoryName(w.categoryChannel.Name)
	w.GuildDB.SetChannelID(w.listenChannel.ID)
	w.GuildDB.SetChannelName(w.listenChannel.Name)

	return nil
}

func (w *Widget) UserJoined(userID, channelID string) {
	if activeChannel, ok := w.activeChannels[channelID]; ok {
		w.currentChannel[userID] = activeChannel
		if activeChannel.ownerID != userID {
			activeChannel.AddVisitor(userID)
		}
		w.log.Debugf("VisitorCount: %d\n", len(activeChannel.visitorIDs))
	}
}

func (w *Widget) UserLeft(userID string) {
	prevChannel, ok := w.currentChannel[userID]
	if !ok {
		return
	}

	delete(w.currentChannel, userID)

	if prevChannel.ownerID == userID {
		w.log.Debugln("User is owner")
		if prevChannel.PopToOwner() {
			w.log.Debugln("Popping to new owner")
			channelName := w.GuildDB.UserChannelName(prevChannel.ownerID)
			w.session.ChannelEditComplex(prevChannel.ID, changeOwnerChannelData(channelName, prevChannel.ownerID))
		} else {
			w.log.Debugln("Empty. Deleting")
			w.session.ChannelDelete(prevChannel.ID)
		}
	} else {
		prevChannel.RemoveVisitor(userID)
	}
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

func (w *Widget) RenameUserChannel(channelID, channelName string) {
	userChan, ok := w.activeChannels[channelID]
	if !ok {
		return
	}

	w.GuildDB.SetUserChannel(userChan.ownerID, userChan.ID, channelName)
	w.log.WithField("channelName", channelName).Debugln("New user channel name")
}

func (w *Widget) RenameListenChannel(channelName string) {
	w.GuildDB.SetChannelName(channelName)
}

func (w *Widget) Close() {
	w.session.ChannelDelete(w.listenChannel.ID)
}

// Returns true if the given `channelID` is the listen channel for this widget
func (w *Widget) IsListenChannel(channelID string) bool {
	return channelID == w.listenChannel.ID
}

// Returns true if the given `channelID` is a managed user channel for this widget
func (w *Widget) IsUserChannel(channelID string) bool {
	for activeID := range w.activeChannels {
		if channelID == activeID {
			return true
		}
	}
	return false
}

// Create a new channel for user
func (w *Widget) newUserChannel(userID string) (*userChannel, error) {
	// Look up the saved channel name for user
	channelName := w.GuildDB.UserChannelName(userID)

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
	channel, err := w.session.GuildChannelCreateComplex(w.GuildDB.GuildID(), userChannelData(channelName, userID, w.categoryChannel.ID))
	if err != nil {
		return nil, err
	}

	return &userChannel{
		Channel:    channel,
		ownerID:    user.ID,
		visitorIDs: []string{},
	}, nil
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

func changeOwnerChannelData(channelName, ownerID string) *discordgo.ChannelEdit {
	return &discordgo.ChannelEdit{
		Name: channelName,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    ownerID,
				Type:  "member",
				Allow: 16,
				Deny:  0,
			},
		},
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
