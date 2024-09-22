package handler

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Schedule uint8

const (
	NoSchedule                Schedule = iota
	MessageSchedule                    // for deleting all messages
	MessageAttachmentSchedule          // for deleting all messages but with attachments
	ChannelSchedule                    // for deleting the channel itself
	ChannelScheduleNSFW                // for deleting the channel itself, but making it nsfw beforehand (idk whats the point honestly i might remove)
)

func (s Schedule) String() string {
	switch s {
	case NoSchedule:
		return "No Schedule"
	case MessageSchedule:
		return "Message Schedule"
	case MessageAttachmentSchedule:
		return "Message Attachment Schedule"
	case ChannelSchedule:
		return "Channel Schedule"
	case ChannelScheduleNSFW:
		return "Channel Schedule NSFW"
	}

	return "Undefined"
}

func (s Schedule) Description() string {
	switch s {
	case NoSchedule:
		return "Schedules nothing."
	case MessageSchedule:
		return "Schedules every message sent to be deleted after a specified amount time."
	case MessageAttachmentSchedule:
		return "Schedules every message sent with any kind of attachment to be deleted after a specified amount time."
	case ChannelSchedule:
		return "Schedules channel to be deleted and replaced with the exact copy after a specified amount of time."
	case ChannelScheduleNSFW:
		return "Schedules channel to be made age restricted and deleted and replaced with the exact copy after a specified amount of time."
	}

	return "Undefined"
}

type TempChannel struct {
	Schedule Schedule
	Time     time.Time
	Duration time.Duration

	Messages map[string]time.Time
}

func NewTempChannel() *TempChannel {
	return &TempChannel{
		Schedule: NoSchedule,
		Time:     time.Now(),
		Duration: time.Minute,

		Messages: map[string]time.Time{},
	}
}

var TempChannels = map[string]*TempChannel{}
var TempTicker = time.NewTicker(time.Minute) // update every minute
var TempInit = false

// ready controls the temp channels and messages
func ready(s *discordgo.Session) {
	if TempInit {
		return
	}
	TempInit = true

	for {
		<-TempTicker.C

		for cID, channel := range TempChannels {
			now := time.Now()

			switch channel.Schedule {
			case MessageSchedule, MessageAttachmentSchedule:
				messages := []string{}

				for mID, time := range channel.Messages {
					if now.Sub(time) < channel.Duration {
						continue
					}

					if messages = append(messages, mID); len(messages) != 100 {
						continue
					}

					s.ChannelMessagesBulkDelete(cID, messages)
					messages = []string{}
				}

				s.ChannelMessagesBulkDelete(cID, messages)

			case ChannelSchedule, ChannelScheduleNSFW:
				if channel.Time.Add(channel.Duration).Before(now) {
					continue
				}

				c, err := s.State.Channel(cID)
				if err != nil {
					if c, err = s.Channel(cID); err != nil {
						s.ChannelMessageSend(cID, err.Error())
						continue
					}

					s.State.ChannelAdd(c)
				}

				if channel.Schedule == ChannelScheduleNSFW {
					c.NSFW = true

					if _, err := s.ChannelEdit(cID, &discordgo.ChannelEdit{NSFW: &c.NSFW}); err != nil {
						s.ChannelMessageSend(cID, err.Error())
						continue
					}
				}

				if _, err = s.ChannelDelete(cID); err != nil {
					s.ChannelMessageSend(cID, err.Error())
					continue
				}

				new, err := s.GuildChannelCreateComplex(c.GuildID, discordgo.GuildChannelCreateData{
					Name:                 c.Name,
					Type:                 c.Type,
					Topic:                c.Topic,
					Bitrate:              c.Bitrate,
					UserLimit:            c.UserLimit,
					RateLimitPerUser:     c.RateLimitPerUser,
					Position:             c.Position,
					PermissionOverwrites: c.PermissionOverwrites,
					ParentID:             c.ParentID,
					NSFW:                 c.NSFW,
				})
				if err != nil {
					s.ChannelMessageSend(cID, err.Error())
					continue
				}

				TempChannels[new.ID] = channel
				channel.Messages = map[string]time.Time{}
				channel.Time = time.Now()
			}
		}
	}
}

func OnTemp(s *discordgo.Session, i interface{}) {
	switch e := i.(type) {
	case *discordgo.Ready:
		ready(s)

	case *discordgo.MessageCreate:
		if channel, ok := TempChannels[e.ChannelID]; ok {
			switch channel.Schedule {
			case MessageAttachmentSchedule:
				if len(e.Attachments) == 0 {
					return
				}

				fallthrough
			case MessageSchedule:
				channel.Messages[e.ID] = time.Now()
			}
		}

	case *discordgo.MessageDelete:
		if channel, ok := TempChannels[e.ChannelID]; ok {
			delete(channel.Messages, e.ID)
		}

	case *discordgo.ChannelDelete:
		delete(TempChannels, e.ID)
	}
}
