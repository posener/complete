package complete

// Args describes command line arguments
type Args struct {
	All           []string
	Completed     []string
	Last          string
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
