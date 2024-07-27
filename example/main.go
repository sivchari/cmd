package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/sivchari/commander"
)

func main() {
	cmd := commander.NewCommandManager().Build()
	print := commander.NewCommand(&PrintCommand{})
	print.Register(&SubCommand{})
	cmd.Register(print)
	cmd.Run(context.Background())
}

var _ commander.Commander = &PrintCommand{}

type PrintCommand struct{}

var name string

func (p *PrintCommand) Run(ctx context.Context) error {
	fmt.Printf("Hello, %s!\n", name)
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

func (p *PrintCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&name, "name", "World", "Name to print")
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

func (s *SubCommand) SetFlags(f *flag.FlagSet) {
	f.IntVar(&number, "number", 0, "Number to print")
}
