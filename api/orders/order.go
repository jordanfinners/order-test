package orders

import (
	"jordanfinners/api/model"
	"math"
)

// isSmallerThanMinPackSize calculates if a pack is smaller than the min pack size available
func isSmallerThanMinPackSize(items int, packs []model.Pack) (bool, []model.Pack) {
	var fulfilment []model.Pack
	isSmaller := false

	minPack := packs[len(packs)-1]

	if items < minPack.Quantity {
		fulfilment = append(fulfilment, minPack)
		isSmaller = true
	}
	return isSmaller, fulfilment
}

// isExactPackMatching calculates if a pack exactly matches the ordered items, either with one pack or a multiple
func isExactPackMatching(items int, packs []model.Pack) (bool, []model.Pack) {
	var fulfilment []model.Pack
	isExactPacks := false

	for _, pack := range packs {
		// if the packs equal just use it and break out
		if items == pack.Quantity {
			fulfilment = append(fulfilment, pack)
			isExactPacks = true
			break
		}
		// if you can make up the pack with exact multiples add those and escape
		if math.Remainder(float64(items), float64(pack.Quantity)) == 0.0 {
			packsToAdd := items / pack.Quantity
			for j := 1; j <= packsToAdd; j++ {
				fulfilment = append(fulfilment, pack)
			}
			isExactPacks = true
			break
		}
	}
	return isExactPacks, fulfilment
}

// packMultiples calculates how many of each pack to add fulfil an order based on multiples of packs
func packMultiples(items int, packs []model.Pack) []model.Pack {
	var fulfilment []model.Pack
	minPack := packs[len(packs)-1]
	quantityRemaining := items

	for _, pack := range packs {
		// calculate how many packs fit into the order, and what the amount left over is
		packsIntoItems := quantityRemaining / pack.Quantity
		remainder := quantityRemaining - (packsIntoItems * pack.Quantity)
		// if less than one pack fulfils the order and the amount extra if we used the pack is less than the smallest pack size
		// use that pack as its the least wasteful way to fulfil the order
		if packsIntoItems < 1 && (pack.Quantity-remainder) < minPack.Quantity {
			fulfilment = append(fulfilment, pack)
			quantityRemaining = 0
			break
		}
		// if more than a single pack fulfils the order then add all the packs needed
		if packsIntoItems > 1 {
			for i := 1; i <= packsIntoItems; i++ {
				fulfilment = append(fulfilment, pack)
			}
			quantityRemaining = remainder
		}
		// if one pack fills the order and there might be items left over, but not over fulfilled then use the pack
		if packsIntoItems == 1 && remainder >= 0 {
			fulfilment = append(fulfilment, pack)
			quantityRemaining = remainder
		}
	}
	// if there is any left over fill with min packs
	if quantityRemaining > 0 {
		packsIntoItems := quantityRemaining / minPack.Quantity
		for i := 1; i <= packsIntoItems; i++ {
			fulfilment = append(fulfilment, minPack)
		}
	}
	return fulfilment
}

// CalculateOrder calculates the packs required to fulfil the ordered items
func CalculateOrder(items int, packs []model.Pack) model.Order {

	isSmaller, fulfilment := isSmallerThanMinPackSize(items, packs)

	if isSmaller {
		return model.Order{
			Items: items,
			Packs: fulfilment,
		}
	}

	isExactPacks, fulfilment := isExactPackMatching(items, packs)

	if isExactPacks {
		return model.Order{
			Items: items,
			Packs: fulfilment,
		}
	}

	fulfilment = packMultiples(items, packs)

	return model.Order{
		Items: items,
		Packs: fulfilment,
	}
}
