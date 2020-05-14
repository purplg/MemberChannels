package widget

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
)

type Widget interface {
	UserSwitchedChannel(user *discordgo.User, channel *discordgo.Channel)
	Close()
}


// A Widget consists of a Category channel and a designated Voice channel
// The Voice channel waits for a User to join. When a user joins,
// the Widget will:
// - Create a new voice channel with parameters:
//   - Name the channel appropriately
//   - Permissions so that the creating user has more control
//   - Parented to the Category channel
// - Save the channel name in case it was generated
// - Move player to new channel
type widget struct {
	session *discordgo.Session
	log     *logrus.Entry
	guildDB database.GuildDatabase

	categoryChannel *discordgo.Channel
	listenChannel   *discordgo.Channel

	currentChannel map[string]*discordgo.Channel // map[userID] -> channel
	userChannels   map[string]*userChannel       // map[channelID] -> userChannel
}

// Only used to initialize a new Widget
type WidgetData struct {
	CategoryID        string
	CategoryName      string
	ListenChannelID   string
	ListenChannelName string
}

type userChannel struct {
	channel   *discordgo.Channel
	owner     *discordgo.User
	userCount uint8
}

// Just initialize values to prepare the widget
func New(session *discordgo.Session, log *logrus.Entry, guildDB database.GuildDatabase, data *WidgetData) (Widget, error) {
	w := &widget{
		session:         session,
		log:             log,
		guildDB:         guildDB,
		categoryChannel: nil,
		listenChannel:   nil,
		currentChannel:  make(map[string]*discordgo.Channel),
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
		channel:   channel,
		owner:     user,
		userCount: 0,
	}, nil
}

func (w *widget) userLeftChannel(channel *discordgo.Channel) {
	fmt.Println("userLeftChannel")
	if userChannel, ok := w.userChannels[channel.ID]; ok {
		userChannel.userCount--
		if userChannel.userCount <= 0 {
			delete(w.userChannels, channel.ID)
			w.session.ChannelDelete(channel.ID)
		}
	}
}

func (w *widget) UserSwitchedChannel(user *discordgo.User, channel *discordgo.Channel) {
	fmt.Println("userSwitchedChannel")
	if lastChannel, ok := w.currentChannel[user.ID]; ok {
		fmt.Println("has previous channel")
		w.userLeftChannel(lastChannel)
	}

	switch {
	case channel == nil:
		fmt.Println("channel == nil")
		delete(w.currentChannel, user.ID)

	case w.isListenChannel(channel):
		fmt.Println("isListenChannel")
		if userChannel, err := w.newUserChannel(user); err != nil {
			w.log.WithError(err).Errorln("Failed to create user channel")
		} else {
			w.userChannels[userChannel.channel.ID] = userChannel
			w.guildDB.SetUserChannel(user.ID, userChannel.channel.ID, userChannel.channel.Name)
			w.session.GuildMemberMove(w.guildDB.GuildID(), user.ID, userChannel.channel.ID)
		}

	case w.isUserChannel(channel):
		fmt.Println("isUserChannel")
		w.currentChannel[user.ID] = channel
		w.userChannels[channel.ID].userCount++
	}

}

func (w *widget) Close() {
	w.session.ChannelDelete(w.listenChannel.ID)
}
