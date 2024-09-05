package handler

import (
	"eco/eco"
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(s *eco.Session, e interface{}) {
	switch r := e.(type) {
	case *discordgo.Ready:
		log.Printf("Logged in as %s#%s", r.User.Username, r.User.Discriminator)
	}
}
