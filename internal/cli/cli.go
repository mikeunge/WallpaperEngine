package cli

import (
    "os"
	"fmt"

	"github.com/akamensky/argparse"
)

type AppInfo struct {
    Name string
    Description string
    Version string
    Author string
}

type Args struct {
    ConfigPath string 
    Debug bool
    Verbose bool
}

func New(app AppInfo) (Args, error) {
    parser := argparse.NewParser(app.Name, app.Description)

	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Prints the version"})
	debug := parser.Flag("d", "debug", &argparse.Options{Required: false, Help: "Enable debug logging"})
	verbose := parser.Flag("r", "verbose", &argparse.Options{Required: false, Help: "Logs more (verbose) information"})
	config := parser.String("c", "config", &argparse.Options{Required: false, Help: "Specify the configuration path"})

	err := parser.Parse(os.Args)
	if err != nil {
		return Args{}, fmt.Errorf("%+v", parser.Usage(err))
	}

    var args = Args{ConfigPath: "", Debug: false, Verbose: false}
	if *version {
		fmt.Printf("v%s\n", app.Version)
		os.Exit(0)
	}

	if *debug {
		args.Debug = true
	}

    if *verbose {
        args.Verbose = true
    }

	if len(*config) >= 1 {
        args.ConfigPath = *config
    }

    return args, nil
}
