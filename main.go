package main

import (
	"net/http"
)

var Deck = []Card{
	{Id: 32, Suit: "Eichel", Rank: "Ober", Value: 3, Player: 0, Place: "Deck"},
	{Id: 1, Suit: "Eichel", Rank: "Unter", Value: 2, Player: 0, Place: "Deck"},
	{Id: 2, Suit: "Eichel", Rank: "König", Value: 4, Player: 0, Place: "Deck"},
	{Id: 3, Suit: "Eichel", Rank: "As", Value: 11, Player: 0, Place: "Deck"},
	{Id: 4, Suit: "Eichel", Rank: "10", Value: 10, Player: 0, Place: "Deck"},
	{Id: 5, Suit: "Eichel", Rank: "9", Value: 0, Player: 0, Place: "Deck"},
	{Id: 6, Suit: "Eichel", Rank: "8", Value: 0, Player: 0, Place: "Deck"},
	{Id: 7, Suit: "Eichel", Rank: "7", Value: 0, Player: 0, Place: "Deck"},
	{Id: 8, Suit: "Gras", Rank: "Ober", Value: 3, Player: 0, Place: "Deck"},
	{Id: 9, Suit: "Gras", Rank: "Unter", Value: 2, Player: 0, Place: "Deck"},
	{Id: 10, Suit: "Gras", Rank: "König", Value: 4, Player: 0, Place: "Deck"},
	{Id: 11, Suit: "Gras", Rank: "As", Value: 11, Player: 0, Place: "Deck"},
	{Id: 12, Suit: "Gras", Rank: "10", Value: 10, Player: 0, Place: "Deck"},
	{Id: 13, Suit: "Gras", Rank: "9", Value: 0, Player: 0, Place: "Deck"},
	{Id: 14, Suit: "Gras", Rank: "8", Value: 0, Player: 0, Place: "Deck"},
	{Id: 15, Suit: "Gras", Rank: "7", Value: 0, Player: 0, Place: "Deck"},
	{Id: 16, Suit: "Herz", Rank: "Ober", Value: 3, Player: 0, Place: "Deck"},
	{Id: 17, Suit: "Herz", Rank: "Unter", Value: 2, Player: 0, Place: "Deck"},
	{Id: 18, Suit: "Herz", Rank: "König", Value: 4, Player: 0, Place: "Deck"},
	{Id: 19, Suit: "Herz", Rank: "As", Value: 11, Player: 0, Place: "Deck"},
	{Id: 20, Suit: "Herz", Rank: "10", Value: 10, Player: 0, Place: "Deck"},
	{Id: 21, Suit: "Herz", Rank: "9", Value: 0, Player: 0, Place: "Deck"},
	{Id: 22, Suit: "Herz", Rank: "8", Value: 0, Player: 0, Place: "Deck"},
	{Id: 23, Suit: "Herz", Rank: "7", Value: 0, Player: 0, Place: "Deck"},
	{Id: 24, Suit: "Schellen", Rank: "Ober", Value: 3, Player: 0, Place: "Deck"},
	{Id: 25, Suit: "Schellen", Rank: "Unter", Value: 2, Player: 0, Place: "Deck"},
	{Id: 26, Suit: "Schellen", Rank: "König", Value: 4, Player: 0, Place: "Deck"},
	{Id: 27, Suit: "Schellen", Rank: "As", Value: 11, Player: 0, Place: "Deck"},
	{Id: 28, Suit: "Schellen", Rank: "10", Value: 10, Player: 0, Place: "Deck"},
	{Id: 29, Suit: "Schellen", Rank: "9", Value: 0, Player: 0, Place: "Deck"},
	{Id: 30, Suit: "Schellen", Rank: "8", Value: 0, Player: 0, Place: "Deck"},
	{Id: 31, Suit: "Schellen", Rank: "7", Value: 0, Player: 0, Place: "Deck"},
}

var players []Player
var game Game

func main() {
	players = []Player{
		{Id: 0, Name: "Tom"},
		{Id: 1, Name: "Max"},
		{Id: 2, Name: "Sibylle"},
		{Id: 3, Name: "Birgit"},
	}

	game.NextPlayer = 0
	game.TrumpSuit = "Herz"

	shuffleDeck()
	dealCards()

	http.HandleFunc("/ws", handleWSClient)



	http.HandleFunc("/render", renderHandler)

	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/trick", trickHandler)
	http.HandleFunc("/finish", finishHandler)
	http.ListenAndServe(":9010", nil)
}
