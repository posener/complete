package main

import (
	"github.com/posener/complete"
)

var (
	build = complete.Command{
		Flags: complete.Flags{
			"-o": complete.PredictFiles("*"),
			"-i": complete.PredictNothing,
		},
	}

	test = complete.Command{
		Flags: complete.Flags{
			"-run":   complete.PredictAnything,
			"-count": complete.PredictAnything,
		},
	}

	gogo = complete.Command{
		Sub: complete.Commands{
			"build": build,
			"test":  test,
		},
		Flags: complete.Flags{
			"-h": complete.PredictNothing,
		},
	}
)

func main() {
	complete.New(gogo).Complete()
}
