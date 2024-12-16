package commander

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/spf13/pflag"
)

type Manage struct {
	Stdout, Stderr io.Writer
	// if true, the help command will be added to the command list
	// Help is used for all commands.
	// default is true
	Help *bool
}

func NewCommandManager() *Manage {
	return &Manage{}
}

func (m *Manage) WithStdout(w io.Writer) *Manage {
	m.Stdout = w
	return m
}

func (m *Manage) WithStderr(w io.Writer) *Manage {
	m.Stderr = w
	return m
}

func (m *Manage) WithHelp(b bool) *Manage {
	m.Help = &b
	return m
}

func (m *Manage) Build() *CommandManager {
	mgr := &CommandManager{
		stdout: os.Stdout,
		stderr: os.Stderr,
		help:   true,
	}
	if m.Stdout != nil {
		mgr.stdout = m.Stdout
	}
	if m.Stderr != nil {
		mgr.stderr = m.Stderr
	}
	if m.Help != nil {
		mgr.help = *m.Help
	}
	return mgr
}

type Commander interface {
	Name() string
	Short() string
	Long() string
	SetFlags(f *pflag.FlagSet)
	Run(ctx context.Context) error
}

type CommandManager struct {
	stdout   io.Writer
	stderr   io.Writer
	help     bool
	commands map[string]Command
}

func (c *CommandManager) Register(cmd Commander) {
	if c.commands == nil {
		c.commands = make(map[string]Command)
	}
	c.commands[cmd.Name()] = Command{Commander: cmd}
}

var (
	ErrNoCommand             = errors.New("No command provided")
	ErrCommandNotImplemented = errors.New("Command not implemented")
	ErrDisableHelp           = errors.New("Help command is disabled")
)

func (c *CommandManager) Run(ctx context.Context) error {
	args := os.Args[1:]

	if len(args) == 0 {
		return ErrNoCommand
	}

	if args[0] == "help" {
		if !c.help {
			return ErrDisableHelp
		}
		if len(args) == 1 {
			fmt.Fprintf(c.stderr, "No command provided. Use 'help <command>' for more information about a command.\n")
			return nil
		}
		c.printHelp(args[1:])
		return nil
	}

	cmd, ok := c.commands[args[0]]
	if !ok {
		return ErrCommandNotImplemented
	}

	pos, sub := cmd.search(args[1:])
	if sub != nil {
		cmd = *sub
		pos += 1
	}

	f := pflag.NewFlagSet(cmd.Name(), pflag.ExitOnError)
	cmd.SetFlags(f)

	if err := f.Parse(args[pos:]); err != nil {
		return err
	}

	if err := cmd.Run(ctx); err != nil {
		return err
	}

	return nil
}

const usage = `
Usage: %s <command> [arguments]

Subcommands:
%s
Flags:
%s
Use "%s help <command>" for more information about a command.
`

func (c *CommandManager) printHelp(args []string) {
	cmd, ok := c.commands[args[0]]
	if !ok {
		return
	}
	if _, sub := cmd.search(args[1:]); sub != nil {
		cmd = *sub
	}
	_c, ok := cmd.Commander.(Command)
	if ok {
	}
	fmt.Fprintf(c.stdout, fmt.Sprintf(usage, os.Args[0], _c.subcommands(), cmd.flags(), os.Args[0]))
}

type Command struct {
	Commander
	help        string
	subCommands map[string]Command
}

func NewCommand(cmd Commander) Command {
	return Command{Commander: cmd}
}

func (c *Command) Register(cmd Commander) {
	if c.subCommands == nil {
		c.subCommands = make(map[string]Command)
	}
	c.subCommands[cmd.Name()] = Command{Commander: cmd}
}

func (c *Command) SetHelp(help string) {
	c.help = help
}

func (c *Command) flags() string {
	f := pflag.NewFlagSet(c.Name(), pflag.ExitOnError)
	c.SetFlags(f)
	buf := bytes.NewBuffer(nil)
	f.SetOutput(buf)
	f.PrintDefaults()
	return buf.String()
}

func (c *Command) subcommands() string {
	subcmds := "    "
	sub := c.subCommands
	keys := make([]string, 0, len(sub))
	for k := range sub {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		subcmds += fmt.Sprintf("%s\n    ", sub[k].Name())
	}
	return subcmds
}

func (c *Command) search(args []string) (int, *Command) {
	var pos int
	cmd := c
	for _, arg := range args {
		c, ok := cmd.Commander.(Command)
		if !ok {
			break
		}
		subc, ok := c.subCommands[arg]
		if !ok {
			break
		}
		cmd = &subc
		pos++
	}
	return pos, cmd
}
