package engines

import "fmt"

type Gnome struct {
	path string
}

type Feh struct {
	path string
}

type Engine interface {
	SetWallpaper() error
	SetWallpaperPath(string)
}

func EngineFactory(engineType string) (Engine, error) {
	switch engineType {
	case "gnome":
		return &Gnome{}, nil
	case "feh":
		return &Feh{}, nil
	default:
		return nil, fmt.Errorf("could not load engine type: %s", engineType)
	}
}
