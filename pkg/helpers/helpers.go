package helpers

import (
	"os"
    "io/fs"
    "fmt"
	"os/user"
	"path/filepath"
	"strings"
    "crypto/sha256"
)

func SanitizePath(path string) string {
	var sPath string

	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" || path == "$HOME" {
		sPath = dir
	} else if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "$HOME/") {
		sPath = filepath.Join(dir, path[2:])
	} else {
		sPath = path
	}

	return sPath
}

func FileExists(path string) bool {
	if info, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
    return true
}

func GetFilesInDir(path string) ([]string, error) {
    var files []string

    if !PathExists(path) {
        return files, fmt.Errorf("Path '%s' does not exist", path)
    }

    // recursivly search for images in provided path
    err := filepath.WalkDir(path, func(path string, dir fs.DirEntry, err error) error {
        if !dir.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    if err != nil {
        return files, err
    }

    return files, nil
}

func CreateHash(str string) string {
    hash := sha256.New()
    hash.Write([]byte(str))
    return fmt.Sprintf("%x", hash.Sum(nil))
}