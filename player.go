package main

import (
	"fmt"
	"sort"
)

type Player struct {
	Id    int
	Name  string
	Score int
}

func (p *Player) Hand() []Card {
	var hand []Card
	for _, card := range Deck {
		if card.Player == p.Id && card.Place == "Hand" {
			hand = append(hand, card)
		}
	}

	// Sort cards
	sort.Slice(hand, func(i, j int) bool {
		return (hand[i].SortOrder > hand[j].SortOrder)
	})
	return hand
}

func (p *Player) getTrick() bool {
	tableCards := getTable()
	if len(tableCards) < 4 {
		return false
	}

	for i := range tableCards {
		card := getCardById(tableCards[i].Id)
		card.Player = p.Id
		card.Place = "Trick"

		p.Score += card.Value
		fmt.Println("Punkte von", p.Name, ":", p.Score)
	}

	game.NextPlayer = p.Id

	return true
}
