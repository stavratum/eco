package commands

import (
	"eco/eco"
	"eco/handler"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var TempDescription = `
> **%s**
> ` + "`%s`" + `

**Duration**: ` + "`%s`" + `
**Updated**: %s 
`

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

		channel, ok := handler.TempChannels[m.ChannelID]
		if ok {
			if channel.Schedule++; channel.Schedule > handler.ChannelScheduleNSFW {
				channel.Schedule = handler.NoSchedule
			}

			channel.Time = time.Now()
			channel.Messages = map[string]time.Time{}
		} else {
			channel = handler.NewTempChannel()
			handler.TempChannels[m.ChannelID] = channel
		}

		s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
			Description: fmt.Sprintf(TempDescription, channel.Schedule.String(), channel.Schedule.Description(), channel.Duration.String(), "<t:"+strconv.FormatInt(channel.Time.Unix(), 10)+":R>"),
			Color:       0xff5858,
		}, m.Reference())
	},
}
