package complete

type args struct {
	all           []string
	completed     []string
	beingTyped    string
	lastCompleted string
}

func newArgs(line []string) args {
	completed := removeLast(line)
	return args{
		all:           line[1:],
		completed:     completed,
		beingTyped:    last(line),
		lastCompleted: last(completed),
	}
}

func (a args) from(i int) args {
	a.all = a.all[i:]
	a.completed = a.completed[i:]
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
