package complete

// Command represents a command line
// It holds the data that enables auto completion of a given typed command line
// Command can also be a sub command.
type Command struct {
	// Sub is map of sub commands of the current command
	// The key refer to the sub command name, and the value is it's
	// Command descriptive struct.
	Sub Commands

	// Flags is a map of flags that the command accepts.
	// The key is the flag name, and the value is it's prediction options.
	Flags Flags

	// Args are extra arguments that the command accepts, those who are
	// given without any flag before.
	Args Predicate
}

// Commands is the type of Sub member, it maps a command name to a command struct
type Commands map[string]Command

// Flags is the type Flags of the Flags member, it maps a flag name to the flag
// prediction options.
type Flags map[string]Predicate

// options returns all available complete options for the given command
// args are all except the last command line arguments relevant to the command
func (c *Command) options(args []string) (options []Matcher, only bool) {

	// remove the first argument, which is the command name
	args = args[1:]
	last := last(args)
	// if prev has something that needs to follow it,
	// it is the most relevant completion
	if predicate, ok := c.Flags[last]; ok && predicate != nil {
		return predicate.predict(last), true
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
		options = append(options, MatchPrefix(flag))
	}

	// add additional expected argument of the command
	options = append(options, c.Args.predict(last)...)

	return
}

// searchSub searches recursively within sub commands if the sub command appear
// in the on of the arguments.
func (c *Command) searchSub(args []string) (sub string, all []Matcher, only bool) {
	for i, arg := range args {
		if cmd, ok := c.Sub[arg]; ok {
			sub = arg
			all, only = cmd.options(args[i:])
			return
		}
	}
	return "", nil, false
}

// suvCommands returns a list of matchers according to the sub command names
func (c *Command) subCommands() []Matcher {
	subs := make([]Matcher, 0, len(c.Sub))
	for sub := range c.Sub {
		subs = append(subs, MatchPrefix(sub))
	}
	return subs
}
