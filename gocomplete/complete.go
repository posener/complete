package main

import (
	"github.com/posener/complete"
)

var (
	build = complete.Command{
		Flags: complete.Flags{
			"-o": complete.FlagUnknownFollow,
			"-i": complete.FlagNoFollow,
		},
	}

	test = complete.Command{
		Flags: complete.Flags{
			"-run":   complete.FlagUnknownFollow,
			"-count": complete.FlagUnknownFollow,
		},
	}

	gogo = complete.Command{
		Sub: complete.Commands{
			"build": build,
			"test":  test,
		},
		Flags: complete.Flags{
			"-h": complete.FlagNoFollow,
		},
	}
)

func main() {
	complete.New(gogo).Complete()
}
