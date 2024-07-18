package engines

import (
	"fmt"
	"os/exec"
	"strings"
)

type Gnome struct {
	Engine
}

func getPictureUri() (string, error) {
	cmd := "gsettings get org.gnome.desktop.interface color-scheme"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	// for gnome (in ubuntu at least) you need to define what colorscheme is used
	// so we need this little helper that checks the selected (prefered) colorscheme and returns the "picture-uri" to set
	if strings.Contains(string(out), "prefer-dark") {
		return "picture-uri-dark", nil
	}
	return "picture-uri", nil
}

func (engine *Gnome) SetWallpaper() error {
	pictureUri, err := getPictureUri()
	if err != nil {
		return fmt.Errorf("could not get picture uri (colorschema), %+v", err)
	}
	cmd := fmt.Sprintf("gsettings set org.gnome.desktop.background %s file://%s", pictureUri, engine.path)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}
	return nil
}

func (engine *Gnome) SetWallpaperPath(path string) {
	engine.path = path
}
