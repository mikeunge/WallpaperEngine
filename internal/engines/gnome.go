package engines

import (
	"fmt"
	"os/exec"
	"strings"
)

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

func (e *Gnome) SetWallpaper() error {
	pictureUri, err := getPictureUri()
	if err != nil {
		return fmt.Errorf("could not get picture uri (colorschema), %+v", err)
	}
	cmd := fmt.Sprintf("gsettings set org.gnome.desktop.background %s file://%s", pictureUri, e.path)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}
	return nil
}

func (e *Gnome) SetWallpaperPath(path string) {
	e.path = path
}
