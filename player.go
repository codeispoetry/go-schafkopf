package main

import (
	"fmt"
	"sort"
)

type Player struct {
	Id    int
	Name  string
	Score int
	HasTrump bool
	HasSuit bool
}

func (p *Player) Hand() []*Card {
	var hand []*Card
	for i := range Deck {
		c := &Deck[i]
		if c.Player == p.Id && c.Place == "Hand" {
			c.Playable = false
			hand = append(hand, c)
		}
	}

	for _, card := range hand {
		if( card.isTrump() ) {
			p.HasTrump = true
		}	
		
		table := getTable()

		if( len(table) > 0 ) {
			leadCard := table[0]
			leadSuit := leadCard.Suit
	
			if( card.Suit == leadSuit ) {
				p.HasSuit = true
			}
		}
	}

	for _, card := range hand {
		card.Playable = p.Id == game.NextPlayer && card.isPlayable()	
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
