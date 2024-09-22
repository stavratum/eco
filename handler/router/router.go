package router

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Value uint

// Values for arguments.
const (
	String Value = iota
	Flag
	UInt
	Int

	Channel
	User
	Emoji

	Duration
	Timestamp
)

// Argument (idk how to describe what it does pls help)
type Argument struct {
	Type    Value
	Aliases []string
}

func NewArgument(v Value, a []string) *Argument {
	return &Argument{
		Type:    v,
		Aliases: a,
	}
}

func Arguments(args []*Argument) (v map[string]*Argument) {
	v = map[string]*Argument{}

	for _, arg := range args {
		for _, alias := range arg.Aliases {
			v[alias] = arg
		}
	}

	return
}

// Command (idk how to describe what it does pls help)
type Command struct {
	Name        string
	Usage       string
	Description string

	Aliases []string

	Arguments map[string]*Argument
	Handler   func(*Context)
}

// Parsing tasks.
const (
	ParseCmd uint = iota
	ParseBody
	ParseArg
	ParseVal
)

// Router is a handler for message commands (and maybe interaction commands)
type Router struct {
	Prefix string

	Commands      map[string]*Command
	CommandsCount uint16
}

func (r *Router) AddCommand(c *Command) {
	for _, v := range c.Aliases {
		r.Commands[v] = c
	}

	r.CommandsCount += 1
}

func (r *Router) RemoveCommand(c *Command) {
	for _, v := range c.Aliases {
		delete(r.Commands, v)
	}

	r.CommandsCount -= 1
}

func (r *Router) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	content, found := strings.CutPrefix(m.Content, r.Prefix)
	if !found {
		return
	}

	var (
		ctx *Context  = NewContext(s, m)
		cmd *Command  = nil
		arg *Argument = nil
	)

	for _, part := range strings.Split(content, " ") {
		switch ctx.Task {
		case ParseCmd:
			if part == "" {
				continue
			}

			ctx.Command = strings.ToLower(part)

			var ok bool
			if cmd, ok = r.Commands[ctx.Command]; !ok {
				return
			}

			ctx.Task = ParseBody

		case ParseBody:
			key, found := strings.CutPrefix(part, "-")
			if !found {
				ctx.Body = append(ctx.Body, part)
				continue
			}

			var ok bool
			if arg, ok = cmd.Arguments[key]; !ok {
				ctx.Body = append(ctx.Body, part)
				continue
			}

			ctx.Task = ParseVal

		case ParseVal:
			if part == "" {
				continue
			}

			switch arg.Type {
			case String:
				ctx.Arguments[arg] = part
			}

			ctx.Task = ParseArg

			if arg.Type != Flag {
				continue
			}

			ctx.Arguments[arg] = nil
			fallthrough

		case ParseArg:
			if part == "" {
				continue
			}

			key, found := strings.CutPrefix(part, "-")
			if !found {
				continue
			}

			var ok bool
			if arg, ok = cmd.Arguments[key]; !ok {
				continue
			}

			ctx.Task = ParseVal
		}

	}

	if command, ok := r.Commands[ctx.Command]; ok {
		command.Handler(ctx)
	}
}

func (r *Router) Handler(s *discordgo.Session, i interface{}) {
	switch e := i.(type) {
	case *discordgo.Ready:
		log.Printf("Router: Loaded %d commands.", r.CommandsCount)

	case *discordgo.MessageCreate:
		r.messageCreate(s, e)
	}
}

func New() *Router {
	return &Router{
		Prefix:   "!",
		Commands: map[string]*Command{},
	}
}
