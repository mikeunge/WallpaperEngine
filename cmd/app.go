package cmd

import (
	"fmt"
	"os"

	"github.com/mikeunge/WallpaperEngine/internal/cli"
	"github.com/mikeunge/WallpaperEngine/internal/config"
	"github.com/mikeunge/WallpaperEngine/internal/engines"
	"github.com/mikeunge/WallpaperEngine/internal/wallpaper"
	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

const (
	AppName        = "Wallpaper-Engine"
	AppDescription = "A simple wallpaper setter."
	AppVersion     = "1.0.0"
	AppAuthor      = "@mikeunge"
)

var (
	configPath = "~/.config/wallpaper-engine/config.json"
	currentWallpaperPath = "~/.wpe"
	wp         string
)

func init() {
    currentWallpaperPath = helpers.SanitizePath(currentWallpaperPath)
	appInfo := cli.AppInfo{Name: AppName, Description: AppDescription, Version: AppVersion, Author: AppAuthor, WallpaperPath: currentWallpaperPath}

	args, err := cli.New(appInfo)
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(0)
	}

	if args.Debug {
		log.SetLogLevel("debug")
		log.Info("Wallpaper-Engine running in debug mode")
	} else if args.Verbose {
		log.SetLogLevel("info")
	}

	if len(args.ConfigPath) > 0 {
		log.Info(fmt.Sprintf("Config path changed '%s' to '%s'", configPath, args.ConfigPath))
		configPath = args.ConfigPath
	}

	if len(args.Wallpaper) > 0 {
		log.Info(fmt.Sprintf("Wallpaper was specified, using: %s", args.Wallpaper))
		wp = args.Wallpaper
	}

	configPath = helpers.SanitizePath(configPath)
}

func sanitizePaths(appConfig *config.Config) {
	appConfig.WallpaperPath = helpers.SanitizePath(appConfig.WallpaperPath)
	appConfig.Remember.RememberPath = helpers.SanitizePath(appConfig.Remember.RememberPath)
}

func App() int {
	appConfig, err := config.Parse(configPath)
	if err != nil {
		log.Error("Config parser returned an error, '%+v'", err)
		log.Warn("Using default config as fallback!")
		appConfig = config.DefaultConfig()
	}
	sanitizePaths(&appConfig)

	// check if wallpaper is set from cli
	if len(wp) <= 0 {
		wp, err = wallpaper.GetWallpaper(&appConfig)
		if err != nil {
			return 1
		}
	}

	// get the engine
	engine, err := engines.LoadEngine(appConfig.Engine)
	if err != nil {
		log.Error("Error while loading engine: %+v", err)
		return 1
	}

	engine.SetWallpaperPath(wp)
	err = engine.SetWallpaper()
	if err != nil {
		log.Error("Something went wrong while setting wallpaper, err: %+v", err)
		return 1
	}

    err = helpers.WriteToFile(currentWallpaperPath, wp)
	if err != nil {
		log.Error("Could not write current wallpaper to file, err: %+v", err)
		return 1
	}
	return 0
}
