package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type calculateOrderTest struct {
	items         int
	expectedPacks []Pack
}

var calculateOrderTests = map[string]calculateOrderTest{
	"1": {
		items:         1,
		expectedPacks: []Pack{{Quantity: 250}},
	},
	"250": {
		items:         250,
		expectedPacks: []Pack{{Quantity: 250}},
	},
	"251": {
		items:         251,
		expectedPacks: []Pack{{Quantity: 500}},
	},
	"501": {
		items:         501,
		expectedPacks: []Pack{{Quantity: 500}, {Quantity: 250}},
	},
	"12001": {
		items:         12001,
		expectedPacks: []Pack{{Quantity: 5000}, {Quantity: 5000}, {Quantity: 2000}, {Quantity: 250}},
	},
}

func TestCalculatingOrders(t *testing.T) {
	packs := []Pack{{Quantity: 5000}, {Quantity: 2000}, {Quantity: 1000}, {Quantity: 500}, {Quantity: 250}}

	for name, test := range calculateOrderTests {
		t.Run(name, func(t *testing.T) {
			computedOrder := calculateOrder(test.items, packs)

			require.Equal(t, test.items, computedOrder.Items)
			require.Equal(t, test.expectedPacks, computedOrder.Packs)
		})
	}
}
