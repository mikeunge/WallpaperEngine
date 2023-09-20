package helpers

import (
    "fmt"
    "strings"

    log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

func FilterImages(imagePaths []string, validExtensions []string) ([]string, error) {
    if len(validExtensions) == 0 || len(imagePaths) == 0 {
        return imagePaths, fmt.Errorf("Cannot filter images, available images: '%d' - available extensions: '%d'", len(imagePaths), len(validExtensions))
    }
    
    var filteredImages []string
    for i := 0; i < len(imagePaths); i++ {
        for j := 0; j < len(validExtensions); j++ {
            if strings.HasSuffix(imagePaths[i], validExtensions[j]) {
                filteredImages = append(filteredImages, imagePaths[i])
                break
            }
            if len(validExtensions) - 1 == j {
                log.Debug("Dropping image '%s'", imagePaths[i])
            }
        }
    }

    return filteredImages, nil
}

