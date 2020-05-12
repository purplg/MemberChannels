package environment

import "github.com/sirupsen/logrus"

const (
	defaultDBFile = "./db"
	defaultCategoryName = "Dynamic Channels"
	defaultChannelName = "[ + ] [ Create channel ]"
)

type Environment struct {
	// Required
	DiscordAPIToken string

	// Optional
	LogLevel            logrus.Level
	DBFile              string
	DefaultCategoryName string
	DefaultChannelName  string
}

func New(discordToken string) *Environment {
	if len(discordToken) == 0 {
		logrus.Fatalln("Missing discord api token")
	}
	return &Environment{
		DiscordAPIToken:     discordToken,
		LogLevel:            logrus.WarnLevel,
		DBFile:              defaultDBFile,
		DefaultCategoryName: defaultCategoryName,
		DefaultChannelName:  defaultChannelName,
	}
}
