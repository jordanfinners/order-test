package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGettingPacks(t *testing.T) {
	packs, err := getPacks()
	require.NoError(t, err)

	expected := []Pack{{Quantity: 5000}, {Quantity: 2000}, {Quantity: 1000}, {Quantity: 500}, {Quantity: 250}}
	require.Equal(t, expected, packs)
}
