package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
		PlayableCards: PlayableCards(),
		Players:       players,
	}

	json.NewEncoder(w).Encode(Info)
}


func finishHandler(w http.ResponseWriter, r *http.Request) {
	if (!prepareResponse(w, r, http.MethodGet)) {
		return
	}

	scores := make(map[int]int)
	for _, player := range players {
		scores[player.Id] = player.getPoints()
	}

	json.NewEncoder(w).Encode(scores)
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

	Info := Info{
		Hand:       players[requestBody.Player].Hand(),
		Table:      getTable(),
		NextPlayer: game.NextPlayer,
		Players:    players,
	}

	json.NewEncoder(w).Encode(Info)
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

	card.Player = -1
	card.Place = "Table"

	game.NextPlayer = (requestBody.Player + 1) % 4

	fmt.Println(players[requestBody.Player].Name, "spielt", card.Suit, card.Rank)
	

	if len(getTable()) == 4 {
		game.NextPlayer = -1
	}

}

func prepareResponse(w http.ResponseWriter, r *http.Request, method string) bool{
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Preflight abfangen
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return false
	}

	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}

	return true
}
