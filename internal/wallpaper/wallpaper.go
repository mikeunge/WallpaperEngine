package wallpaper

import (
	"math/rand"

	"github.com/mikeunge/WallpaperEngine/internal/config"
	"github.com/mikeunge/WallpaperEngine/pkg/cache"
	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

func GetWallpaper(appConfig *config.Config, wallpaper string) (string, error) {
	// Check if wp was provided & check if it exists
	if len(wallpaper) > 0 {
		if helpers.FileExists(wallpaper) {
			log.Info("Provided wallpaper was found, %s", wallpaper)
			return wallpaper, nil
		}

		wallpaper = helpers.GetFileName(wallpaper)
		log.Debug("Retrieved filename: %s", wallpaper)
	}

	files, err := helpers.GetFilesInDir(appConfig.WallpaperPath)
	if err != nil {
		log.Error("%+v", err)
		return "", err
	}

	images, err := helpers.FilterImages(files, appConfig.ValidExtensions)
	if err != nil {
		log.Error("%+v", err)
		return "", err
	}

	if len(wallpaper) > 0 {
		image, err := helpers.FindImageInImages(wallpaper, images)
		if err != nil {
			log.Warn("Filtered through all the images but there seems to be no an error: %+v", err)
			log.Info("Returning random image")
			return getRandomImage(images), nil
		}
		log.Debug("Found image! %s", image)
		return image, nil
	} else if !appConfig.Remember.RememberSetWallpapers {
		return getRandomImage(images), nil
	} else if !appConfig.RandomWallpaper && len(appConfig.Wallpapers) <= 0 {
		log.Warn("Choosing random image because no wallpapers are defined!")
		return getRandomImage(images), nil
	} else if appConfig.Remember.MaxRotations >= len(images) {
		log.Warn("Choosing random image because store size is bigger than all the available images, to fix this adapt your config!")
		return getRandomImage(images), nil
	}

	return getImageWithCacheCheck(images, appConfig.Remember.RememberPath, appConfig.Remember.MaxRotations), nil
}

func SaveCurrentWallpaper(path string, data string) error {
	return helpers.WriteToFile(path, data)
}

func getRandomImage(imagePaths []string) string {
	return imagePaths[rand.Intn(len(imagePaths))]
}

func getImageWithCacheCheck(images []string, rememberPath string, maxRotations int) string {
	var image string

	for {
		image = getRandomImage(images)
		log.Debug("Selected image: %s", image)
		fileHash := helpers.CreateHash(image)

		if !cache.Find(rememberPath, fileHash, maxRotations) {
			log.Info("Using image: %s", image)
			break
		}
		log.Debug("Image '%s' found in cache, looking for another one", image)
	}

	return image
}
