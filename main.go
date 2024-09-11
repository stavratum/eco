package main

import (
	"eco/commands"
	"eco/eco"
	"eco/handler"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Discord Authentication token")

	flag.Parse()
}

func panic(str string) {
	os.Stdout.WriteString(str)

	if !strings.HasSuffix(str, "\n") {
		os.Stdout.WriteString("\n")
	}

	os.Stdout.WriteString("\nPress any key to exit...")
	os.Stdin.Read(make([]byte, 1))
}

func main() {
	if Token == "" {
		os.Stdout.WriteString("Token must be defined within command line arguments.\n\n")

		os.Stdout.WriteString("Usage:\n")
		fmt.Printf("\t%s -t (token)\n", os.Args[0])

		os.Stdout.WriteString("\nPress any key to exit...")
		os.Stdin.Read(make([]byte, 1))
		return
	}

	s, err := eco.New(Token)
	if err != nil {
		panic(err.Error())
		return
	}

	r := handler.NewRouter()
	s.AddHandler(r.Handler)

	r.AddCommand(commands.Help)
	r.AddCommand(commands.New)

	s.AddHandler(handler.OnTemp)
	r.AddCommand(commands.Temp)

	s.AddHandler(handler.OnReady)

	err = s.Open()
	defer s.Close()

	if err != nil {
		log.Println(err)
		return
	}

	<-make(chan struct{})
}
