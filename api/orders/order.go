package orders

import (
	"jordanfinners/api/model"
	"log"
	"math"
	"sort"
)

// exactPackMatching calculates if a pack exactly matches the ordered items, either with one pack or a multiple
func exactPackMatching(items int, packs []model.Pack) (bool, []model.Pack) {
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

// multiplesPackMatching calculates if the ordered items can be made up of multiples of a pack
// with less remainder than the smallest pack size, which is the most efficient packs
func multiplesPackMatching(items int, packs []model.Pack) (bool, []model.Pack) {
	var fulfilment []model.Pack
	isExactMultiplePacks := false

	minPack := packs[len(packs)-1]

	for _, pack := range packs {

		itemsOrdered := float64(items)
		quantity := float64(pack.Quantity)
		packsIntoItems := math.Ceil(itemsOrdered / quantity)
		// if the left over quantity after adding all the packs is less than the smallest pack size
		// use that as there is less waste
		log.Printf("packsIntoItems %v left over %v, min %v", packsIntoItems, (packsIntoItems*quantity)-itemsOrdered, float64(minPack.Quantity))
		if (packsIntoItems*quantity)-itemsOrdered < float64(minPack.Quantity) {
			for j := 1; j <= int(packsIntoItems); j++ {
				fulfilment = append(fulfilment, pack)
			}
			isExactMultiplePacks = true
			break
		}
	}
	return isExactMultiplePacks, fulfilment
}

// sizePackMatching calculates how we can make up the ordered items with packs of various sizes
func sizePackMatching(items int, packs []model.Pack) []model.Pack {
	var fulfilment []model.Pack

	quantityRemaining := items
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
	return fulfilment
}

// combinePacks calculates if we can combine packs to make the dispatched fulfilments more efficient by combinding together packs into larger packs
func combinePacks(fulfilment []model.Pack, packs []model.Pack) []model.Pack {

	maxPack := packs[0]

	// we can't combine the largest packs down further for lets keep all them to the side
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
		return lesserPacks[i].Quantity < lesserPacks[j].Quantity
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
	return combinedPacks
}

// CalculateOrder calculates the packs required to fulfil the ordered items
func CalculateOrder(items int, packs []model.Pack) model.Order {

	isExactPacks, fulfilment := exactPackMatching(items, packs)

	if isExactPacks {
		return model.Order{
			Items: items,
			Packs: fulfilment,
		}
	}

	isExactMultiplePacks, fulfilment := multiplesPackMatching(items, packs)
	log.Printf("isExactMultiplePacks %v", isExactMultiplePacks)
	if !isExactMultiplePacks {
		fulfilment = sizePackMatching(items, packs)
	}
	log.Printf("ful %v", fulfilment)

	combinedPacks := combinePacks(fulfilment, packs)

	return model.Order{
		Items: items,
		Packs: combinedPacks,
	}
}
