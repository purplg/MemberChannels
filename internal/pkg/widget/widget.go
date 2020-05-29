package widget

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
)

// A Widget consists of a category channel and a designated Voice channel (listen channel)
// The listen channel waits for a user to join.
// When a user joins, a memberChannel is created which:
// - Is a voice channel parented to the category channel
// - Named after the user (default) or by the user
// - Gives the user addition permissions the memberChannel
// Then the user is automatically moved into this new channel

type Widget struct {
	session *discordgo.Session
	log     *logrus.Entry
	GuildDB database.GuildDatabase

	guildID string

	categoryChannel *discordgo.Channel
	listenChannel   *discordgo.Channel

	currentChannel map[string]*memberChannel // map[userID] -> memberChannel
	activeChannels map[string]*memberChannel // map[channelID] -> memberChannel
}

func New(session *discordgo.Session, log *logrus.Entry, guildDB database.GuildDatabase) *Widget {
	return &Widget{
		session:         session,
		log:             log,
		GuildDB:         guildDB,
		guildID:         guildDB.GuildID(),
		categoryChannel: nil,
		listenChannel:   nil,
		currentChannel:  make(map[string]*memberChannel),
		activeChannels:  make(map[string]*memberChannel),
	}
}

func (w *Widget) Spawn(categoryID, categoryName, listenChannelName string) error {
	var err error

	// Resolve existing categoryChannel or create a new one
	if w.categoryChannel, err = w.session.State.Channel(categoryID); err != nil {
		if w.categoryChannel, err = w.session.Channel(categoryID); err != nil {
			w.categoryChannel, err = w.session.GuildChannelCreateComplex(w.GuildDB.GuildID(), categoryChannelData(categoryName))
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Resolved channel from API")
		}
	} else {
		fmt.Println("Resolved channel from State")
	}

	w.listenChannel, err = w.session.GuildChannelCreateComplex(w.GuildDB.GuildID(), listenChannelData(listenChannelName, w.categoryChannel.ID))
	if err != nil {
		return err
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

func (w *Widget) UserLeft(userID, channelID string) {
	channel := w.activeChannels[channelID]

	delete(w.currentChannel, userID)

	if channel == nil {
		return
	}

	if channel.ownerID == userID {
		w.log.Debugln("User is owner")
		if channel.PopToOwner() {
			w.log.Debugln("Popping to new owner")
			channelName := w.GuildDB.MemberChannelName(channel.ownerID)
			w.session.ChannelEditComplex(channel.ID, changeOwnerChannelData(channelName, channel.ownerID, 99))
		} else {
			w.log.Debugln("Empty. Deleting")
			w.session.ChannelDelete(channel.ID)
		}
	} else {
		channel.RemoveVisitor(userID)
	}
}

func (w *Widget) UserRequestChannel(userID string) {
	memberChannel, err := w.newMemberChannel(userID)
	if err != nil {
		w.log.WithError(err).Errorln("Error creating new member channel")
		return
	}

	w.activeChannels[memberChannel.ID] = memberChannel
	w.session.GuildMemberMove(w.guildID, userID, memberChannel.ID)
}

func (w *Widget) RenameMemberChannel(channelID, channelName string) {
	memberChan, ok := w.activeChannels[channelID]
	if !ok {
		return
	}

	w.GuildDB.SetMemberChannel(memberChan.ownerID, channelName)
	w.log.WithField("channelName", channelName).Debugln("New member channel name")
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

// Returns true if the given `channelID` is a managed member channel for this widget
func (w *Widget) IsMemberChannel(channelID string) bool {
	for activeID := range w.activeChannels {
		if channelID == activeID {
			return true
		}
	}
	return false
}

// Create a new channel for user
func (w *Widget) newMemberChannel(userID string) (*memberChannel, error) {
	// Look up the saved channel name for user
	channelName := w.GuildDB.MemberChannelName(userID)

	var user *discordgo.User

	if member, err := w.session.State.Member(w.guildID, userID); err == nil {
		user = member.User
		fmt.Println("Resolved user from State")
	} else if user, err = w.session.User(userID); err != nil {
		w.log.WithError(err).Errorf("Could not resolve userID: %s\n", userID)
		return nil, err
	} else {
		fmt.Println("Resolved user from API")
	}

	// Or generate one if none found
	if channelName == "" {
		channelName = fmt.Sprintf("%s's channel", user.Username)
	}

	// Send API request to create the voice channel
	channel, err := w.session.GuildChannelCreateComplex(w.GuildDB.GuildID(), memberChannelData(channelName, userID, w.categoryChannel.ID))
	if err != nil {
		return nil, err
	}

	return &memberChannel{
		ID:         channel.ID,
		ownerID:    user.ID,
		visitorIDs: []string{},
	}, nil
}

func memberChannelData(channelName, userID, parentID string) discordgo.GuildChannelCreateData {
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

func changeOwnerChannelData(channelName, ownerID string, position int) *discordgo.ChannelEdit {
	return &discordgo.ChannelEdit{
		Name:     channelName,
		Position: position,
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
