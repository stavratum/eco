package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/stavratum/eco/handler"
	"github.com/stavratum/eco/handler/router"

	"github.com/bwmarrin/discordgo"
)

var TempDescription = `
> **%s**
> ` + "`%s`" + `

**Duration**: ` + "`%s`" + `
**Updated**: %s 
`

var Temp = &router.Command{
	Name:        "temp",
	Usage:       "temp",
	Description: "Sets this channel to be temporary or not if already is.",

	Aliases: []string{"temp"},

	Handler: func(ctx *router.Context) {
		permissions, err := ctx.UserChannelPermissions(ctx.Author.ID, ctx.ChannelID)
		if err != nil {
			ctx.ChannelMessageSendReply(ctx.ChannelID, err.Error(), ctx.Reference())
			return
		}

		if permissions&discordgo.PermissionManageChannels == 0 {
			ctx.ChannelMessageSendReply(ctx.ChannelID, "Not enough permissions.", ctx.Reference())
			return
		}

		channel, ok := handler.TempChannels[ctx.ChannelID]
		if ok {
			if channel.Schedule++; channel.Schedule > handler.ChannelScheduleNSFW {
				channel.Schedule = handler.NoSchedule
			}

			channel.Time = time.Now()
			channel.Messages = map[string]time.Time{}
		} else {
			channel = handler.NewTempChannel()
			handler.TempChannels[ctx.ChannelID] = channel
		}

		ctx.ChannelMessageSendEmbedReply(ctx.ChannelID, &discordgo.MessageEmbed{
			Description: fmt.Sprintf(TempDescription, channel.Schedule.String(), channel.Schedule.Description(), channel.Duration.String(), "<t:"+strconv.FormatInt(channel.Time.Unix(), 10)+":R>"),
			Color:       0xff5858,
		}, ctx.Reference())
	},
}
