package variables

import "github.com/sirupsen/logrus"

const (
	defaultDBFile = "./db"
	defaultCategoryName = "Dynamic Channels"
	defaultChannelName = "[ + ] [ Create channel ]"
)

type Variables struct {
	// Required
	DiscordAPIToken string

	// Optional
	LogLevel            logrus.Level
	DBFile              string
	DefaultCategoryName string
	DefaultChannelName  string
}

func New(discordToken string) *Variables {
	if len(discordToken) == 0 {
		logrus.Fatalln("Missing discord api token")
	}
	return &Variables{
		DiscordAPIToken:     discordToken,
		LogLevel:            logrus.WarnLevel,
		DBFile:              defaultDBFile,
		DefaultCategoryName: defaultCategoryName,
		DefaultChannelName:  defaultChannelName,
	}
}
