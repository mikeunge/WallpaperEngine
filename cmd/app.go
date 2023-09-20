package cmd

import (
	"fmt"
	"os"

	"github.com/mikeunge/WallpaperEngine/internal/cli"
	"github.com/mikeunge/WallpaperEngine/internal/config"
	"github.com/mikeunge/WallpaperEngine/internal/wallpaper"
	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

const (
    AppName = "Wallpaper-Engine"
    AppDescription = "A simple wallpaper setter."
    AppVersion = "0.1"
    AppAuthor = "@mikeunge"
)

var (
    configPath = "~/.config/WallpaperEngine/config.json"
)

func init() {
    appInfo := cli.AppInfo{Name: AppName, Description: AppDescription, Version: AppVersion, Author: AppAuthor}

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

    configPath = helpers.SanitizePath(configPath)
}

func sanitizePaths(appConfig *config.Config) {
    appConfig.WallpaperPath = helpers.SanitizePath(appConfig.WallpaperPath)
    appConfig.Remember.RememberPath = helpers.SanitizePath(appConfig.Remember.RememberPath)
}


func App() error {
	appConfig, err := config.Parse(configPath)
	if err != nil {
        log.Error("Config parser returned an error, '%+v'", err)
        log.Warn("Using default config as fallback!")
		appConfig = config.DefaultConfig()
	}
    sanitizePaths(&appConfig)

    wallpaper, err := wallpaper.GetWallpaper(&appConfig)
    if err != nil {
        return err
    }

    fmt.Println(wallpaper)

    return  nil
}
