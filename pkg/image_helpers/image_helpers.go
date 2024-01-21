package image_helpers

import (
	"fmt"
	"strings"

	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

func FilterImages(imagePaths []string, validExtensions []string, blacklist []string) ([]string, error) {
	var filteredImages []string
	for _, image := range imagePaths {
		if len(validExtensions) > 0 && !hasValidExtension(image, validExtensions) {
			log.Debug("Skipping image '%s' - has no valid extensions", image)
			continue
		}

		if isBlacklisted(image, blacklist) {
			log.Info("Skipping image '%s' - image/path is blacklisted", image)
			continue
		}
		filteredImages = append(filteredImages, image)
	}

	return filteredImages, nil
}

func FindImageInImages(image string, images []string) (string, error) {
	hasExtension := true
	if len(strings.Split(image, ".")) == 1 {
		log.Warn("Looking for image WITHOUT provided extension. Be aware that this can lead to unexpected behaviour!")
		hasExtension = false
	}

	for _, img := range images {
		var file string

		if hasExtension {
			file = helpers.GetFileName(img)
		} else {
			file = helpers.GetFileNameWithoutExtension(img)
		}

		if file == image {
			return img, nil
		}
	}
	return "", fmt.Errorf("could not find image, %s", image)
}

func hasValidExtension(image string, extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}

	for _, extension := range extensions {
		if strings.HasSuffix(image, extension) {
			return true
		}
	}
	return false
}

func isBlacklisted(path string, blacklist []string) bool {
	if len(blacklist) == 0 {
		return false
	}

	for i := 0; i < len(blacklist); i++ {
		if path == blacklist[i] || helpers.GetFileName(path) == blacklist[i] {
			return true
		}
	}

	return false
}
