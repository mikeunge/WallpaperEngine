package engines

import (
	"fmt"
	"os/exec"
)

type Feh struct {
	Engine
}

func (engine *Feh) SetWallpaper() error {
	cmd := fmt.Sprintf("feh --bg-scale %s &", engine.path)
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}
	return nil
}

func (engine *Feh) SetWallpaperPath(path string) {
	engine.path = path
}
