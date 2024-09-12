package commands

import (
	"eco/eco"
	"eco/handler"

	"github.com/bwmarrin/discordgo"
)

var New = &handler.Command{
	Name:        "new",
	Usage:       "new",
	Description: "Recreates this channel with the same data as before.",

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

		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			if channel, err = s.Channel(m.ChannelID); err != nil {
				s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
				return
			}

			s.State.ChannelAdd(channel)
		}

		nchannel, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name:                 channel.Name,
			Type:                 channel.Type,
			Topic:                channel.Topic,
			Bitrate:              channel.Bitrate,
			UserLimit:            channel.UserLimit,
			RateLimitPerUser:     channel.RateLimitPerUser,
			Position:             channel.Position,
			PermissionOverwrites: channel.PermissionOverwrites,
			ParentID:             channel.ParentID,
			NSFW:                 channel.NSFW,
		})

		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
			return
		}

		if guild, ok := s.Guilds[m.GuildID]; ok {
			if _, ok = guild.TempChannels[m.ChannelID]; ok {
				guild.TempChannels[nchannel.ID] = true
			}
		}

		if _, err = s.ChannelDelete(m.ChannelID); err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
			return
		}
	},
}
