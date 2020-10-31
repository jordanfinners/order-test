package orders

import (
	"jordanfinners/api/model"
	"math"
	"sort"
)

// CalculateOrder calculates the packs required to fulfil the ordered items
// It does this by starting with the largest pack size and working down
// If the items ordered after adding a pack size is still greater than that pack size, then it starts again, using the largest pack sizes possible to fill an order
// If the order is smaller than the smallest pack size available, then use the smallest pack size as the packs cannot be broken up
// Once calculated the packs required to fulfil the order, check if we can use a single pack to fulfil the order instead of multiple as this is more efficient for shipping. This has been kept simple for now but in future should be expanded to check all combinations of packs to see if the order can be simplified.
func CalculateOrder(items int, packs []model.Pack) model.Order {
	var fulfilment []model.Pack
	var quantityRemaining = items
	exactPacks := false

	for _, pack := range packs {
		// if the packs equal just use it and break out
		if quantityRemaining == pack.Quantity {
			fulfilment = append(fulfilment, pack)
			quantityRemaining = 0
			exactPacks = true
			break
		}
		// if you can make up the pack with exact multiples add those and escape
		if math.Remainder(float64(quantityRemaining), float64(pack.Quantity)) == 0.0 {
			packsToAdd := quantityRemaining / pack.Quantity
			for j := 1; j <= packsToAdd; j++ {
				fulfilment = append(fulfilment, pack)
			}
			quantityRemaining = 0
			exactPacks = true
			break
		}
	}

	if exactPacks {
		return model.Order{
			Items: items,
			Packs: fulfilment,
		}
	}

	for quantityRemaining > 0 {
		for i, pack := range packs {
			if quantityRemaining < 0 {
				break
			}
			// otherwise build up if packs can fulfil it
			if quantityRemaining >= pack.Quantity {
				fulfilment = append(fulfilment, pack)
				quantityRemaining = quantityRemaining - pack.Quantity
				if quantityRemaining >= pack.Quantity {
					break
				}
			}
			// if there is a small amount remaining, less than the smallest pack and were at the smallest pack just use that
			if quantityRemaining < pack.Quantity && quantityRemaining > 0 && i == len(packs)-1 {
				fulfilment = append(fulfilment, pack)
				quantityRemaining = quantityRemaining - pack.Quantity
				break
			}
		}
	}

	// we can't combine the largest packs down further for lets keep all them to the side
	maxPack := packs[0]
	positionLesserPacks := len(fulfilment)
	for i, pack := range fulfilment {
		if pack.Quantity != maxPack.Quantity {
			positionLesserPacks = i
			break
		}
	}
	combinedPacks := fulfilment[:positionLesserPacks]

	// grab all the smaller packs to start combining
	lesserPacks := fulfilment[positionLesserPacks:]
	sort.SliceStable(lesserPacks, func(i, j int) bool {
		return packs[i].Quantity < packs[j].Quantity
	})

	for i := 0; len(lesserPacks) > i; i++ {
		for {
			fulfilmentPack := lesserPacks[i]
			nextIndex := i + 1

			var sum int
			if i != len(lesserPacks)-1 {
				sum = fulfilmentPack.Quantity + lesserPacks[nextIndex].Quantity

				isModified := false
				for _, pack := range packs {
					if sum == pack.Quantity {
						lesserPacks[nextIndex] = pack
						lesserPacks = lesserPacks[nextIndex:]
						isModified = true
					}
				}
				if !isModified {
					break
				}
			} else {
				break
			}
		}
	}
	combinedPacks = append(combinedPacks, lesserPacks...)

	return model.Order{
		Items: items,
		Packs: combinedPacks,
	}
}
