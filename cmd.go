package commander

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
)

type Manage struct {
	Stdout io.Writer
	Stderr io.Writer
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

func (m *Manage) Build() *CommandManager {
	mgr := &CommandManager{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	if m.Stdout != nil {
		mgr.stdout = m.Stdout
	}
	if m.Stderr != nil {
		mgr.stderr = m.Stderr
	}
	return mgr
}

type Commander interface {
	Name() string
	Short() string
	Long() string
	SetFlags(f *flag.FlagSet)
	Run(ctx context.Context) error
}

type CommandManager struct {
	stdout   io.Writer
	stderr   io.Writer
	commands map[string]Command
}

func (c *CommandManager) Register(cmd Commander) {
	if c.commands == nil {
		c.commands = make(map[string]Command)
	}
	c.commands[cmd.Name()] = Command{Commander: cmd}
}

type Command struct {
	Commander
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

var ErrNoCommand = errors.New("No command provided")
var ErrCommandNotImplemented = errors.New("Command not implemented")

func (c *CommandManager) Run(ctx context.Context) error {
	args := os.Args[1:]

	if len(args) == 0 {
		return ErrNoCommand
	}

	cmd, ok := c.commands[args[0]]
	if !ok {
		return ErrCommandNotImplemented
	}

	parsePos := 1
	for _, arg := range args[1:] {
		c, ok := cmd.Commander.(Command)
		if !ok {
			break
		}
		sub, ok := c.subCommands[arg]
		if !ok {
			break
		}
		cmd = sub
		parsePos++
	}

	f := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.SetFlags(f)

	if err := f.Parse(args[parsePos:]); err != nil {
		return err
	}

	if err := cmd.Run(ctx); err != nil {
		return err
	}

	return nil
}
