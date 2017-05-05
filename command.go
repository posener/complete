package complete

type Command struct {
	Sub   map[string]Command
	Flags map[string]FlagOptions
}

// options returns all available complete options for the given command
// args are all except the last command line arguments relevant to the command
func (c *Command) options(args []string) (options []string, only bool) {

	// remove the first argument, which is the command name
	args = args[1:]

	// if prev has something that needs to follow it,
	// it is the most relevant completion
	if options, ok := c.Flags[last(args)]; ok && options.HasFollow {
		return options.FollowsOptions, true
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
		options = append(options, flag)
	}

	return
}

func (c *Command) searchSub(args []string) (sub string, all []string, only bool) {
	for i, arg := range args {
		if cmd, ok := c.Sub[arg]; ok {
			sub = arg
			all, only = cmd.options(args[i:])
			return
		}
	}
	return "", nil, false
}

func (c *Command) subCommands() []string {
	subs := make([]string, 0, len(c.Sub))
	for sub := range c.Sub {
		subs = append(subs, sub)
	}
	return subs
}

