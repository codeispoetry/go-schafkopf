package main

import (
	"fmt"
	"sort"
)

type Card struct {
	Id     int
	Suit   string
	Rank   string
	Value  int
	Player int
	Place  string
}

type Player struct {
	Id   int
	Name string
}

type Info struct {
	Hand          []Card
	Table         []Card
	NextPlayer    int
	Players       []Player
	PlayableCards []Card
}

type Game struct {
	Table      []Card
	NextPlayer int
	TrumpSuit  string
}

func (p Player) Hand() []Card {
	var hand []Card
	for _, card := range Deck {
		if card.Player == p.Id && card.Place == "Hand" {
			hand = append(hand, card)
		}
	}

	// Sort cards by trump status and then by suit/rank order
	sort.Slice(hand, func(i, j int) bool {
		cardI, cardJ := hand[i], hand[j]

		// Trump cards come first
		if cardI.isTrump() && !cardJ.isTrump() {
			return true
		}
		if !cardI.isTrump() && cardJ.isTrump() {
			return false
		}

		// Both trump or both non-trump - sort by suit then rank
		if cardI.Suit != cardJ.Suit {
			return cardI.Suit < cardJ.Suit
		}

		rankOrder := map[string]int{
			"Ober": 0, "Unter": 1, "As": 2,
			"10": 3, "König": 4, "9": 5, "8": 6, "7": 7,
		}
		return rankOrder[cardI.Rank] < rankOrder[cardJ.Rank]
	})
	return hand
}

func (p Player) getTrick() bool {
	tableCards := getTable()
	if len(tableCards) < 4 {
		return false
	}

	for i := range tableCards {
		card := getCardById(tableCards[i].Id)
		card.Player = p.Id
		card.Place = "Trick"
	}

	game.NextPlayer = p.Id

	fmt.Println("Player", p.Id, "takes the trick")
	return true
}

func (p Player) getPoints() int {
	points := 0
	for _, card := range Deck {
		if card.Player == p.Id && card.Place == "Trick" {
			points += card.Value
		}
	}
	return points
}

func (c Card) isTrump() bool {
	return c.Suit == game.TrumpSuit || c.Rank == "Ober" || c.Rank == "Unter"
}

func PlayableCards() []Card {
	if game.NextPlayer == -1 {
		return []Card{}
	}

	var allowed []Card
	tableCards := getTable()
	if len(tableCards) == 0 {
		return players[game.NextPlayer].Hand()
	}

	leadCard := tableCards[0]

	for _, card := range players[game.NextPlayer].Hand() {
		if leadCard.isTrump() {
			if card.isTrump() {
				allowed = append(allowed, card)
			}
		} else {
			if card.Suit == leadCard.Suit && !card.isTrump() {
				allowed = append(allowed, card)
			}
		}
	}

	if len(allowed) == 0 {
		return players[game.NextPlayer].Hand()
	}

	return allowed
}
