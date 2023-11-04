package cmd

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/mikeunge/WallpaperEngine/internal/cli"
	"github.com/mikeunge/WallpaperEngine/internal/config"
	"github.com/mikeunge/WallpaperEngine/internal/engines"
	"github.com/mikeunge/WallpaperEngine/internal/wallpaper"
	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

const (
	AppName        = "WallpaperEngine"
	AppDescription = "A simple wallpaper setter."
	AppVersion     = "1.0.3"
	AppAuthor      = "@mikeunge"
)

var (
	configPath           = "~/.config/wallpaper-engine/config.json"
	currentWallpaperPath = "~/.wpe"
	wp                   string
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
		log.SetOutput(false)
		log.Info("WallpaperEngine running in debug mode")
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

	// if the user has wallpapers specified - clean them ;)
	if len(appConfig.Wallpapers) > 0 {
		var cleanWallpapers []string
		for _, w := range appConfig.Wallpapers {
			cleanWallpapers = append(cleanWallpapers, helpers.SanitizePath(w))
		}
		appConfig.Wallpapers = cleanWallpapers
	}
}

func App() int {
	appConfig, err := config.Parse(configPath)
	if err != nil {
		log.Error("Config parser returned an error, '%+v'", err)
		log.Warn("Using default config as fallback!")
		appConfig = config.DefaultConfig()
	}
	sanitizePaths(&appConfig)

	/**
	 * Conditions for the wallpapers specified in the config:
	 *  - Make sure __NO__ wallpaper was set (-s);
	 *  - If there is at least __ONE__ wallpaper in the wallpapers array;
	 *  - And if random wallpaper is __NOT__ active;
	 */
	if !(len(wp) > 0) && len(appConfig.Wallpapers) > 0 && !appConfig.RandomWallpaper {
		log.Info("Wallpaper(s) provided - let's choose one")
		wp = appConfig.Wallpapers[rand.Intn(len(appConfig.Wallpapers))]
		// set the config to false because we don't remember pre-defined wallpapers
		appConfig.Remember.RememberSetWallpapers = false
	}

	wp, err = wallpaper.GetWallpaper(&appConfig, wp)
	if err != nil {
		return 1
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

	// safe current wallpaper to tmp file
	err = wallpaper.SaveCurrentWallpaper(currentWallpaperPath, wp)
	if err != nil {
		log.Error("Could not write current wallpaper to file, err: %+v", err)
		return 1
	}
	return 0
}
