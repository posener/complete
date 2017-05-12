package complete

// Args describes command line arguments
type Args struct {
	// All lists of all arguments in command line (not including the command itself)
	All           []string
	// Completed lists of all completed arguments in command line,
	// If the last one is still being typed - no space after it,
	// it won't appear in this list of arguments.
	Completed     []string
	// Last argument in command line, the one being typed, if the last
	// character in the command line is a space, this argument will be empty,
	// otherwise this would be the last word.
	Last          string
	// LastCompleted is the last argument that was fully typed.
	// If the last character in the command line is space, this would be the
	// last word, otherwise, it would be the word before that.
	LastCompleted string
}

func newArgs(line []string) Args {
	completed := removeLast(line)
	return Args{
		All:           line[1:],
		Completed:     completed,
		Last:          last(line),
		LastCompleted: last(completed),
	}
}

func (a Args) from(i int) Args {
	a.All = a.All[i:]
	a.Completed = a.Completed[i:]
	return a
}

func removeLast(a []string) []string {
	if len(a) > 0 {
		return a[:len(a)-1]
	}
	return a
}

func last(args []string) (last string) {
	if len(args) > 0 {
		last = args[len(args)-1]
	}
	return
}
