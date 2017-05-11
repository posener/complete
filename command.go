package complete

import "github.com/posener/complete/match"

// Command represents a command line
// It holds the data that enables auto completion of command line
// Command can also be a sub command.
type Command struct {
	// Sub is map of sub commands of the current command
	// The key refer to the sub command name, and the value is it's
	// Command descriptive struct.
	Sub Commands

	// Flags is a map of flags that the command accepts.
	// The key is the flag name, and the value is it's prediction predict.
	Flags Flags

	// Args are extra arguments that the command accepts, those who are
	// given without any flag before.
	Args Predicate
}

// Commands is the type of Sub member, it maps a command name to a command struct
type Commands map[string]Command

// Flags is the type Flags of the Flags member, it maps a flag name to the flag
// prediction predict.
type Flags map[string]Predicate

// predict returns all available complete predict for the given command
// all are all except the last command line arguments relevant to the command
func (c *Command) predict(a args) (options []match.Matcher, only bool) {

	// if wordCompleted has something that needs to follow it,
	// it is the most relevant completion
	if predicate, ok := c.Flags[a.lastCompleted]; ok && predicate != nil {
		Log("Predicting according to flag %s", a.beingTyped)
		return predicate.predict(a.beingTyped), true
	}

	sub, options, only := c.searchSub(a)
	if only {
		return
	}

	// if no sub command was found, return a list of the sub commands
	if sub == "" {
		options = append(options, c.subCommands()...)
	}

	// add global available complete predict
	for flag := range c.Flags {
		options = append(options, match.Prefix(flag))
	}

	// add additional expected argument of the command
	options = append(options, c.Args.predict(a.beingTyped)...)

	return
}

// searchSub searches recursively within sub commands if the sub command appear
// in the on of the arguments.
func (c *Command) searchSub(a args) (sub string, all []match.Matcher, only bool) {
	for i, arg := range a.completed {
		if cmd, ok := c.Sub[arg]; ok {
			sub = arg
			all, only = cmd.predict(a.from(i))
			return
		}
	}
	return
}

// suvCommands returns a list of matchers according to the sub command names
func (c *Command) subCommands() []match.Matcher {
	subs := make([]match.Matcher, 0, len(c.Sub))
	for sub := range c.Sub {
		subs = append(subs, match.Prefix(sub))
	}
	return subs
}
