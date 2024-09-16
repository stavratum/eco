package handler

import (
	"eco/eco"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	Manual         = iota // for manual deleting of messages (bulk)
	Everything            // planned for deleting the channel itself
	EverythingNSFW        // planned for deleting the channel itself, but making it nsfw beforehand (idk whats the point honestly i might remove)
)

// temporary channel that stores umm yea
type TempChannel struct {
	Type  uint8
	Timed time.Time

	Messages map[string]time.Time
}

func NewTempChannel() *TempChannel {
	return &TempChannel{
		Type:     Manual,
		Messages: map[string]time.Time{},
	}
}

var TempChannels = map[string]*TempChannel{}
var TempDuration = time.Minute               // how long temp messages or channels stay
var TempTicker = time.NewTicker(time.Minute) // update every minute
var TempInit = false

// ready controls the temp channels and messages
func ready(s *eco.Session) {
	if TempInit {
		return
	}
	TempInit = true

	for {
		<-TempTicker.C

		for cID, channel := range TempChannels {
			switch channel.Type {
			case Manual:
				var (
					now      = time.Now()
					messages = []string{}
				)

				for mID, timed := range channel.Messages {
					if !now.After(timed) {
						continue
					}

					if messages = append(messages, mID); len(messages) != 100 {
						continue
					}

					s.ChannelMessagesBulkDelete(cID, messages)
					messages = []string{}
				}

				s.ChannelMessagesBulkDelete(cID, messages)
			case Everything:
				// TODO: recreate a channel and delete old one

			case EverythingNSFW:
				// TODO: recreate a channel, make old one nsfw and delete it

			}
		}
	}
}

// when message gets created make it temp as long as its configured
func messageCreate(m *discordgo.MessageCreate) {
	channel, ok := TempChannels[m.ChannelID]
	if !ok {
		return
	}

	if channel.Type != Manual {
		return
	}

	channel.Messages[m.ID] = time.Now().Add(TempDuration)
}

// when message gets deleted forget it
func messageDelete(m *discordgo.MessageDelete) {
	channel, ok := TempChannels[m.ChannelID]
	if !ok {
		return
	}

	if channel.Type != 0 {
		return
	}

	delete(channel.Messages, m.ID)
}

// when channel deletes forget everything about it
func channelDelete(c *discordgo.ChannelDelete) {
	if _, ok := TempChannels[c.ID]; !ok {
		return
	}

	delete(TempChannels, c.ID)
}

func OnTemp(s *eco.Session, i interface{}) {
	switch e := i.(type) {
	case *discordgo.Ready:
		ready(s)
	case *discordgo.MessageCreate:
		messageCreate(e)
	case *discordgo.MessageDelete:
		messageDelete(e)
	case *discordgo.ChannelDelete:
		channelDelete(e)
	}
}
