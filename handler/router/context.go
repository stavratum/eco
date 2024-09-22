package router

import (
	"github.com/bwmarrin/discordgo"
)

// Context is a session, message and command data.
type Context struct {
	*discordgo.Session
	*discordgo.MessageCreate

	Task uint

	Command   string                    // Command trigger
	Body      []string                  // Command body
	Arguments map[*Argument]interface{} // Command arguments
}

func (ctx *Context) SendReply(content string) {
	ctx.Session.ChannelMessageSendReply(ctx.MessageCreate.ChannelID, content, ctx.MessageCreate.Reference())
}

// NewContext creates a new context.
func NewContext(s *discordgo.Session, m *discordgo.MessageCreate) *Context {
	return &Context{
		Session:       s,
		MessageCreate: m,

		Task: 0,

		Command:   "",
		Body:      []string{},
		Arguments: map[*Argument]interface{}{},
	}
}
