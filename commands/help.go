package commands

import (
	"github.com/stavratum/eco/handler/router"
)

var TestArg = router.NewArgument(router.String, []string{"testarg"})

var Help = &router.Command{
	Name:        "help",
	Usage:       "help [command]",
	Description: "Returns information about a command and it's features.",

	Aliases: []string{"help", "h"},
	Arguments: router.Arguments([]*router.Argument{
		TestArg,
	}),

	Handler: func(ctx *router.Context) {
		if v, ok := ctx.Arguments[TestArg]; ok {
			ctx.SendReply("testarg: " + v.(string))
			return
		}

		ctx.SendReply("helping... 10% 20% 30% Blah blah blah blah system error")
	},
}
