package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"
)

type GameOption struct {
	Game string
	Suit string
}
type Info struct {
	Hand          []*Card
	Table         []Card
	NextPlayer    int
	TrickWinner   int
	Status 	      string
	Players       []*Player
	Scores        []int
	FinishLine    string
	GameOptions   []GameOption
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
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

	Hand := players[requestBody.Player].Hand()

	if requestBody.Player == getNextPlayer() {
		Hand = setPlayableCards(Hand)
	}

	
	var gameOptions []GameOption
	if !isGameDefined() {
		for _, card := range Hand {
			if card.Suit == "Herz" {
				continue
			}
			if slices.Contains([]string{"Sieben", "Acht", "Neun", "König", "Zehn"}, card.Rank){
				for _, option := range gameOptions {
					if option.Game == "Sauspiel" && option.Suit == card.Suit {
						goto NextCard
					}
				}
				gameOptions = append(gameOptions, GameOption{Game: "Sauspiel", Suit: card.Suit})
			}
		NextCard:
		}
	}

	status := ""
	if(isGameDefined()) {
		status = "defined"
	}
	if isFinished() {
		status = "finished"
	}
	

	Info := Info{
		Hand:          Hand,
		Table:         getTable(),
		NextPlayer:    getNextPlayer(),
		TrickWinner:   getTrickWinner(),
		Status:		   status,
		Players:       players,
		Scores:        calculateScores(),
		FinishLine:    getFinishLine(),
		GameOptions:   gameOptions,
	}

	json.NewEncoder(w).Encode(Info)
}

func getFinishLine() string {
	if !isFinished() || !isGameDefined() {
		return ""
	}

	gamers := []string{}
	gamersScore := 0

	for _, player := range players {
		if player.Gamer {
			gamers = append(gamers, player.Name)
			gamersScore += player.Points
		}
	}

	wonOrLost := map[bool]string{true: "gewonnen", false: "verloren"}[gamersScore >= 61]
	verb := map[bool]string{true: "haben", false: "hat"}[len(gamers) > 1]

	return fmt.Sprintf("%s %s mit %d Punkten %s.", strings.Join(gamers, " und "), verb, gamersScore, wonOrLost)
}

func calculateScores() []int {
	if !isFinished() || !isGameDefined() {
		return []int{0, 0}
	}

	gamerPoints := 0
	nonGamerPoints := 0

	for _, player := range players {
		if player.Gamer {
			gamerPoints += player.Points
		} else {
			nonGamerPoints += player.Points
		}
	}

	return []int{gamerPoints, nonGamerPoints}

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

	// if no gamer, all cards unplayable
	gameDefined := false
	for _, player := range players {
		if player.Gamer {
			gameDefined = true
			break
		}
	}
	if !gameDefined {
		return hand
	}

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
	} else {
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

	// Sort table cards by position
	sort.Slice(table, func(i, j int) bool {
		return table[i].Position < table[j].Position
	})
	return table
}

func getNextPlayer() int {
	if isFinished() {
		for _,card := range Deck {
			if card.Position == 1 {
				return card.Player
			}
		}
	}

	for _, player := range players {
		if player.IsNext {
			return player.Id
		}
	}
	return -1
}

func getTrickWinner() int {
	if len(getTable()) < 4 {
		return -1
	}

	leadCard := getTable()[0]
	winnerCard := leadCard
	for _, card := range getTable()[1:] {
		if (card.Trump && !winnerCard.Trump) || (card.Trump == winnerCard.Trump && card.FightOrder > winnerCard.FightOrder) {
			winnerCard = card
		}
	}
	return winnerCard.Player
}

func isGameDefined() bool {
	for _, player := range players {
		if player.Gamer {
			return true
		}
	}
	return false
}

func countTricks() int {
	counter := 0
	for _, card := range Deck {
		if card.Place == "Trick"{
			counter++
		}
	}
	return counter
}