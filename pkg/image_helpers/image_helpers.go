package image_helpers

import (
	"fmt"
	"strings"

	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

func FilterImages(imagePaths []string, validExtensions []string, blacklist []string) ([]string, error) {
	if len(validExtensions) == 0 || len(imagePaths) == 0 {
		return imagePaths, fmt.Errorf("cannot filter images, available images: '%d' - available extensions: '%d'", len(imagePaths), len(validExtensions))
	}

	var filteredImages []string
	for i := 0; i < len(imagePaths); i++ {
		for j := 0; j < len(validExtensions); j++ {
			if !strings.HasSuffix(imagePaths[i], validExtensions[j]) {
				log.Debug("Dropping file '%s' - image has no valid extension", imagePaths[i])
				continue
			}

			if isBlacklisted(imagePaths[i], blacklist) {
				log.Info("Dropping image '%s' - image/path is blacklisted", imagePaths[i])
				continue
			}
			filteredImages = append(filteredImages, imagePaths[i])
		}
	}

	return filteredImages, nil
}

func FindImageInImages(image string, images []string) (string, error) {
	for i := 0; i < len(images); i++ {
		filename := helpers.GetFileName(images[i])
		if filename == image {
			return images[i], nil
		}
	}

	return "", fmt.Errorf("could not find image, %s", image)
}

func isBlacklisted(path string, blacklist []string) bool {
	// Check if path is in blacklsit and also check for only image name in blacklist
	for i := 0; i < len(blacklist); i++ {
		if path == blacklist[i] {
			return true
		} else if helpers.GetFileName(path) == blacklist[i] {
			return true
		}
	}

	return false
}
