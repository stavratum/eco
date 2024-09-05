package commands

import (
	"eco/eco"
	"eco/handler"

	"github.com/bwmarrin/discordgo"
)

var Help = &handler.Command{
	Name:        "help",
	Usage:       "help [command]",
	Description: "Returns information about a command and it's features.",

	Handler: func(s *eco.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSendReply(m.ChannelID, "unimplemented", m.Reference())
	},
}
