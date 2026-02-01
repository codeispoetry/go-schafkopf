package main

import (
	"encoding/json"
	"net/http"
	"sort"
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
	card.Playable = false


	players[requestBody.Player].IsNext = false

	if(len(getTable()) < 4) {
		playerId := (requestBody.Player + 1) % len(players)
		players[playerId].IsNext = true
	}


	w.WriteHeader(http.StatusOK)
	pingAllClients()
}

func (p *Player) setHasTrump() {
	for _, card := range Deck {
		if card.Player == p.Id && card.Place == "Hand" && card.Trump {
			p.HasTrump = true
			return
		}
	}
	p.HasTrump = false
}

func (p *Player) setHasSuit(suit string) {
	if getTable()[0].Trump {
		p.HasSuit = false
		return
	}
	for _, card := range Deck {
		if card.Player == p.Id && card.Place == "Hand" && card.Suit == suit && !card.Trump {
			p.HasSuit = true
			return
		}
	}
	p.HasSuit = false
}

func (p *Player) Hand() []*Card {
	var hand []*Card
	for i := range Deck {
		if Deck[i].Player == p.Id && Deck[i].Place == "Hand" {
			Deck[i].Playable = p.IsNext && isPlayable(Deck[i], p)
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

func isPlayable(card *Card, p *Player) bool {
	if(len(getTable()) == 0) {
		return true
	}

	leadCard := getTable()[0]
	
	if leadCard.Trump {
		if( p.HasTrump && card.Trump ) || (!p.HasTrump) {
			return true
		}
	}

	if !leadCard.Trump {
		if( p.HasSuit && card.Suit == leadCard.Suit && !card.Trump) || (!p.HasSuit) {
			return true
		}
	}

	
	
	return false
}

func (p *Player) reset() {
	p.HasTrump = false
	p.HasSuit = false
	p.IsNext = false
}


