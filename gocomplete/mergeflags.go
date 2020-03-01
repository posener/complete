package main

import "github.com/posener/complete/v2"

type Flags = map[string]complete.Predictor

// mergeFlags returns a new complete.Flags that merges all flgs.
func mergeFlags(flgs ...Flags) Flags {
	var size int
	for _, flg := range flgs {
		size += len(flg)
	}
	f := make(Flags, size)
	for _, flg := range flgs {
		for k, v := range flg {
			f[k] = v
		}
	}
	return f
}
