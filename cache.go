package osutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CacheObjectType defines a unique type for a cache object.
// This type should be descriptive to the object that is cached.
type CacheObjectType string

// CacheObjectWrite write data with the given unique identifier and type to the cache.
// The meta data will be written as a human-readable data to identify the cache object.
func CacheObjectWrite(cachePath string, identifier string, cacheObjectType CacheObjectType, data any, meta map[string]string) (err error) {
	cacheObjectPathRelative := cacheObjectPath(identifier)
	cacheObjectPath := filepath.Join(cachePath, cacheObjectPathRelative)

	if err := os.MkdirAll(cacheObjectPath, 0755); err != nil {
		return err
	}

	dataFile, err := os.Create(filepath.Join(cacheObjectPath, string(cacheObjectType)+".gob"))
	if err != nil {
		return err
	}
	defer func() {
		if e := dataFile.Close(); e != nil {
			if err == nil {
				err = e
			} else {
				err = errors.Join(err, e)
			}
		}
	}()
	if err := gob.NewEncoder(dataFile).Encode(data); err != nil {
		return err
	}
	if err := dataFile.Sync(); err != nil {
		return err
	}

	metaFile, err := os.Create(filepath.Join(cacheObjectPath, string(cacheObjectType)+".json"))
	if err != nil {
		return err
	}
	defer func() {
		if e := metaFile.Close(); e != nil {
			if err == nil {
				err = e
			} else {
				err = errors.Join(err, e)
			}
		}
	}()
	if err := json.NewEncoder(metaFile).Encode(meta); err != nil {
		return err
	}
	if err := metaFile.Sync(); err != nil {
		return err
	}

	return nil
}

// CacheObjectRead reads data with the given unique identifier and type from the cache.
func CacheObjectRead(cachePath string, identifier string, cacheObjectType CacheObjectType, data any) (exists bool, err error) {
	cacheObjectPathRelative := cacheObjectPath(identifier)
	cacheObjectPath := filepath.Join(cachePath, cacheObjectPathRelative)

	raw, err := os.ReadFile(filepath.Join(cacheObjectPath, string(cacheObjectType)+".gob"))
	if err != nil {
		return false, nil
	}

	if err := gob.NewDecoder(bytes.NewReader(raw)).Decode(data); err != nil {
		return false, err
	}

	return true, nil
}

// cacheObjectPath returns the relative object path for the given identifier.
func cacheObjectPath(identifier string) (cacheObjectPathRelative string) {
	checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(identifier)))

	return filepath.Join(checksum[0:2], checksum[2:])
}
