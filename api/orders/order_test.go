package orders

import (
	"jordanfinners/api/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type calculateOrderTest struct {
	items         int
	expectedPacks []model.Pack
}

var calculateOrderTests = map[string]calculateOrderTest{
	"1": {
		items:         1,
		expectedPacks: []model.Pack{{Quantity: 250}},
	},
	"250": {
		items:         250,
		expectedPacks: []model.Pack{{Quantity: 250}},
	},
	"251": {
		items:         251,
		expectedPacks: []model.Pack{{Quantity: 500}},
	},
	"501": {
		items:         501,
		expectedPacks: []model.Pack{{Quantity: 500}, {Quantity: 250}},
	},
	"12001": {
		items:         12001,
		expectedPacks: []model.Pack{{Quantity: 5000}, {Quantity: 5000}, {Quantity: 2000}, {Quantity: 250}},
	},
	"10900": {
		items:         10900,
		expectedPacks: []model.Pack{{Quantity: 5000}, {Quantity: 5000}, {Quantity: 1000}},
	},
	"90909": {
		items:         90909,
		expectedPacks: []model.Pack{{Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 5000}, {Quantity: 1000}},
	},
}

func TestCalculatingOrders(t *testing.T) {
	packs := []model.Pack{{Quantity: 5000}, {Quantity: 2000}, {Quantity: 1000}, {Quantity: 500}, {Quantity: 250}}

	for name, test := range calculateOrderTests {
		t.Run(name, func(t *testing.T) {
			computedOrder := CalculateOrder(test.items, packs)

			require.Equal(t, test.items, computedOrder.Items)
			require.ElementsMatch(t, test.expectedPacks, computedOrder.Packs)
		})
	}
}

func TestCalculatingOrdersWith99EdgeCase(t *testing.T) {
	packs := []model.Pack{{Quantity: 90}, {Quantity: 33}}
	computedOrder := CalculateOrder(99, packs)

	require.Equal(t, 99, computedOrder.Items)
	require.ElementsMatch(t, []model.Pack{{Quantity: 33}, {Quantity: 33}, {Quantity: 33}}, computedOrder.Packs)
}
