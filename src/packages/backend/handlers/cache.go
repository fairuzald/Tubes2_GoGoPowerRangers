package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// CacheManager is a singleton to manage cache and JSON operations.
type CacheManager struct {
	cache map[string][]string
	mu    sync.Mutex
}

var (
	globalCache      *CacheManager
	cacheInitialized bool
)

// InitCache initializes the global cache if it's not already initialized.
func InitCache() {
	if !cacheInitialized {
		globalCache = &CacheManager{
			cache: make(map[string][]string),
		}
		LoadCacheFromJSON("cache.json")
		cacheInitialized = true
	}
}

// SaveCacheToJSON saves the cache to a JSON file.
func SaveCacheToJSON(filename string) error {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	data, err := json.Marshal(globalCache.cache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadCacheFromJSON loads the cache from a JSON file.
func LoadCacheFromJSON(filename string) {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("No cache file found. Starting with empty cache.")
		return
	}

	err = json.Unmarshal(file, &globalCache.cache)
	if err != nil {
		fmt.Println("Error reading cache file:", err)
		return
	}
}

// GetLinksFromCache retrieves links from the cache if available.
func GetLinksFromCache(key string) ([]string, bool) {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	value, ok := globalCache.cache[key]
	return value, ok
}

// SetLinksToCache sets links to the cache.
func SetLinksToCache(key string, value []string) {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	globalCache.cache[key] = value
}

func SaveCacheHandler(c *gin.Context) {

	err := SaveCacheToJSON("cache.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cache saved to " + "cache.json"})
}
