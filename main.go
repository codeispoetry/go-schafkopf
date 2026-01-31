package main

import (
	"net/http"
)

var Deck = []Card{
	{Id: 32, Suit: "Eichel", Rank: "Ober", Value: 3, Player: 0, Place: "Deck", FightOrder: 33, SortOrder: 73},
	{Id: 1, Suit: "Eichel", Rank: "Unter", Value: 2, Player: 0, Place: "Deck", FightOrder: 23, SortOrder: 63},

	{Id: 3, Suit: "Eichel", Rank: "As", Value: 11, Player: 0, Place: "Deck", FightOrder: 6, SortOrder: 45},
	{Id: 4, Suit: "Eichel", Rank: "10", Value: 10, Player: 0, Place: "Deck", FightOrder: 5, SortOrder: 44},
	{Id: 2, Suit: "Eichel", Rank: "König", Value: 4, Player: 0, Place: "Deck", FightOrder: 4, SortOrder: 43},
	{Id: 5, Suit: "Eichel", Rank: "9", Value: 0, Player: 0, Place: "Deck", FightOrder: 3, SortOrder: 42},
	{Id: 6, Suit: "Eichel", Rank: "8", Value: 0, Player: 0, Place: "Deck", FightOrder: 2, SortOrder: 41},
	{Id: 7, Suit: "Eichel", Rank: "7", Value: 0, Player: 0, Place: "Deck", FightOrder: 1, SortOrder: 40},

	{Id: 8, Suit: "Gras", Rank: "Ober", Value: 3, Player: 0, Place: "Deck", FightOrder: 32, SortOrder: 72},
	{Id: 9, Suit: "Gras", Rank: "Unter", Value: 2, Player: 0, Place: "Deck", FightOrder: 22, SortOrder: 62},

	{Id: 11, Suit: "Gras", Rank: "As", Value: 11, Player: 0, Place: "Deck", FightOrder: 6, SortOrder: 35},
	{Id: 12, Suit: "Gras", Rank: "10", Value: 10, Player: 0, Place: "Deck", FightOrder: 5, SortOrder: 34},
	{Id: 10, Suit: "Gras", Rank: "König", Value: 4, Player: 0, Place: "Deck", FightOrder: 4, SortOrder: 33},
	{Id: 13, Suit: "Gras", Rank: "9", Value: 0, Player: 0, Place: "Deck", FightOrder: 3, SortOrder: 32},
	{Id: 14, Suit: "Gras", Rank: "8", Value: 0, Player: 0, Place: "Deck", FightOrder: 2, SortOrder: 31},
	{Id: 15, Suit: "Gras", Rank: "7", Value: 0, Player: 0, Place: "Deck", FightOrder: 1, SortOrder: 30},

	{Id: 16, Suit: "Herz", Rank: "Ober", Value: 3, Player: 0, Place: "Deck", FightOrder: 31, SortOrder: 71},
	{Id: 17, Suit: "Herz", Rank: "Unter", Value: 2, Player: 0, Place: "Deck", FightOrder: 21, SortOrder: 61},

	{Id: 19, Suit: "Herz", Rank: "As", Value: 11, Player: 0, Place: "Deck", FightOrder: 16, SortOrder: 55},
	{Id: 20, Suit: "Herz", Rank: "10", Value: 10, Player: 0, Place: "Deck", FightOrder: 15, SortOrder: 54},
	{Id: 18, Suit: "Herz", Rank: "König", Value: 4, Player: 0, Place: "Deck", FightOrder: 14, SortOrder: 53},
	{Id: 21, Suit: "Herz", Rank: "9", Value: 0, Player: 0, Place: "Deck", FightOrder: 13, SortOrder: 52},
	{Id: 22, Suit: "Herz", Rank: "8", Value: 0, Player: 0, Place: "Deck", FightOrder: 12, SortOrder: 51},
	{Id: 23, Suit: "Herz", Rank: "7", Value: 0, Player: 0, Place: "Deck", FightOrder: 11, SortOrder: 50},


	{Id: 24, Suit: "Schellen", Rank: "Ober", Value: 3, Player: 0, Place: "Deck", FightOrder: 30, SortOrder: 70},
	{Id: 25, Suit: "Schellen", Rank: "Unter", Value: 2, Player: 0, Place: "Deck", FightOrder: 20, SortOrder: 60},

	{Id: 26, Suit: "Schellen", Rank: "As", Value: 11, Player: 0, Place: "Deck", FightOrder: 6, SortOrder: 15},
	{Id: 27, Suit: "Schellen", Rank: "10", Value: 10, Player: 0, Place: "Deck", FightOrder: 5, SortOrder: 14},
	{Id: 28, Suit: "Schellen", Rank: "König", Value: 4, Player: 0, Place: "Deck", FightOrder: 4, SortOrder: 13},
	{Id: 29, Suit: "Schellen", Rank: "9", Value: 0, Player: 0, Place: "Deck", FightOrder: 3, SortOrder: 12},
	{Id: 30, Suit: "Schellen", Rank: "8", Value: 0, Player: 0, Place: "Deck", FightOrder: 2, SortOrder: 11},
	{Id: 31, Suit: "Schellen", Rank: "7", Value: 0, Player: 0, Place: "Deck", FightOrder: 1, SortOrder: 10},
}

var players []*Player
var game Game

type Game struct {
	Table      []Card
	NextPlayer int
	TrumpSuit  string
}


func main() {
	players = []*Player{
		&Player{Id: 0, Name: "Tom", Score: 0},
		&Player{Id: 1, Name: "Max", Score: 0},
		&Player{Id: 2, Name: "Sibylle", Score: 0},
		&Player{Id: 3, Name: "Birgit", Score: 0},
	}

	game.NextPlayer = 0
	game.TrumpSuit = "Herz"

	shuffleDeck()
	dealCards()

	http.HandleFunc("/ws", handleWSClient)

	http.HandleFunc("/render", renderHandler)

	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/trick", trickHandler)
	http.ListenAndServe(":9010", nil)


}
