package cache

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	log "github.com/mikeunge/WallpaperEngine/pkg/logger"
)

func Find(path string, hash string, size int) bool {
	if !helpers.FileExists(path) {
		update(path, hash, size)
		return false
	}

	data, err := read(path)
	if err != nil {
		log.Error("Cannot read from cache file, %s", path)
		return false
	}

	for _, elem := range data {
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
		write(path, hash)
		return nil
	}

	cache, err := read(path)
	if err != nil {
		return fmt.Errorf("Cannot read from cache file, %s", path)
	}

	if len(cache) >= size {
		cache = pop(cache)
	}

	cache = append(cache, hash)
	write(path, strings.Join(cache, ";"))
	return nil
}

func read(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot read cache from %s, %+v", path, err)
	}
	return strings.Split(string(data), ";"), nil
}

func write(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

func pop(cache []string) []string {
	return cache[1:]
}
