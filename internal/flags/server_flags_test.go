package flags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerURL(t *testing.T) {
	t.Run("GetServerUrlFromEnv", func(t *testing.T) {
		os.Setenv("ADDRESS", "localhost:8080")
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:8080", GetServerURL())
	})
	t.Run("GetEmptyServerURL", func(t *testing.T) {
		flagRunAddr = ""
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:8080", GetServerURL())
	})

	t.Run("GetCustomServerURL", func(t *testing.T) {
		os.Unsetenv("ADDRESS")
		flagRunAddr = "localhost:9000"
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:9000", GetServerURL())
	})

}
func TestSetServerFlags(t *testing.T) {
	t.Run("SetServerFlags", func(t *testing.T) {
		flagRunAddr = "localhost:8080"
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:8080", flagRunAddr)
	})
	t.Run("SetEmptyServerFlags", func(t *testing.T) {
		flagRunAddr = ""
		// Check that the flags are set correctly
		assert.Equal(t, "", flagRunAddr)
	})
}
