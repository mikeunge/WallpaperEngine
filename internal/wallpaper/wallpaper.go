package wallpaper

import (
	"fmt"
	"math/rand"

	"github.com/mikeunge/WallpaperEngine/internal/config"
	"github.com/mikeunge/WallpaperEngine/pkg/cache"
	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	"github.com/mikeunge/WallpaperEngine/pkg/image_helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

func GetWallpaper(appConfig *config.Config, wallpaper string) (string, error) {
	// Check if wallpaper with full path was provided & check if it the file exists
	if len(wallpaper) > 0 && helpers.FileExists(wallpaper) {
		log.Info("Provided wallpaper was found, %s", wallpaper)
		return wallpaper, nil
	}

	images, err := getFilteredWallpapers(appConfig)
	if err != nil {
		return "", err
	}

	// Check if we want to cache the wallpapers or not or we want to choose a random one
	if appConfig.RandomWallpaper || !appConfig.Remember.RememberSetWallpapers || appConfig.Remember.MaxRotations >= len(images) {
		log.Info("Choosing random image.")
		return getRandomImage(images), nil
	} else if len(wallpaper) > 0 {
		image, err := image_helpers.FindImageInImages(wallpaper, images)
		if err != nil {
			log.Warn("Filtered through all the images but there seems to be no an error: %+v", err)
			log.Info("Returning random image")
			return getRandomImage(images), nil
		}
		log.Debug("Found image! %s", image)
		return image, nil
	} else {
		return getImageWithCacheCheck(images, appConfig.Remember.RememberPath, appConfig.Remember.MaxRotations), nil
	}
}

func SaveCurrentWallpaper(path string, data string) error {
	return helpers.WriteToFile(path, data)
}

func getRandomImage(imagePaths []string) string {
	return imagePaths[rand.Intn(len(imagePaths))]
}

func getFilteredWallpapers(appConfig *config.Config) ([]string, error) {
	files, err := helpers.GetFilesInDir(appConfig.WallpaperPath)
	if err != nil {
		log.Error("%+v", err)
		return []string{}, err
	} else if len(files) == 0 {
		return []string{}, fmt.Errorf("no files found in %s", appConfig.WallpaperPath)
	}

	images, err := image_helpers.FilterImages(files, appConfig.ValidExtensions, appConfig.Blacklist)
	if err != nil {
		log.Error("%+v", err)
		return []string{}, err
	}

	return images, nil
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
