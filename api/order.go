package api

// Pack represents packs of items that can be sent to the customer
type Pack struct {
	Quantity int `json:"quantity"`
}

// Order represents a customers order and the pack(s) that fulfil it
type Order struct {
	Items int    `json:"items"`
	Packs []Pack `json:"packs"`
}

func calculateOrder(items int, packs []Pack) Order {
	var fulfilment []Pack
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
	for _, v := range fulfilment {
		sum += v.Quantity
	}
	for _, p := range packs {
		if sum == p.Quantity {
			fulfilment = []Pack{p}
		}
	}

	return Order{
		Items: items,
		Packs: fulfilment,
	}
}
