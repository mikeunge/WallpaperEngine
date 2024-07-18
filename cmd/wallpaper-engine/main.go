package main

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
	AppVersion     = "1.0.6"
	AppAuthor      = "@mikeunge"
)

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

func main() {
	var selectedWallpaper string

	configPath := helpers.SanitizePath("~/.config/wallpaper_engine/config.json")
	currentWallpaperPath := helpers.SanitizePath("~/.wpe")
	cliArgs, err := cli.New(cli.AppInfo{Name: AppName, Description: AppDescription, Version: AppVersion, Author: AppAuthor, WallpaperPath: currentWallpaperPath})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(0)
	}

	if cliArgs.Debug {
		log.SetLogLevel("debug")
		log.SetOutput(false)
		log.Info("WallpaperEngine running in debug mode")
	}

	if len(cliArgs.ConfigPath) > 0 {
		log.Info(fmt.Sprintf("Config path changed '%s' to '%s'", configPath, cliArgs.ConfigPath))
		configPath = cliArgs.ConfigPath
	}

	if len(cliArgs.Wallpaper) > 0 {
		log.Info(fmt.Sprintf("Wallpaper was specified, using: %s", cliArgs.Wallpaper))
		selectedWallpaper = cliArgs.Wallpaper
	}

	appConfig, err := config.Parse(configPath)
	if err != nil {
		log.Error("Config parser returned an error, '%+v'", err)
		log.Warn("Using default config as fallback!")
		appConfig = config.DefaultConfig()
	}
	sanitizePaths(&appConfig)

	if cliArgs.ResetCache {
		log.Warn("Resetting cache (%s)", appConfig.Remember.RememberPath)
		if err = helpers.WriteToFile(appConfig.Remember.RememberPath, ""); err != nil {
			fmt.Printf("Could not reset cache: %s\n%+v\n", appConfig.Remember.RememberPath, err)
		}
	}

	engine, err := engines.EngineFactory(appConfig.Engine)
	if err != nil {
		log.Error("Error while loading engine: %+v", err)
		os.Exit(1)
	}

	/**
	 * Conditions for the wallpapers specified in the config:
	 *  - Make sure __NO__ wallpaper was set (-s);
	 *  - If there is at least __ONE__ wallpaper in the wallpapers array;
	 *  - And if random wallpaper is __NOT__ active;
	 */
	if !(len(selectedWallpaper) > 0) && len(appConfig.Wallpapers) > 0 && !appConfig.RandomWallpaper {
		log.Info("Wallpaper(s) provided - let's choose one")
		selectedWallpaper = appConfig.Wallpapers[rand.Intn(len(appConfig.Wallpapers))]
		// set the config to false because we don't remember pre-defined wallpapers
		appConfig.Remember.RememberSetWallpapers = false
	}

	selectedWallpaper, err = wallpaper.GetWallpaper(&appConfig, selectedWallpaper)
	if err != nil {
		os.Exit(1)
	}

	engine.SetWallpaperPath(selectedWallpaper)
	if err = engine.SetWallpaper(); err != nil {
		log.Error("Something went wrong while setting wallpaper, err: %+v", err)
		os.Exit(1)
	}

	// safe current wallpaper to tmp file
	if err = wallpaper.SaveCurrentWallpaper(currentWallpaperPath, selectedWallpaper); err != nil {
		log.Error("Could not write current wallpaper to file, err: %+v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
