package main

import (
	"math/rand"
)

type Card struct {
	Id         int
	Suit       string // e.g., "Eichel", "Gras", "Herz", "Schellen"
	Rank       string // e.g., "Ober", "Unter", "König", "As", "10", "9", "8", "7"
	Value      int    // point value of the card, 0, 2, 3, 4, 10, or 11
	FightOrder int    // order in which cards win tricks
	SortOrder  int    // order for sorting in hand
	Player     int    // ID of the player who currently holds the card, or has won it
	Place      string // "Deck", "Hand", "Table", "Trick"
	Position   int    // position on the table when played
	Playable   bool   // whether the card is currently playable
	Trump 	   bool   // whether the card is a trump card
}


func dealCards() {
	rand.Shuffle(len(Deck), func(i, j int) { Deck[i], Deck[j] = Deck[j], Deck[i] })

	for i,_ := range Deck {
		Deck[i].Player = i%4
		Deck[i].Place = "Hand"
	}
}

func getCardById(cardId int) *Card {
	for _, card := range Deck {
		if card.Id == cardId {
			return card
		}
	}
	return nil
}

