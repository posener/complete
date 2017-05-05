package complete

type Commands map[string]Command

type Flags map[string]Predicate

type Command struct {
	Sub   Commands
	Flags Flags
	Args  Predicate
}

// options returns all available complete options for the given command
// args are all except the last command line arguments relevant to the command
func (c *Command) options(args []string) (options []Option, only bool) {

	// remove the first argument, which is the command name
	args = args[1:]

	// if prev has something that needs to follow it,
	// it is the most relevant completion
	if predicate, ok := c.Flags[last(args)]; ok && !predicate.ExpectsNothing {
		return predicate.predict(), true
	}

	sub, options, only := c.searchSub(args)
	if only {
		return
	}

	// if no subcommand was entered in any of the args, add the
	// subcommands as complete options.
	if sub == "" {
		options = append(options, c.subCommands()...)
	}

	// add global available complete options
	for flag := range c.Flags {
		options = append(options, Arg(flag))
	}

	// add additional expected argument of the command
	if !c.Args.ExpectsNothing {
		options = append(options, c.Args.predict()...)
	}

	return
}

func (c *Command) searchSub(args []string) (sub string, all []Option, only bool) {
	for i, arg := range args {
		if cmd, ok := c.Sub[arg]; ok {
			sub = arg
			all, only = cmd.options(args[i:])
			return
		}
	}
	return "", nil, false
}

func (c *Command) subCommands() []Option {
	subs := make([]Option, 0, len(c.Sub))
	for sub := range c.Sub {
		subs = append(subs, Arg(sub))
	}
	return subs
}
