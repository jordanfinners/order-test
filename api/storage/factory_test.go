package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGettingDBConnection(t *testing.T) {
	os.Setenv("DATABASE_NAME", "TESTING")
	client := NewClient()

	name := client.db.Name()
	require.Equal(t, "TESTING", name)
}
