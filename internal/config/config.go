package config

import (
	"encoding/json"
	"fmt"

	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
)

type IRemember struct {
	RememberSetWallpapers bool   `json:"remember_set_wallpapers"`
	MaxRotations          int    `json:"max_rotations"`
	RememberPath          string `json:"remember_path"`
}

type Config struct {
	Engine          string    `json:"engine"`
	WallpaperPath   string    `json:"wallpaper_path"`
	Wallpapers      []string  `json:"wallpapers"`
	Blacklist       []string  `json:"blacklist"`
	RandomWallpaper bool      `json:"random_wallpaper"`
	Remember        IRemember `json:"remember"`
	ValidExtensions []string  `json:"valid_extensions"`
}

func Parse(configPath string) (Config, error) {
	var config Config

	fileExists := helpers.FileExists(configPath)
	if !(fileExists) {
		return config, fmt.Errorf("config not found: %s", configPath)
	}

	data, err := helpers.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("could not read config: %+v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error while parsing: %+v, make sure the file is a correct json file", err)
	}

	return config, nil
}

func DefaultConfig() Config {
	return Config{
		Engine:          "gnome",
		WallpaperPath:   "~/Pictures/",
		Wallpapers:      []string{},
		Blacklist:       []string{},
		RandomWallpaper: true,
		Remember: IRemember{
			RememberSetWallpapers: true,
			MaxRotations:          5,
			RememberPath:          "~/.wpe.store",
		},
		ValidExtensions: []string{
			".jpg",
			".jpeg",
			".png",
		},
	}
}
