package main

import (
	"fmt"
	"strings"

	avisha "github.com/jackmordaunt/avisha-fn"
)

// Command executes the associated action and that of any matching children.
type Command struct {
	Name     string
	Action   func(app *avisha.App, args []string) error
	Children []Command
}

// Handle executes the command action and invokes any children that match the
// next cmd string.
func (cmd Command) Handle(app *avisha.App, args []string) error {
	if cmd.Action != nil {
		if err := cmd.Action(app, args); err != nil {
			return fmt.Errorf("%s: %w", cmd.Name, err)
		}
	}

	if len(args) < 1 {
		return nil
	}

	for _, child := range cmd.Children {
		if child.Name == args[0] {
			if err := child.Handle(app, args[1:]); err != nil {
				return Err{Ctx: cmd.Name, Err: err}
			}
			return nil
		}
	}

	return nil
}

// ArgMap associates argument handlers with argument names.
// Allows for `name=value` syntax.
type ArgMap struct {
	Handlers map[string]func(value string)
}

// Match performs the argument matching.
func (m ArgMap) Match(args []string) {
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		name, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if handler, ok := m.Handlers[name]; ok {
			handler(value)
		}
	}
}

// Err wraps and error and adds context.
// Ignores empty context.
type Err struct {
	Ctx string
	Err error
}

func (f Err) Error() string {
	if len(f.Ctx) > 0 {
		return fmt.Sprintf("%s: %s", f.Ctx, f.Err)
	}
	return fmt.Sprintf("%s", f.Err)
}

func (f Err) Unwrap() error {
	return f.Err
}
