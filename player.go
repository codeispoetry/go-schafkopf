package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"log"
)

type Player struct {
	Id       int
	Name     string // player's name
	Points    int // points in this round
	Tricks   int // number of tricks won
	IsNext   bool
	Gamer    bool // is gamer or non-gamer
}

func defineHandler(w http.ResponseWriter, r *http.Request) {
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

	for _, player := range players {
		player.Gamer = false
	}
	players[requestBody.Player].Gamer = true
	players[2].Gamer = true

	
	w.WriteHeader(http.StatusOK)
	pingAllClients()
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


	players[requestBody.Player].passOnToTheNext()

	log.Printf("Player %d played card %d", requestBody.Player, requestBody.Card)
	w.WriteHeader(http.StatusOK)
	pingAllClients()
}

func (p *Player) passOnToTheNext() {
	p.IsNext = false;
	if(len(getTable()) < 4) {
		playerId := (p.Id + 1) % len(players)
		players[playerId].IsNext = true
	}
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
	p.Gamer = false
}

