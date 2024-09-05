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

/*

func New(file string) (s *Session, err error) {
	L := lua.NewState()
	defer L.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	L.SetContext(ctx)
	if err = L.DoFile(file); err != nil {
		return
	}

	s = &Session{}

	switch lToken := L.Env.RawGetString("Token"); lToken.Type() {
	case lua.LTString:
		s.Session, err = discordgo.New(lToken.String())
		if err != nil {
			return
		}

		s.Session.AddHandler(s.handleHandlers)
	case lua.LTNil:
		err = fmt.Errorf("%s: Token must be defined", file)
		return
	default:
		err = fmt.Errorf("%s: Token must be a string", file)
		return
	}

	switch prefix := L.Env.RawGetString("Prefix"); prefix.Type() {
	case lua.LTString:
		s.Prefix = prefix.String()
	case lua.LTNil:
		s.Prefix = ","
	default:
		err = fmt.Errorf("%s: Prefix must be a string", file)
	}

	return
}

*/
