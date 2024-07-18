package engines

import "fmt"

type Engine struct {
	path string
}

type EngineImpl interface {
	SetWallpaper() error
	SetWallpaperPath(string)
}

func EngineFactory(engineType string) (EngineImpl, error) {
	switch engineType {
	case "gnome":
		return &Gnome{}, nil
	case "feh":
		return &Feh{}, nil
	default:
		return nil, fmt.Errorf("could not load engine type: %s", engineType)
	}
}
