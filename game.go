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

func whoWonTrick() int {
	if(len(getTable()) < 4) {
		return -1
	}
	
	tableCards := getTable()
	highestCard := tableCards[0]
	leadSuit := tableCards[0].Suit
	
	for _, card := range tableCards[1:] {
		if(card.isTrump()) {
			if(card.FightOrder > highestCard.FightOrder) {
				highestCard = card
			}
		 } else {
			if(!highestCard.isTrump() && card.Suit == leadSuit && card.FightOrder > highestCard.FightOrder) {
				highestCard = card
			}
		 }
	}

	
	return highestCard.Player
}
