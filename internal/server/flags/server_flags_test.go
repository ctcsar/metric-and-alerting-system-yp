package server

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetServerURL(t *testing.T) {
	f := &serverFlags{}
	t.Run("GetServerUrlFromEnv", func(t *testing.T) {
		os.Setenv("ADDRESS", "localhost:9010")
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:9010", f.GetServerURL())
	})
	t.Run("GetEmptyServerURL", func(t *testing.T) {
		os.Unsetenv("ADDRESS")
		f.url = ""
		// Check that the flags are set correctly
		assert.Equal(t, "", f.GetServerURL())
	})

	t.Run("GetCustomServerURL", func(t *testing.T) {
		f.url = "localhost:9000"
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:9000", f.GetServerURL())
	})

}

func TestGetStoragePath(t *testing.T) {
	f := &serverFlags{}
	t.Run("GetStoragePathFromEnv", func(t *testing.T) {
		os.Setenv("FILE_STORAGE_PATH", "../../data/storage_test.json")
		// Check that the flags are set correctly
		assert.Equal(t, "../../data/storage_test.json", f.GetStoragePath())
	})
	t.Run("GetEmptyStoragePath", func(t *testing.T) {
		os.Unsetenv("FILE_STORAGE_PATH")
		// Check that the flags are set correctly
		assert.Equal(t, "../../data/storage_test.json", f.GetStoragePath())
	})

	t.Run("GetCustomStoragePath", func(t *testing.T) {
		f.storagePath = "../../data/storage_not_empty.json"
		// Check that the flags are set correctly
		assert.Equal(t, "../../data/storage_not_empty.json", f.GetStoragePath())
	})
}

func TestGetRestore(t *testing.T) {
	f := &serverFlags{}

	t.Run("GetFalseRestore", func(t *testing.T) {
		f.restore = false
		// Check that the flags are set correctly
		assert.Equal(t, false, f.GetRestore())
	})

	t.Run("GetTrueRestore", func(t *testing.T) {
		f.restore = true
		// Check that the flags are set correctly
		assert.Equal(t, true, f.GetRestore())
	})
}

func TestGetStoreInterval(t *testing.T) {
	f := &serverFlags{}

	t.Run("GetStoreIntervalFromEnv", func(t *testing.T) {
		os.Setenv("STORE_INTERVAL", "100")
		// Check that the flags are set correctly
		assert.Equal(t, time.Duration(100), f.GetStoreInterval())
	})

	t.Run("GetCustomStoreInterval", func(t *testing.T) {
		os.Unsetenv("STORE_INTERVAL")
		f.storeInterval = 120
		// Check that the flags are set correctly
		assert.Equal(t, time.Duration(120), f.GetStoreInterval())
	})
}
