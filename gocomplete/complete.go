package main

import (
	"github.com/posener/complete"
)

var completer = complete.New(complete.Command{
	Sub: map[string]complete.Command{
		"build": {
			Flags: map[string]complete.FlagOptions {
				"-o": complete.FlagUnknownFollow,
			}, "-i": complete.FlagNoFollow,
		},
		"test": {
			Flags: map[string]complete.FlagOptions{
				"-run":   complete.FlagUnknownFollow,
				"-count": complete.FlagUnknownFollow,
			},
		},
	},
	Flags: map[string]complete.FlagOptions{
		"-h": complete.FlagNoFollow,
	},
})

func main() {
	completer.Complete()
}
