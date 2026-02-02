package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"log"
)

type Player struct {
	Id       int
	Name     string
	Points    int
	Tricks   int
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
	card.Playable = false


	players[requestBody.Player].IsNext = false

	if(len(getTable()) < 4) {
		playerId := (requestBody.Player + 1) % len(players)
		players[playerId].IsNext = true
	}

	log.Printf("Player %d played card %d", requestBody.Player, requestBody.Card)
	w.WriteHeader(http.StatusOK)
	pingAllClients()
}



func (p *Player) Hand() []*Card {
	var hand []*Card
	for i := range Deck {
		if Deck[i].Player == p.Id && Deck[i].Place == "Hand" {
			Deck[i].Playable = false
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

		p.Points += card.Value
	}
	
	p.Tricks++
	p.IsNext = true
}

func (p *Player) reset() {
	p.IsNext = false
	p.Points = 0
	p.Tricks = 0
}

