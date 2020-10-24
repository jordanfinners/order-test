package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

func getPacks() ([]Pack, error) {
	response, err := http.Get("https://raw.githubusercontent.com/jordanfinners/order-test/master/data/packs.json")
	if err != nil {
		log.Printf("Failed to load pack data: %v", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var packs []Pack
	err = json.Unmarshal(body, &packs)
	if err != nil {
		log.Printf("Failed unmarshal pack data: %v", err)
		return nil, err
	}

	sort.SliceStable(packs, func(i, j int) bool {
		return packs[i].Quantity > packs[j].Quantity
	})
	return packs, nil
}
