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
	CachePath     string
}

type Args struct {
	ConfigPath string
	Debug      bool
	Wallpaper  string
	ResetCache bool
}

func New(app AppInfo) (Args, error) {
	parser := argparse.NewParser(app.Name, app.Description)

	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Prints the version"})
	debug := parser.Flag("d", "debug", &argparse.Options{Required: false, Help: "Enable debug logging"})
	resetCache := parser.Flag("r", "reset", &argparse.Options{Required: false, Help: "Reset/Clear the wpe (cache) store"})
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

	return Args{ConfigPath: *config, Debug: *debug, Wallpaper: *wallpaper, ResetCache: *resetCache}, nil
}
