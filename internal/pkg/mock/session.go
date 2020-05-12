package mock

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

func MockSession() *discordgo.Session {
	return &discordgo.Session{
		State:        discordgo.NewState(),
		StateEnabled: true,
		Client:       restClient(),
		Ratelimiter:  discordgo.NewRatelimiter(),
	}
}
