// Package main is complete tool for the go command line
package main

import (
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

var (
	ellipsis   = predict.Set{"./..."}
	anyPackage = complete.PredictFunc(predictPackages)
	goFiles    = predict.Files("*.go")
	anyFile    = predict.Files("*")
	anyGo      = predict.Or(goFiles, anyPackage, ellipsis)
)

func main() {
	build := &complete.Command{
		Flags: map[string]complete.Predictor{
			"o": anyFile,
			"i": predict.Nothing,

			"a":             predict.Nothing,
			"n":             predict.Nothing,
			"p":             predict.Something,
			"race":          predict.Nothing,
			"msan":          predict.Nothing,
			"v":             predict.Nothing,
			"work":          predict.Nothing,
			"x":             predict.Nothing,
			"asmflags":      predict.Something,
			"buildmode":     predict.Something,
			"compiler":      predict.Something,
			"gccgoflags":    predict.Set{"gccgo", "gc"},
			"gcflags":       predict.Something,
			"installsuffix": predict.Something,
			"ldflags":       predict.Something,
			"linkshared":    predict.Nothing,
			"pkgdir":        anyPackage,
			"tags":          predict.Something,
			"toolexec":      predict.Something,
		},
		Args: anyGo,
	}

	run := &complete.Command{
		Flags: map[string]complete.Predictor{
			"exec": predict.Something,
		},
		Args: goFiles,
	}

	test := &complete.Command{
		Flags: map[string]complete.Predictor{
			"args": predict.Something,
			"c":    predict.Nothing,
			"exec": predict.Something,

			"bench":     predictBenchmark,
			"benchtime": predict.Something,
			"count":     predict.Something,
			"cover":     predict.Nothing,
			"covermode": predict.Set{"set", "count", "atomic"},
			"coverpkg":  predict.Dirs("*"),
			"cpu":       predict.Something,
			"run":       predictTest,
			"short":     predict.Nothing,
			"timeout":   predict.Something,

			"benchmem":             predict.Nothing,
			"blockprofile":         predict.Files("*.out"),
			"blockprofilerate":     predict.Something,
			"coverprofile":         predict.Files("*.out"),
			"cpuprofile":           predict.Files("*.out"),
			"memprofile":           predict.Files("*.out"),
			"memprofilerate":       predict.Something,
			"mutexprofile":         predict.Files("*.out"),
			"mutexprofilefraction": predict.Something,
			"outputdir":            predict.Dirs("*"),
			"trace":                predict.Files("*.out"),
		},
		Args: anyGo,
	}

	fmt := &complete.Command{
		Flags: map[string]complete.Predictor{
			"n": predict.Nothing,
			"x": predict.Nothing,
		},
		Args: anyGo,
	}

	get := &complete.Command{
		Flags: map[string]complete.Predictor{
			"d":        predict.Nothing,
			"f":        predict.Nothing,
			"fix":      predict.Nothing,
			"insecure": predict.Nothing,
			"t":        predict.Nothing,
			"u":        predict.Nothing,
		},
		Args: anyGo,
	}

	generate := &complete.Command{
		Flags: map[string]complete.Predictor{
			"n":   predict.Nothing,
			"x":   predict.Nothing,
			"v":   predict.Nothing,
			"run": predict.Something,
		},
		Args: anyGo,
	}

	vet := &complete.Command{
		Flags: map[string]complete.Predictor{
			"n": predict.Nothing,
			"x": predict.Nothing,
		},
		Args: anyGo,
	}

	list := &complete.Command{
		Flags: map[string]complete.Predictor{
			"e":    predict.Nothing,
			"f":    predict.Something,
			"json": predict.Nothing,
		},
		Args: predict.Or(anyPackage, ellipsis),
	}

	doc := &complete.Command{
		Flags: map[string]complete.Predictor{
			"c":   predict.Nothing,
			"cmd": predict.Nothing,
			"u":   predict.Nothing,
		},
		Args: anyPackage,
	}

	tool := &complete.Command{
		Flags: map[string]complete.Predictor{
			"n": predict.Nothing,
		},
		Sub: map[string]*complete.Command{
			"addr2line": {
				Args: anyFile,
			},
			"asm": {
				Flags: map[string]complete.Predictor{
					"D":        predict.Something,
					"I":        predict.Dirs("*"),
					"S":        predict.Nothing,
					"V":        predict.Nothing,
					"debug":    predict.Nothing,
					"dynlink":  predict.Nothing,
					"e":        predict.Nothing,
					"o":        anyFile,
					"shared":   predict.Nothing,
					"trimpath": predict.Nothing,
				},
				Args: predict.Files("*.s"),
			},
			"cgo": {
				Flags: map[string]complete.Predictor{
					"debug-define":       predict.Nothing,
					"debug-gcc":          predict.Nothing,
					"dynimport":          anyFile,
					"dynlinker":          predict.Nothing,
					"dynout":             anyFile,
					"dynpackage":         anyPackage,
					"exportheader":       predict.Dirs("*"),
					"gccgo":              predict.Nothing,
					"gccgopkgpath":       predict.Dirs("*"),
					"gccgoprefix":        predict.Something,
					"godefs":             predict.Nothing,
					"import_runtime_cgo": predict.Nothing,
					"import_syscall":     predict.Nothing,
					"importpath":         predict.Dirs("*"),
					"objdir":             predict.Dirs("*"),
					"srcdir":             predict.Dirs("*"),
				},
				Args: goFiles,
			},
			"compile": {
				Flags: map[string]complete.Predictor{
					"%":              predict.Nothing,
					"+":              predict.Nothing,
					"B":              predict.Nothing,
					"D":              predict.Dirs("*"),
					"E":              predict.Nothing,
					"I":              predict.Dirs("*"),
					"K":              predict.Nothing,
					"N":              predict.Nothing,
					"S":              predict.Nothing,
					"V":              predict.Nothing,
					"W":              predict.Nothing,
					"asmhdr":         anyFile,
					"bench":          anyFile,
					"buildid":        predict.Nothing,
					"complete":       predict.Nothing,
					"cpuprofile":     anyFile,
					"d":              predict.Nothing,
					"dynlink":        predict.Nothing,
					"e":              predict.Nothing,
					"f":              predict.Nothing,
					"h":              predict.Nothing,
					"i":              predict.Nothing,
					"importmap":      predict.Something,
					"installsuffix":  predict.Something,
					"j":              predict.Nothing,
					"l":              predict.Nothing,
					"largemodel":     predict.Nothing,
					"linkobj":        anyFile,
					"live":           predict.Nothing,
					"m":              predict.Nothing,
					"memprofile":     predict.Nothing,
					"memprofilerate": predict.Something,
					"msan":           predict.Nothing,
					"nolocalimports": predict.Nothing,
					"o":              anyFile,
					"p":              predict.Dirs("*"),
					"pack":           predict.Nothing,
					"r":              predict.Nothing,
					"race":           predict.Nothing,
					"s":              predict.Nothing,
					"shared":         predict.Nothing,
					"traceprofile":   anyFile,
					"trimpath":       predict.Something,
					"u":              predict.Nothing,
					"v":              predict.Nothing,
					"w":              predict.Nothing,
					"wb":             predict.Nothing,
				},
				Args: goFiles,
			},
			"cover": {
				Flags: map[string]complete.Predictor{
					"func": predict.Something,
					"html": predict.Something,
					"mode": predict.Set{"set", "count", "atomic"},
					"o":    anyFile,
					"var":  predict.Something,
				},
				Args: anyFile,
			},
			"dist": {
				Sub: map[string]*complete.Command{
					"banner":    {Flags: map[string]complete.Predictor{"v": predict.Nothing}},
					"bootstrap": {Flags: map[string]complete.Predictor{"v": predict.Nothing}},
					"clean":     {Flags: map[string]complete.Predictor{"v": predict.Nothing}},
					"env":       {Flags: map[string]complete.Predictor{"v": predict.Nothing, "p": predict.Nothing}},
					"install":   {Flags: map[string]complete.Predictor{"v": predict.Nothing}, Args: predict.Dirs("*")},
					"list":      {Flags: map[string]complete.Predictor{"v": predict.Nothing, "json": predict.Nothing}},
					"test":      {Flags: map[string]complete.Predictor{"v": predict.Nothing, "h": predict.Nothing}},
					"version":   {Flags: map[string]complete.Predictor{"v": predict.Nothing}},
				},
			},
			"doc": doc,
			"fix": {
				Flags: map[string]complete.Predictor{
					"diff":  predict.Nothing,
					"force": predict.Something,
					"r":     predict.Set{"context", "gotypes", "netipv6zone", "printerconfig"},
				},
				Args: anyGo,
			},
			"link": {
				Flags: map[string]complete.Predictor{
					"B":              predict.Something, // note
					"D":              predict.Something, // address (default -1)
					"E":              predict.Something, // entry symbol name
					"H":              predict.Something, // header type
					"I":              predict.Something, // linker binary
					"L":              predict.Dirs("*"), // directory
					"R":              predict.Something, // quantum (default -1)
					"T":              predict.Something, // address (default -1)
					"V":              predict.Nothing,
					"X":              predict.Something,
					"a":              predict.Something,
					"buildid":        predict.Something, // build id
					"buildmode":      predict.Something,
					"c":              predict.Nothing,
					"cpuprofile":     anyFile,
					"d":              predict.Nothing,
					"debugtramp":     predict.Something, // int
					"dumpdep":        predict.Nothing,
					"extar":          predict.Something,
					"extld":          predict.Something,
					"extldflags":     predict.Something, // flags
					"f":              predict.Nothing,
					"g":              predict.Nothing,
					"importcfg":      anyFile,
					"installsuffix":  predict.Something, // dir suffix
					"k":              predict.Something, // symbol
					"libgcc":         predict.Something, // maybe "none"
					"linkmode":       predict.Something, // mode
					"linkshared":     predict.Nothing,
					"memprofile":     anyFile,
					"memprofilerate": predict.Something, // rate
					"msan":           predict.Nothing,
					"n":              predict.Nothing,
					"o":              predict.Something,
					"pluginpath":     predict.Something,
					"r":              predict.Something, // "dir1:dir2:..."
					"race":           predict.Nothing,
					"s":              predict.Nothing,
					"tmpdir":         predict.Dirs("*"),
					"u":              predict.Nothing,
					"v":              predict.Nothing,
					"w":              predict.Nothing,
					// "h":           predict.Something, // halt on error
				},
				Args: predict.Or(
					predict.Files("*.a"),
					predict.Files("*.o"),
				),
			},
			"nm": {
				Flags: map[string]complete.Predictor{
					"n":    predict.Nothing,
					"size": predict.Nothing,
					"sort": predict.Something,
					"type": predict.Nothing,
				},
				Args: anyGo,
			},
			"objdump": {
				Flags: map[string]complete.Predictor{
					"s": predict.Something,
					"S": predict.Nothing,
				},
				Args: anyFile,
			},
			"pack": {
				/* this lacks the positional aspect of all these params */
				Flags: map[string]complete.Predictor{
					"c":  predict.Nothing,
					"p":  predict.Nothing,
					"r":  predict.Nothing,
					"t":  predict.Nothing,
					"x":  predict.Nothing,
					"cv": predict.Nothing,
					"pv": predict.Nothing,
					"rv": predict.Nothing,
					"tv": predict.Nothing,
					"xv": predict.Nothing,
				},
				Args: predict.Or(
					predict.Files("*.a"),
					predict.Files("*.o"),
				),
			},
			"pprof": {
				Flags: map[string]complete.Predictor{
					"callgrind":     predict.Nothing,
					"disasm":        predict.Something,
					"dot":           predict.Nothing,
					"eog":           predict.Nothing,
					"evince":        predict.Nothing,
					"gif":           predict.Nothing,
					"gv":            predict.Nothing,
					"list":          predict.Something,
					"pdf":           predict.Nothing,
					"peek":          predict.Something,
					"png":           predict.Nothing,
					"proto":         predict.Nothing,
					"ps":            predict.Nothing,
					"raw":           predict.Nothing,
					"svg":           predict.Nothing,
					"tags":          predict.Nothing,
					"text":          predict.Nothing,
					"top":           predict.Nothing,
					"tree":          predict.Nothing,
					"web":           predict.Nothing,
					"weblist":       predict.Something,
					"output":        anyFile,
					"functions":     predict.Nothing,
					"files":         predict.Nothing,
					"lines":         predict.Nothing,
					"addresses":     predict.Nothing,
					"base":          predict.Something,
					"drop_negative": predict.Nothing,
					"cum":           predict.Nothing,
					"seconds":       predict.Something,
					"nodecount":     predict.Something,
					"nodefraction":  predict.Something,
					"edgefraction":  predict.Something,
					"sample_index":  predict.Nothing,
					"mean":          predict.Nothing,
					"inuse_space":   predict.Nothing,
					"inuse_objects": predict.Nothing,
					"alloc_space":   predict.Nothing,
					"alloc_objects": predict.Nothing,
					"total_delay":   predict.Nothing,
					"contentions":   predict.Nothing,
					"mean_delay":    predict.Nothing,
					"runtime":       predict.Nothing,
					"focus":         predict.Something,
					"ignore":        predict.Something,
					"tagfocus":      predict.Something,
					"tagignore":     predict.Something,
					"call_tree":     predict.Nothing,
					"unit":          predict.Something,
					"divide_by":     predict.Something,
					"buildid":       predict.Something,
					"tools":         predict.Dirs("*"),
					"help":          predict.Nothing,
				},
				Args: anyFile,
			},
			"tour": {
				Flags: map[string]complete.Predictor{
					"http":        predict.Something,
					"openbrowser": predict.Nothing,
				},
			},
			"trace": {
				Flags: map[string]complete.Predictor{
					"http":  predict.Something,
					"pprof": predict.Set{"net", "sync", "syscall", "sched"},
				},
				Args: anyFile,
			},
			"vet": {
				Flags: map[string]complete.Predictor{
					"all":                 predict.Nothing,
					"asmdecl":             predict.Nothing,
					"assign":              predict.Nothing,
					"atomic":              predict.Nothing,
					"bool":                predict.Nothing,
					"buildtags":           predict.Nothing,
					"cgocall":             predict.Nothing,
					"composites":          predict.Nothing,
					"compositewhitelist":  predict.Nothing,
					"copylocks":           predict.Nothing,
					"httpresponse":        predict.Nothing,
					"lostcancel":          predict.Nothing,
					"methods":             predict.Nothing,
					"nilfunc":             predict.Nothing,
					"printf":              predict.Nothing,
					"printfuncs":          predict.Something,
					"rangeloops":          predict.Nothing,
					"shadow":              predict.Nothing,
					"shadowstrict":        predict.Nothing,
					"shift":               predict.Nothing,
					"structtags":          predict.Nothing,
					"tags":                predict.Something,
					"tests":               predict.Nothing,
					"unreachable":         predict.Nothing,
					"unsafeptr":           predict.Nothing,
					"unusedfuncs":         predict.Something,
					"unusedresult":        predict.Nothing,
					"unusedstringmethods": predict.Something,
					"v":                   predict.Nothing,
				},
				Args: anyGo,
			},
		},
	}

	clean := &complete.Command{
		Flags: map[string]complete.Predictor{
			"i":         predict.Nothing,
			"r":         predict.Nothing,
			"n":         predict.Nothing,
			"x":         predict.Nothing,
			"cache":     predict.Nothing,
			"testcache": predict.Nothing,
			"modcache":  predict.Nothing,
		},
		Args: predict.Or(anyPackage, ellipsis),
	}

	env := &complete.Command{
		Args: predict.Something,
	}

	bug := &complete.Command{}
	version := &complete.Command{}

	fix := &complete.Command{
		Args: anyGo,
	}

	modDownload := &complete.Command{
		Flags: map[string]complete.Predictor{
			"json": predict.Nothing,
		},
		Args: anyPackage,
	}

	modEdit := &complete.Command{
		Flags: map[string]complete.Predictor{
			"fmt":    predict.Nothing,
			"module": predict.Nothing,
			"print":  predict.Nothing,

			"exclude":     anyPackage,
			"dropexclude": anyPackage,
			"replace":     anyPackage,
			"dropreplace": anyPackage,
			"require":     anyPackage,
			"droprequire": anyPackage,
		},
		Args: predict.Files("go.mod"),
	}

	modGraph := &complete.Command{}

	modInit := &complete.Command{
		Args: predict.Something,
	}

	modTidy := &complete.Command{
		Flags: map[string]complete.Predictor{
			"v": predict.Nothing,
		},
	}

	modVendor := &complete.Command{
		Flags: map[string]complete.Predictor{
			"v": predict.Nothing,
		},
	}

	modVerify := &complete.Command{}

	modWhy := &complete.Command{
		Flags: map[string]complete.Predictor{
			"m":      predict.Nothing,
			"vendor": predict.Nothing,
		},
		Args: anyPackage,
	}

	modHelp := &complete.Command{
		Sub: map[string]*complete.Command{
			"download": &complete.Command{},
			"edit":     &complete.Command{},
			"graph":    &complete.Command{},
			"init":     &complete.Command{},
			"tidy":     &complete.Command{},
			"vendor":   &complete.Command{},
			"verify":   &complete.Command{},
			"why":      &complete.Command{},
		},
	}

	mod := &complete.Command{
		Sub: map[string]*complete.Command{
			"download": modDownload,
			"edit":     modEdit,
			"graph":    modGraph,
			"init":     modInit,
			"tidy":     modTidy,
			"vendor":   modVendor,
			"verify":   modVerify,
			"why":      modWhy,
			"help":     modHelp,
		},
	}

	help := &complete.Command{
		Sub: map[string]*complete.Command{
			"bug":         &complete.Command{},
			"build":       &complete.Command{},
			"clean":       &complete.Command{},
			"doc":         &complete.Command{},
			"env":         &complete.Command{},
			"fix":         &complete.Command{},
			"fmt":         &complete.Command{},
			"generate":    &complete.Command{},
			"get":         &complete.Command{},
			"install":     &complete.Command{},
			"list":        &complete.Command{},
			"mod":         modHelp,
			"run":         &complete.Command{},
			"test":        &complete.Command{},
			"tool":        &complete.Command{},
			"version":     &complete.Command{},
			"vet":         &complete.Command{},
			"buildmode":   &complete.Command{},
			"c":           &complete.Command{},
			"cache":       &complete.Command{},
			"environment": &complete.Command{},
			"filetype":    &complete.Command{},
			"go.mod":      &complete.Command{},
			"gopath":      &complete.Command{},
			"gopath-get":  &complete.Command{},
			"goproxy":     &complete.Command{},
			"importpath":  &complete.Command{},
			"modules":     &complete.Command{},
			"module-get":  &complete.Command{},
			"packages":    &complete.Command{},
			"testflag":    &complete.Command{},
			"testfunc":    &complete.Command{},
		},
	}

	// commands that also accepts the build flags
	for name, options := range build.Flags {
		test.Flags[name] = options
		run.Flags[name] = options
		list.Flags[name] = options
		vet.Flags[name] = options
		get.Flags[name] = options
	}

	gogo := &complete.Command{
		Sub: map[string]*complete.Command{
			"build":    build,
			"install":  build, // install and build have the same flags
			"run":      run,
			"test":     test,
			"fmt":      fmt,
			"get":      get,
			"generate": generate,
			"vet":      vet,
			"list":     list,
			"doc":      doc,
			"tool":     tool,
			"clean":    clean,
			"env":      env,
			"bug":      bug,
			"fix":      fix,
			"version":  version,
			"mod":      mod,
			"help":     help,
		},
		Flags: map[string]complete.Predictor{
			"h": predict.Nothing,
		},
	}

	gogo.Complete("go")
}
