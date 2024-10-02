# commander
commander is framework to make simple command. It is suitable for that has plenty of subcmd

## Usage

See [example](./_example)

```go
package main

import (
    "fmt"
    "os"

    "github.com/sivchari/commander"
)

func main() {
    cmd := commander.NewCommandManager().Build()
    print := commander.NewCommand(&PrintCommand{})
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
```
