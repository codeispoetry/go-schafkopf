package main

import (
	"encoding/json"
	"net/http"
)

type Info struct {
	Hand          []*Card
	Table         []Card
	NextPlayer    int
	Players       []*Player
	TrickWinner   int
	Scores        map[int]int
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

	if len(getTable()) == 4 {
		game.NextPlayer = -1
	}

	Info := Info{
		Hand:          players[requestBody.Player].Hand(),
		Table:         getTable(),
		NextPlayer:    game.NextPlayer,
		Players:       players,
		TrickWinner:   whoWonTrick(),
	}

	json.NewEncoder(w).Encode(Info)
}

