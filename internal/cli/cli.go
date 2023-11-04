package cli

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
)

type AppInfo struct {
	Name          string
	Description   string
	Version       string
	Author        string
	WallpaperPath string
}

type Args struct {
	ConfigPath string
	Debug      bool
	Verbose    bool
	Wallpaper  string
}

func New(app AppInfo) (Args, error) {
	parser := argparse.NewParser(app.Name, app.Description)

	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Prints the version"})
	debug := parser.Flag("d", "debug", &argparse.Options{Required: false, Help: "Enable debug logging"})
	verbose := parser.Flag("r", "verbose", &argparse.Options{Required: false, Help: "Logs more (verbose) information"})
	config := parser.String("c", "config", &argparse.Options{Required: false, Help: "Specify the configuration path"})
	wallpaper := parser.String("s", "set-wallpaper", &argparse.Options{Required: false, Help: "Specify a wallpaper that should get set"})
	currentWallpaper := parser.Flag("g", "get-wallpaper", &argparse.Options{Required: false, Help: "Get the current wallpaper"})

	err := parser.Parse(os.Args)
	if err != nil {
		return Args{}, fmt.Errorf("%+v", parser.Usage(err))
	}

	if *version {
		fmt.Printf("v%s\n", app.Version)
		os.Exit(0)
	}

	if *currentWallpaper {
        if !helpers.FileExists(app.WallpaperPath) {
            fmt.Printf("Could not get current wallpaper, file does not exist: %s\n", app.WallpaperPath)
            os.Exit(1)
        }

        data, err := helpers.ReadFile(app.WallpaperPath)
        if err != nil {
            fmt.Printf("Could not read data: %+v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Current wallpaper: %s\n", string(data))
		os.Exit(0)
	}

	var args = Args{ConfigPath: *config, Debug: false, Verbose: false, Wallpaper: *wallpaper}
	if *debug {
		args.Debug = true
	}

	if *verbose && !*debug {
		args.Verbose = true
	}
	return args, nil
}
