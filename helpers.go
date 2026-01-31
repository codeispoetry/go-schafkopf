package main

import (
	"sort"
)

func getPlayerNames() map[int]string {
	info := make(map[int]string)
	for _, player := range players {
		info[player.Id] = player.Name
	}
	return info
}

func getCardById(id int) *Card {
	for i := range Deck {
		if Deck[i].Id == id {
			return &Deck[i]
		}
	}
	return nil
}

func getTable() []Card {
	var table []Card
	for _, card := range Deck {
		if card.Place == "Table" {
			table = append(table, card)
		}
	}

	// Sort table cards by player ID
	sort.Slice(table, func(i, j int) bool {
		return table[i].Position < table[j].Position
	})
	return table
}

