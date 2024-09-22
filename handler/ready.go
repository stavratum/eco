package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(s *discordgo.Session, e interface{}) {
	switch r := e.(type) {
	case *discordgo.Ready:
		log.Printf("Logged in as %s#%s", r.User.Username, r.User.Discriminator)

	}
}
