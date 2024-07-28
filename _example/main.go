package main

import (
	"context"
	"fmt"

	"github.com/sivchari/commander"
	"github.com/spf13/pflag"
)

func main() {
	cmd := commander.NewCommandManager().Build()
	print := commander.NewCommand(&PrintCommand{})
	print.Register(&SubCommand{})
	print.Register(&Sub2Command{})
	cmd.Register(print)
	cmd.Run(context.Background())
}

var _ commander.Commander = &PrintCommand{}

type PrintCommand struct{}

var last *string
var first *string

func (p *PrintCommand) Run(ctx context.Context) error {
	fmt.Printf("Hello, %s %s!\n", *first, *last)
	return nil
}

func (p *PrintCommand) Name() string {
	return "print"
}

func (p *PrintCommand) Short() string {
	return "Prints Hello, World!"
}

func (p *PrintCommand) Long() string {
	return "Prints Hello, World!"
}

func (p *PrintCommand) SetFlags(f *pflag.FlagSet) {
	first = f.StringP("first", "f", "Hello", "Name to print")
	last = f.StringP("last", "l", "World", "Name to print")
}

var _ commander.Commander = &SubCommand{}

type SubCommand struct{}

var number int

func (s *SubCommand) Run(ctx context.Context) error {
	fmt.Printf("Number: %d\n", number)
	return nil
}

func (s *SubCommand) Name() string {
	return "sub"
}

func (s *SubCommand) Short() string {
	return "Subcommand"
}

func (s *SubCommand) Long() string {
	return "Subcommand"
}

func (s *SubCommand) SetFlags(f *pflag.FlagSet) {
	f.IntVar(&number, "number", 0, "Number to print")
}

var _ commander.Commander = &Sub2Command{}

type Sub2Command struct{}

var message string

func (s *Sub2Command) Run(ctx context.Context) error {
	fmt.Printf("Message: %s\n", message)
	return nil
}

func (s *Sub2Command) Name() string {
	return "sub2"
}

func (s *Sub2Command) Short() string {
	return "Subcommand 2"
}

func (s *Sub2Command) Long() string {
	return "Subcommand 2"
}

func (s *Sub2Command) SetFlags(f *pflag.FlagSet) {
	f.StringVar(&message, "message", "Hello, World!", "Message to print")
}
