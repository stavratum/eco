package handler

import (
	"eco/eco"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnTemp(s *eco.Session, i interface{}) {
	switch e := i.(type) {
	case *discordgo.ChannelDelete:
		guild, ok := s.Guilds[e.GuildID]
		if !ok {
			return
		}

		_, ok = guild.TempChannels[e.ID]
		if !ok {
			return
		}

		delete(guild.TempChannels, e.ID)
	case *discordgo.MessageCreate:
		guild, ok := s.Guilds[e.GuildID]
		if !ok {
			return
		}

		temp, ok := guild.TempChannels[e.ChannelID]
		if !ok {
			return
		}

		if !temp {
			return
		}

		time.AfterFunc(1*time.Hour, func() {
			temp, ok := guild.TempChannels[e.ChannelID]
			if !ok {
				return
			}

			if !temp {
				return
			}

			s.ChannelMessageDelete(e.ChannelID, e.ID)
		})
	}
}
