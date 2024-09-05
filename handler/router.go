package handler

import (
	"eco/eco"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Usage       string
	Description string

	Handler func(*eco.Session, *discordgo.MessageCreate)
}

type Router struct {
	Prefix string

	Commands      map[string]*Command
	CommandsCount uint16
}

func (r *Router) AddCommand(c *Command) error {
	if r.Commands[c.Name] == nil {
		r.Commands[c.Name] = c
		r.CommandsCount += 1

		return nil
	}

	return fmt.Errorf("command with name %s already exists", c.Name)
}

func (r *Router) RemoveCommand(name string) {
	if r.Commands[name] == nil {
		return
	}

	delete(r.Commands, name)
	r.CommandsCount -= 1
}

func (r *Router) Handler(s *eco.Session, i interface{}) {
	switch e := i.(type) {
	case *discordgo.Ready:
		log.Printf("Router: Loaded %d commands.", r.CommandsCount)
	case *discordgo.MessageCreate:
		if e.Author.Bot {
			return
		}

		content, found := strings.CutPrefix(e.Content, r.Prefix)
		if !found {
			return
		}

		arguments := strings.Split(content, " ")
		name := strings.ToLower(arguments[0])

		if command, ok := r.Commands[name]; ok {
			command.Handler(s, e)
		}
	}
}

func NewRouter() *Router {
	return &Router{
		Prefix:   "!",
		Commands: map[string]*Command{},
	}
}
