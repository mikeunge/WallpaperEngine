package cmd

import (
	"fmt"
	"os"

	"github.com/mikeunge/Scripts/wallpaper-engine/internal/cli"
	"github.com/mikeunge/Scripts/wallpaper-engine/internal/config"
	"github.com/mikeunge/Scripts/wallpaper-engine/pkg/cache"
	"github.com/mikeunge/Scripts/wallpaper-engine/pkg/helpers"
	log "github.com/mikeunge/Scripts/wallpaper-engine/pkg/logger"
)

const (
    AppName = "Wallpaper-Engine"
    AppDescription = "A simple wallpaper setter."
    AppVersion = "0.1"
    AppAuthor = "@mikeunge"
)

var (
    configPath = "~/.config/wallpaper-engine/config.json"
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

func sanitizeConfigPaths(appConfig *config.Config) {
    appConfig.WallpaperPath = helpers.SanitizePath(appConfig.WallpaperPath)
    appConfig.Remember.RememberPath = helpers.SanitizePath(appConfig.Remember.RememberPath)
}

func getWallpaper(appConfig *config.Config) (string, error) {
    files, err := helpers.GetFilesInDir(appConfig.WallpaperPath)
    if err != nil {
        log.Error("%+v", err)
        return "", err
    }

    images, err := helpers.FilterImages(files, appConfig.ValidExtensions)
    if err != nil {
        log.Error("%+v", err)
        return "", err
    }

    // we don't remember the last images, so fuck them
    if !appConfig.Remember.RememberSetWallpapers {
        return helpers.GetRandomImage(images), nil
    }

    image := helpers.GetRandomImage(images)
    log.Debug("Selected image: %s", image)
    fileHash := helpers.CreateHash(image)
    cache.Find(appConfig.Remember.RememberPath, fileHash, appConfig.Remember.MaxRotations)

    return "", nil
}



func App() error {
	appConfig, err := config.Parse(configPath)
	if err != nil {
        log.Error("Config parser returned an error, '%+v'", err)
        log.Warn("Using default config as fallback!")
		appConfig = config.DefaultConfig()
	}
    sanitizeConfigPaths(&appConfig)

    wallpaper, err := getWallpaper(&appConfig)
    if err != nil {
        return err
    }

    fmt.Println(wallpaper)

    return  nil
}
