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
	IsFinished	  bool
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

	Hand := setPlayableCards(players[requestBody.Player].Hand())	
	
	Info := Info{
		Hand:          Hand,
		Table:         getTable(),
		NextPlayer:    getNextPlayer(),
		TrickWinner:   getTrickWinner(),
		IsFinished:    isFinished(),
		Players:       players,
	}

	json.NewEncoder(w).Encode(Info)
}

func isFinished() bool {
	playedCards := 0
	for _, card := range Deck {
		if card.Place == "Trick" {
			playedCards++
		}
	}

	return playedCards == len(Deck)
}

func setPlayableCards(hand []*Card) []*Card {
	table := getTable()

	// If no cards on table, all cards are playable
	if len(table) == 0 {
		for _, card := range hand {
			card.Playable = true
		}
		return hand
	}

	leadCard := table[0]
	if leadCard.Trump {
		// Lead card is trump, player must play trump if possible
		hasTrump := false
		for _, card := range hand {
			if card.Trump {
				hasTrump = true
				break
			}
		}
		if hasTrump {
			for _, card := range hand {
				if card.Trump {
					card.Playable = true
				}
			}
			return hand
		}
	}else{
		// Lead card is non-trump, player must follow suit if possible
		leadSuit := leadCard.Suit
		hasSuit := false
		for _, card := range hand {
			if !card.Trump && card.Suit == leadSuit {
				hasSuit = true
				break
			}
		}
		if hasSuit {
			for _, card := range hand {
				if !card.Trump && card.Suit == leadSuit {
					card.Playable = true
				}
			}
			return hand
		}
	}



	// If code comes here, no playable card found, so all are playable
	for _, card := range hand {
		card.Playable = true
	}
	return hand
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
