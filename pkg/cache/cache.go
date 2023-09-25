package cache

import (
	"fmt"
	"strings"

	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

func Find(path string, hash string, size int) bool {
	if !helpers.FileExists(path) {
		update(path, hash, size)
		return false
	}

	data, err := helpers.ReadFile(path)
	if err != nil {
		log.Error("Cannot read from cache file, %s", path)
		return false
	}

	cacheData := strings.Split(string(data), ";")
	for _, elem := range cacheData {
		if elem == hash {
			log.Debug("Found wallpaper in cache - hash: %s", hash)
			return true
		}
	}
	log.Debug("Hash not found, creating entry - hash: %s", hash)

	err = update(path, hash, size)
	if err != nil {
		log.Error("Error while updating cache file: %s, %+v", path, err)
		return false
	}
	return false
}

func update(path string, hash string, size int) error {
	if !helpers.FileExists(path) {
		log.Warn("Cache file does not exist, %s, trying to create it", path)
		helpers.WriteToFile(path, hash)
		return nil
	}

	data, err := helpers.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read from cache file, %s", path)
	}

	cache := strings.Split(string(data), ";")
	for {
		// pop elements from store till there is space for the new cache entry
		// this for loop ensures that the cache is not bigger than the size we defined in the config
		// because this has led to a bug where wpe ran in a endless loop because all elements where stored in cache
		// and no image could be applied - so we need to keep the cache after config downsize at the provided size
		if len(cache) >= size {
			cache = pop(cache)
			continue
		}
		break
	}

	cache = append(cache, hash)
	helpers.WriteToFile(path, strings.Join(cache, ";"))
	return nil
}

func pop(cache []string) []string {
	return cache[1:]
}
