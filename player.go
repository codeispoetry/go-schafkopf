package main

import (
	"fmt"
	"sort"
	"net/http"
	"encoding/json"
)

type Player struct {
	Id    int
	Name  string
	Score int
	HasTrump bool
	HasSuit bool
}

func trickHandler(w http.ResponseWriter, r *http.Request) {
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

	player := players[requestBody.Player]
	success := player.getTrick()
	if !success {
		http.Error(w, "Not enough cards on table", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	pingAllClients()
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if (!prepareResponse(w, r, http.MethodPost)) {
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

	if requestBody.Card == 0 || requestBody.Player != game.NextPlayer || len(getTable()) == 4 {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}
	
	card := getCardById(requestBody.Card)

	if card.Player != requestBody.Player || card.Place != "Hand" {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	card.Position = len(getTable())
	card.Place = "Table"

	game.NextPlayer = (requestBody.Player + 1) % 4

	fmt.Println(players[requestBody.Player].Name, "spielt", card.Suit, card.Rank)
	

	if len(getTable()) == 4 {

		game.NextPlayer = -1
	}

		
	w.WriteHeader(http.StatusOK)
	pingAllClients()
}

func (p *Player) Hand() []*Card {
	var hand []*Card
	for i := range Deck {
		c := &Deck[i]
		if c.Player == p.Id && c.Place == "Hand" {
			c.Playable = false
			hand = append(hand, c)
		}
	}

	for _, card := range hand {
		if( card.isTrump() ) {
			p.HasTrump = true
		}	
		
		table := getTable()

		if( len(table) > 0 ) {
			leadCard := table[0]
			leadSuit := leadCard.Suit
	
			if( card.Suit == leadSuit ) {
				p.HasSuit = true
			}
		}
	}

	for _, card := range hand {
		card.Playable = p.Id == game.NextPlayer && card.isPlayable()	
	}

	// Sort cards
	sort.Slice(hand, func(i, j int) bool {
		return (hand[i].SortOrder > hand[j].SortOrder)
	})

	return hand
}


func (p *Player) getTrick() bool {
	tableCards := getTable()
	if len(tableCards) < 4 {
		return false
	}

	for i := range tableCards {
		card := getCardById(tableCards[i].Id)
		card.Player = p.Id
		card.Place = "Trick"

		p.Score += card.Value
		fmt.Println("Punkte von", p.Name, ":", p.Score)
	}

	game.NextPlayer = p.Id

	return true
}


func getPlayerNames() map[int]string {
	info := make(map[int]string)
	for _, player := range players {
		info[player.Id] = player.Name
	}
	return info
}
