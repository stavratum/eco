package main

import (
	"github.com/bwmarrin/discordgo"
)

type Session struct {
	Handlers []func(*discordgo.Session, interface{})

	*discordgo.Session
}

func (eco *Session) AddHandler(h func(*discordgo.Session, interface{})) {
	eco.Handlers = append(eco.Handlers, h)
}

func (eco *Session) handleHandlers(s *discordgo.Session, e interface{}) {
	for _, handler := range eco.Handlers {
		handler(s, e)
	}
}

func New(token string) (s *Session, err error) {
	s = &Session{
		Handlers: []func(*discordgo.Session, interface{}){},
	}

	s.Session, err = discordgo.New(token)
	if err != nil {
		return
	}

	s.Session.AddHandler(s.handleHandlers)
	return
}
