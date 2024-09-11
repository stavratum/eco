package commands

import (
	"eco/eco"
	"eco/handler"

	"github.com/bwmarrin/discordgo"
)

var Temp = &handler.Command{
	Name:        "temp",
	Usage:       "temp",
	Description: "Sets this channel to be temporary or not if already is.",

	Handler: func(s *eco.Session, m *discordgo.MessageCreate) {
		permissions, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
			return
		}

		if permissions&discordgo.PermissionManageChannels == 0 {
			s.ChannelMessageSendReply(m.ChannelID, "Not enough permissions.", m.Reference())
			return
		}

		guild, ok := s.Guilds[m.GuildID]
		if !ok {
			guild = eco.NewGuild()
			s.Guilds[m.GuildID] = guild
		}

		_, ok = guild.TempChannels[m.ChannelID]
		if !ok {
			guild.TempChannels[m.ChannelID] = true

			s.ChannelMessageSendReply(m.ChannelID, "Channel was set to be temporary.", m.Reference())
		} else {
			delete(guild.TempChannels, m.ChannelID)

			s.ChannelMessageSendReply(m.ChannelID, "Channel was set to not be temporary anymore.", m.Reference())
		}
	},
}
