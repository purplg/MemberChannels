package environment

import (
	"errors"
	"flag"
	"os"
)

const (
	defaultDBDirectory  = "./db"
	defaultCategoryName = "Member Channels"
	defaultListenName   = "[ + New ]"
)

type Environment struct {
	// Required
	DiscordAPIToken string

	// Optional
	Verbose             bool
	DBFile              string
	DefaultCategoryName string
	DefaultListenName   string
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
	flag.BoolVar(&env.Verbose, "v", false, "Enable debug logging")
	flag.StringVar(&env.DBFile, "store-dir", defaultDBDirectory, "Directory to save database")
	flag.StringVar(&env.DefaultCategoryName, "category-channel-name", defaultCategoryName, "The default name for created the created cateogories")
	flag.StringVar(&env.DefaultListenName, "listen-channel-name", defaultListenName, "The default name for created the created cateogories")
	flag.Parse()
}
