package commands

import (
	"eco/eco"
	"eco/handler"
	"log"

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

		channel, ok := handler.TempChannels[m.ChannelID]
		if !ok {
			channel = handler.NewTempChannel()
			handler.TempChannels[m.ChannelID] = channel

			s.ChannelMessageSendReply(m.ChannelID, "Temp channel type: Manual.", m.Reference())
			return
		}

		channel.Type++
		content := "Temp channel type: Everything."

		switch channel.Type {
		case handler.Everything:
		case handler.EverythingNSFW:
			content = "Temp channel type: Everything (NSFW)."
		case 3:
			delete(handler.TempChannels, m.ChannelID)
			content = "Channel is not temp anymore."
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   content,
			Reference: m.Reference(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							// Label is what the user will see on the button.
							Label: "Yes ($999.99)",
							// Style provides coloring of the button. There are not so many styles tho.
							Style: discordgo.SuccessButton,
							// Disabled allows bot to disable some buttons for users.
							Disabled: false,
							// CustomID is a thing telling Discord which data to send when this button will be pressed.
							CustomID: "fd_yes",

							Emoji: &discordgo.ComponentEmoji{
								Name: "🤑",
							},
						},
						discordgo.Button{
							Label:    "No ($299.99)",
							Style:    discordgo.DangerButton,
							Disabled: false,
							CustomID: "fd_no",

							Emoji: &discordgo.ComponentEmoji{
								Name: "😡",
							},
						},
					},
				},
			},
		})

		if err != nil {
			log.Println(err.Error())
		}
	},
}
