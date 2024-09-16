package eco

import (
	"github.com/bwmarrin/discordgo"
)

type Session struct {
	Handlers []Handler

	*discordgo.Session
}

func New(token string) (s *Session, err error) {
	s = &Session{
		Handlers: []Handler{},
	}

	s.Session, err = discordgo.New(token)
	if err != nil {
		return
	}

	s.Session.AddHandler(s.handleHandlers)
	return
}
