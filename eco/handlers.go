package eco

import "github.com/bwmarrin/discordgo"

type Handler func(*Session, interface{})

func (s *Session) AddHandler(h Handler) {
	s.Handlers = append(s.Handlers, h)
}

func (s *Session) handleHandlers(_ *discordgo.Session, e interface{}) {
	for _, handler := range s.Handlers {
		handler(s, e)
	}
}
