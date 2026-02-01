package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

type Info struct {
	Hand          []*Card
	Table         []Card
	NextPlayer    int
	TrickWinner   int
	Players       []*Player

}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	if (!prepareResponse(w, r, http.MethodPost)) {
		return
	}

	var requestBody struct {
		Player int `json:"player"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if(len(getTable())) == 1 {
		players[requestBody.Player].setHasTrump()
		players[requestBody.Player].setHasSuit(getTable()[0].Suit)
	}
	
	Info := Info{
		Hand:          players[requestBody.Player].Hand(),
		Table:         getTable(),
		NextPlayer:    getNextPlayer(),
		TrickWinner:   getTrickWinner(),
		Players:       players,
	}

	json.NewEncoder(w).Encode(Info)
}

func getTable() []Card {
	var table []Card
	for _, card := range Deck {
		if card.Place == "Table" {
			table = append(table, *card)
		}
	}

	// Sort table cards by player ID
	sort.Slice(table, func(i, j int) bool {
		return table[i].Position < table[j].Position
	})
	return table
}

func getNextPlayer() int {
	for _, player := range players {
		if player.IsNext {
			return player.Id
		}
	}
	return -1
}

func getTrickWinner() int {
	if(len(getTable()) < 4) {
		return -1
	}
	
	leadCard := getTable()[0]
	winnerCard := leadCard
	for _, card := range getTable()[1:] {
		if(card.Trump && !winnerCard.Trump) || (card.Trump == winnerCard.Trump && card.FightOrder > winnerCard.FightOrder) {
			winnerCard = card
		}
	}
	return 	winnerCard.Player
}
