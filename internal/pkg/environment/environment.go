package environment

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	defaultDBDirectory  = "./db"
	defaultCategoryName = "Dynamic Channels"
	defaultListenName  = "[ + ] [ Create channel ]"
)

type Environment struct {
	// Required
	DiscordAPIToken string

	// Optional
	LogLevel            logrus.Level
	DBFile              string
	DefaultCategoryName string
	DefaultListenName  string
}

func New() (*Environment, error) {
	env := &Environment{}
	env.optional()
	if err := env.required(); err != nil {
		return nil, err
	}
	return env, nil
}

func (env *Environment) required() error {
	discordToken := os.Getenv("DISCORD_TOKEN")
	if len(discordToken) == 0 {
		return errors.New("Missing Discord API token. Set env var DISCORD_TOKEN")
	}
	env.DiscordAPIToken = discordToken
	return nil
}

func (env *Environment) optional() {
	loglevel := flag.String("loglevel", "WARN", "(DEBUG|INFO|WARN|ERROR)")
	flag.StringVar(&env.DBFile, "store-dir", defaultDBDirectory, "Directory to save database")
	flag.StringVar(&env.DefaultCategoryName, "category-channel-name", defaultCategoryName, "The default name for created the created cateogories")
	flag.StringVar(&env.DefaultListenName, "listen-channel-name", defaultListenName, "The default name for created the created cateogories")
	flag.Parse()
	env.LogLevel = parseLogLevel(*loglevel)
}

func parseLogLevel(level string) logrus.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	default:
		return logrus.WarnLevel
	}
}
