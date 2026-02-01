package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"fmt"
)

type Player struct {
	Id       int
	Name     string
	Score    int
	HasTrump bool
	HasSuit  bool
	IsNext   bool
}

func trickHandler(w http.ResponseWriter, r *http.Request) {
	if !prepareResponse(w, r, http.MethodPost) {
		return
	}

	var requestBody struct {
		Player int `json:"player"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	players[requestBody.Player].getTrick()
	
	w.WriteHeader(http.StatusOK)
	pingAllClients()
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if !prepareResponse(w, r, http.MethodPost) {
		return
	}

	var requestBody struct {
		Player int `json:"player"`
		Card   int `json:"card"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	card := getCardById(requestBody.Card)
	
	card.Place = "Table"
	card.Position = len(getTable())


	fmt.Println(players[requestBody.Player].Name, "spielt", card.Suit, card.Rank)

	players[requestBody.Player].IsNext = false

	if(len(getTable()) < 4) {
		playerId := (requestBody.Player + 1) % len(players)
		players[playerId].IsNext = true
	}


	w.WriteHeader(http.StatusOK)
	pingAllClients()
}

func (p *Player) Hand() []*Card {
	
	var hand []*Card
	for i := range Deck {
		if Deck[i].Player == p.Id && Deck[i].Place == "Hand" {
			Deck[i].Playable = true
			hand = append(hand, Deck[i])
		}
	}

	sort.Slice(hand, func(i, j int) bool {
		return hand[i].SortOrder > hand[j].SortOrder
	})

	return hand
}

func (p *Player) getTrick()  {
	
	tableCards := getTable()
	if len(tableCards) < 4 {
		return
	}

	for _, card := range tableCards {
		card := getCardById(card.Id)
		card.Place = "Trick"
		card.Position = -1
		card.Playable = false
		card.Player = p.Id

		p.Score += card.Value
	}

	p.IsNext = true

}

func getPlayerNames() map[int]string {
	info := make(map[int]string)
	for _, player := range players {
		info[player.Id] = player.Name
	}
	return info
}


