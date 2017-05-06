package complete_test

import "github.com/posener/complete"

func main() {

	// create a Command object, that represents the command we want
	// to complete.
	run := complete.Command{

		// Name must be exactly as the binary that we want to complete
		Name: "run",

		// Sub defines a list of sub commands of the program,
		// this is recursive, since every command is of type command also.
		Sub: complete.Commands{

			// add a build sub command
			"build": complete.Command{

				// define flags of the build sub command
				Flags: complete.Flags{
					// build sub command has a flag '-fast', which
					// does not expects anything after it.
					"-fast": complete.PredictNothing,
				},
			},
		},

		// define flags of the 'run' main command
		Flags: complete.Flags{

			// a flag '-h' which does not expects anything after it
			"-h": complete.PredictNothing,

			// a flag -o, which expects a file ending with .out after
			// it, the tab completion will auto complete for files matching
			// the given pattern.
			"-o": complete.PredictFiles("*.out"),
		},
	}

	// run the command completion, as part of the main() function.
	// this triggers the autocompletion when needed.
	complete.Run(run)
}
