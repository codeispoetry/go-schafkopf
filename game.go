package main

import (
	"math/rand"
)

func dealCards() {
	for i := 0; i < 32; i++ {
		Deck[i].Player = i % 4
		Deck[i].Place = "Hand"
	}
}

func shuffleDeck() {
	for i := len(Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		Deck[i], Deck[j] = Deck[j], Deck[i]
	}
}

func getHighestCard(cards []Card) Card {
	highest := cards[0]
	for _, card := range cards[1:] {
		if card.isTrump() {
			if !highest.isTrump() || card.Value > highest.Value {
				highest = card
			}
		} else {
			if !highest.isTrump() && card.Suit == highest.Suit && card.Value > highest.Value {
				highest = card
			}
		}
	}
	return highest
}

