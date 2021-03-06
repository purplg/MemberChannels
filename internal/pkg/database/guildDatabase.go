package database

import (
	"github.com/sirupsen/logrus"
	"github.com/xujiajun/nutsdb"
)

// Field names
const (
	userPrefix    = "USER_"
	U_CHANNELNAME = userPrefix + "channelName"

	guildPrefix    = "GUILD_"
	G_CATEGORYID   = guildPrefix + "categoryID"
	G_CATEGORYNAME = guildPrefix + "categoryName"
	G_CHANNELNAME  = guildPrefix + "channelName"
	G_CHANNELID    = guildPrefix + "channelID"
)

type GuildDatabase interface {
	Log() *logrus.Entry
	GuildID() string

	// User data
	MemberChannelName(userID string) string
	SetMemberChannel(userid, channelName string)

	// * Guild data
	CategoryID() string
	CategoryName() string
	ChannelID() string
	ChannelName() string
	SetCategoryID(string)
	SetCategoryName(string)
	SetChannelID(string)
	SetChannelName(string)
}

type guildDB struct {
	*DB

	log *logrus.Entry

	// Acts as bucket id for database
	guildID string
}

func (g guildDB) Log() *logrus.Entry {
	return g.log
}

func (g guildDB) GuildID() string {
	return g.guildID
}

// -----------------------------------------------------------------------------
// User data
// -----------------------------------------------------------------------------
func (g guildDB) MemberChannelName(userID string) string {
	return g.getValue(U_CHANNELNAME + userID)
}

func (g guildDB) SetMemberChannel(userID string, channelName string) {
	g.setValue(U_CHANNELNAME+userID, channelName)
}

// -----------------------------------------------------------------------------
// Guild data
// -----------------------------------------------------------------------------
func (g guildDB) CategoryID() string {
	return g.getValue(G_CATEGORYID)
}

func (g guildDB) CategoryName() string {
	return g.getValue(G_CATEGORYNAME)
}

func (g guildDB) ChannelID() string {
	return g.getValue(G_CHANNELID)
}

func (g guildDB) ChannelName() string {
	return g.getValue(G_CHANNELNAME)
}

func (g guildDB) SetCategoryID(value string) {
	g.setValue(G_CATEGORYID, value)
}

func (g guildDB) SetCategoryName(value string) {
	g.setValue(G_CATEGORYNAME, value)
}

func (g guildDB) SetChannelID(value string) {
	g.setValue(G_CHANNELID, value)
}

func (g guildDB) SetChannelName(value string) {
	g.setValue(G_CHANNELNAME, value)
}

// -----------------------------------------------------------------------------
// Util functions
// -----------------------------------------------------------------------------
func (g guildDB) getValue(key string) string {
	var value string
	if err := g.View(
		func(tx *nutsdb.Tx) error {
			if entry, err := tx.Get(g.guildID, []byte(key)); err != nil {
				return err
			} else {
				value = string(entry.Value)
				return nil
			}
		}); err != nil {
		return ""
	}

	return value
}

func (g guildDB) setValue(key string, value string) {
	if err := g.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put(g.guildID, []byte(key), []byte(value), 0); err != nil {
				return err
			}
			return nil
		}); err != nil {
		g.Log().WithError(err).WithFields(logrus.Fields{
			"Key":   key,
			"Value": value,
		}).Errorln("Error setting value in datastore")
	}
}
