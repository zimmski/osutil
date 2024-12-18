package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheObject(t *testing.T) {
	temporaryPath := t.TempDir()

	dataToBeCached := map[string]string{
		"A": "1",
		"B": "2",
	}
	identifier := "some-identifier"
	typ := CacheObjectType("some-type")

	{
		var dataToBeRead map[string]string
		exists, err := CacheObjectRead(temporaryPath, identifier, typ, &dataToBeRead)
		assert.NoError(t, err)
		assert.False(t, exists)
	}

	assert.NoError(t, CacheObjectWrite(temporaryPath, identifier, typ, dataToBeCached, nil))

	{
		var dataToBeRead map[string]string
		exists, err := CacheObjectRead(temporaryPath, identifier, typ, &dataToBeRead)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, dataToBeCached, dataToBeRead)
	}
}
