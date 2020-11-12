package client

import "strings"

const (
	CachePath        = "o/%s/e/%s/caches"
	CachePathGet     = CachePath + "/%s"
	CacheIdSeparator = ":"
)

type Cache struct {
	EnvironmentName string
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	ExpirySettings  Expiration `json:"expirySettings,omitempty"`
	//OverflowToDisk                    bool       `json:"overflowToDisk,omitempty"`
	SkipCacheIfElementSizeInKBExceeds int `json:"skipCacheIfElementSizeInKBExceeds,omitempty"`
}

type Expiration struct {
	TimeoutInSec *ExpiryValue `json:"timeoutInSec,omitempty"`
	TimeOfDay    *ExpiryValue `json:"timeOfDay,omitempty"`
	ExpiryDate   *ExpiryValue `json:"expiryDate,omitempty"`
}

type ExpiryValue struct {
	Value string `json:"value,omitempty"`
}

func (c *Cache) CacheEncodeId() string {
	return c.EnvironmentName + CacheIdSeparator + c.Name
}

func CacheDecodeId(s string) (string, string) {
	tokens := strings.Split(s, CacheIdSeparator)
	return tokens[0], tokens[1]
}
