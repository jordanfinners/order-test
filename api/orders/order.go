package orders

import (
	"jordanfinners/api/model"
)

// CalculateOrder calculates the packs required to fulfil the ordered items
// It does this by starting with the largest pack size and working down
// If the items ordered after adding a pack size is still greater than that pack size, then it starts again, using the largest pack sizes possible to fill an order
// If the order is smaller than the smallest pack size available, then use the smallest pack size as the packs cannot be broken up
// Once calculated the packs required to fulfil the order, check if we can use a single pack to fulfil the order instead of multiple as this is more efficient for shipping. This has been kept simple for now but in future should be expanded to check all combinations of packs to see if the order can be simplified.
func CalculateOrder(items int, packs []model.Pack) model.Order {
	var fulfilment []model.Pack
	var quantityRemaining = items
	for quantityRemaining > 0 {
		for i, pack := range packs {
			if quantityRemaining < 0 {
				break
			}
			if quantityRemaining >= pack.Quantity {
				fulfilment = append(fulfilment, pack)
				quantityRemaining = quantityRemaining - pack.Quantity
				if quantityRemaining >= pack.Quantity {
					break
				}
			}
			if quantityRemaining < pack.Quantity && quantityRemaining > 0 && i == len(packs)-1 {
				fulfilment = append(fulfilment, pack)
				quantityRemaining = quantityRemaining - pack.Quantity
				break
			}
		}
	}

	sum := 0
	for _, pack := range fulfilment {
		sum += pack.Quantity
	}
	for _, pack := range packs {
		if sum == pack.Quantity {
			fulfilment = []model.Pack{pack}
		}
	}

	return model.Order{
		Items: items,
		Packs: fulfilment,
	}
}
