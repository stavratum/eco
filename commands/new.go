package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stavratum/eco/handler/router"
)

var New = &router.Command{
	Name:        "new",
	Usage:       "new",
	Description: "Recreates this channel with the same data as before.",

	Aliases: []string{"new"},

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

		channel, err := ctx.State.Channel(ctx.ChannelID)
		if err != nil {
			if channel, err = ctx.Channel(ctx.ChannelID); err != nil {
				ctx.ChannelMessageSendReply(ctx.ChannelID, err.Error(), ctx.Reference())
				return
			}

			ctx.State.ChannelAdd(channel)
		}

		_, err = ctx.GuildChannelCreateComplex(ctx.GuildID, discordgo.GuildChannelCreateData{
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
			ctx.ChannelMessageSendReply(ctx.ChannelID, err.Error(), ctx.Reference())
			return
		}

		if _, err = ctx.ChannelDelete(ctx.ChannelID); err != nil {
			ctx.ChannelMessageSendReply(ctx.ChannelID, err.Error(), ctx.Reference())
			return
		}
	},
}
