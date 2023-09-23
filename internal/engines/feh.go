package engines

import (
	"fmt"
	"os/exec"
)

func (e *Feh) SetWallpaper() error {
	cmd := fmt.Sprintf("feh --bg-scale %s &", e.path)
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}
	return nil
}

func (e *Feh) SetWallpaperPath(path string) {
	e.path = path
}
