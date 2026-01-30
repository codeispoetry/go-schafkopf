package main

import (
	"net/http"
	"encoding/json"
	"math/rand"
	"fmt"
	"sort"
)


type Card struct {
	Id   int
	Suit  string
	Rank  string
	Value int
	Player int
	Place string
}

type Player struct {
	Id   int
	Name string
}

type Info struct {
	Hand []Card
	Table []Card
	NextPlayer int
	Players []Player
	PlayableCards []Card
}

type Game struct {
	Table []Card
	NextPlayer int
	TrumpSuit string
}

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
		{Id: 0, Name: "Alice"},
		{Id: 1, Name: "Bob"},
		{Id: 2, Name: "Charlie"},
		{Id: 3, Name: "Diana"},
	}

	game.NextPlayer = 0
	game.TrumpSuit = "Herz"

	shuffleDeck()
	dealCards()

	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/trick", trickHandler)
	http.HandleFunc("/finish", finishHandler)
	http.ListenAndServe(":9010", nil)
}

func dealCards() {
	for i := 0; i < 32; i++ {
		Deck[i].Player = i%4
		Deck[i].Place = "Hand"
	}
}

func shuffleDeck() {
	for i := len(Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		Deck[i], Deck[j] = Deck[j], Deck[i]
	}
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Preflight abfangen
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	scores := make(map[int]int)
	for _, player := range players {
		scores[player.Id] = player.getPoints()
	}

	json.NewEncoder(w).Encode(scores)
}

func trickHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	// Preflight abfangen
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		Hand: players[requestBody.Player].Hand(),
		Table: getTable(),
		NextPlayer: game.NextPlayer,
		Players: players,
	}

	json.NewEncoder(w).Encode(Info)
}

func playHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Preflight abfangen
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	if(requestBody.Card > 0 && requestBody.Player == game.NextPlayer &&  len(getTable()) < 4) {
		card := getCardById(requestBody.Card)

		if(card.Player != requestBody.Player || card.Place != "Hand") {
			http.Error(w, "Invalid move", http.StatusBadRequest)
			return
		}

		card.Player = -1
		card.Place = "Table"

		game.NextPlayer = (requestBody.Player + 1) % 4

		fmt.Println("Player", requestBody.Player, "plays card", requestBody.Card)
	} 

	if len(getTable()) == 4 {
		game.NextPlayer = -1
	}

	Info := Info{
		Hand: players[requestBody.Player].Hand(),
		Table: getTable(),
		NextPlayer: game.NextPlayer,
		PlayableCards: PlayableCards(),
		Players: players,
	}

	json.NewEncoder(w).Encode(Info)
}

func getPlayerNames() map[int]string {
	info := make(map[int]string)
	for _, player := range players {
		info[player.Id] = player.Name
	}
	return info
}

func getCardById(id int) *Card {
	for i := range Deck {
		if Deck[i].Id == id {
			return &Deck[i]
		}
	}
	return nil
}

func getTable() []Card {
	var table []Card
	for _, card := range Deck {
		if card.Place == "Table" {
			table = append(table, card)
		}
	}
	return table
}

func (p Player) Hand() []Card {
	var hand []Card
	for _, card := range Deck {
		if card.Player == p.Id && card.Place == "Hand" {
			hand = append(hand, card)
		}
	}

	// Sort cards by trump status and then by suit/rank order
	sort.Slice(hand, func(i, j int) bool {
		cardI, cardJ := hand[i], hand[j]
		
		// Trump cards come first
		if cardI.isTrump() && !cardJ.isTrump() {
			return true
		}
		if !cardI.isTrump() && cardJ.isTrump() {
			return false
		}
		
		// Both trump or both non-trump - sort by suit then rank
		if cardI.Suit != cardJ.Suit {
			return cardI.Suit < cardJ.Suit
		}
		
		rankOrder := map[string]int{
			"Ober": 0, "Unter": 1,  "As": 2, 
			"10": 3, "König": 4, "9": 5, "8": 6, "7": 7,
		}
		return rankOrder[cardI.Rank] < rankOrder[cardJ.Rank]
	})
	return hand
}

func (p Player) getTrick() bool{
	tableCards := getTable()
	if len(tableCards) < 4 {
		return false
	}
		
	for i := range tableCards {
		card := getCardById(tableCards[i].Id)
		card.Player = p.Id
		card.Place = "Trick"
	}

	game.NextPlayer = p.Id

	fmt.Println("Player", p.Id, "takes the trick")
	return true
}

func (p Player) getPoints() int {
	points := 0
	for _, card := range Deck {
		if card.Player == p.Id && card.Place == "Trick" {
			points += card.Value
		}
	}
	return points
}

func PlayableCards() []Card {
	if(game.NextPlayer == -1) {
		return []Card{}
	}

	var allowed []Card
	tableCards := getTable()
	if len(tableCards) == 0 {
		return players[game.NextPlayer].Hand()
	}

	leadCard := tableCards[0]

	for _, card := range players[game.NextPlayer].Hand() {
		if leadCard.isTrump() {
			if card.isTrump() {
				allowed = append(allowed, card)
			}
		} else {
			if card.Suit == leadCard.Suit && !card.isTrump() {
				allowed = append(allowed, card)
			}
		}
	}

	if len(allowed) == 0 {
		return players[game.NextPlayer].Hand()
	}

	return allowed
}

func (c Card) isTrump() bool {
	return c.Suit == game.TrumpSuit || c.Rank == "Ober" || c.Rank == "Unter"
}