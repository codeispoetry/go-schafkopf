package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var url = "http://localhost:9010/"

func main() {

	// Connect to websocket
	wsURL := "ws://localhost:9010/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Listen for messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			break
		}
		fmt.Println(string(message))

		for i := 0; i <= 3; i++ {

			cardToPlay, trickWinner, nextPlayer := getInfo(i)

			fmt.Println(cardToPlay, trickWinner, nextPlayer)
			if(trickWinner != -1) {
				takeTrick(trickWinner)
				continue
			}

			if(nextPlayer != i || nextPlayer == -1) {
				continue
			}

			fmt.Printf("Playing for %d\n", i)
			time.Sleep(10 * time.Millisecond)
			if cardToPlay != -1 {
				play(i, cardToPlay)
			}

		}
	}
}

func play(player int, cardId int) {
	if cardId == -1 {
		fmt.Println("No playable card found")
		return
	}

	payload := map[string]int{"player": player, "card": cardId}
	jsonData, _ := json.Marshal(payload)

	route := "play"
	resp, err := http.Post(url+route, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func takeTrick(player int) {
	payload := map[string]int{"player": player}
	jsonData, _ := json.Marshal(payload)

	route := "trick"
	resp, err := http.Post(url+route, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Player %d took the trick\n", player)
}

func getInfo(player int) (int, int, int) {
	payload := map[string]int{"player": player}
	jsonData, _ := json.Marshal(payload)

	route := "render"
	resp, err := http.Post(url+route, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return -1, -1,-1
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		fmt.Printf("JSON decode error: %v\n", err)
		return -1, -1,-1
	}

	cardToPlay := -1
	if hand, ok := responseData["Hand"].([]interface{}); ok {
		for _, c := range hand {
			card, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			playable, _ := card["Playable"].(bool)
			if !playable {
				continue
			}
			if id, ok := card["Id"].(float64); ok {
				cardToPlay = int(id)
			}
		}
	}

	trickWinner := -1
	if tw, ok := responseData["TrickWinner"].(float64); ok {
		trickWinner = int(tw)
	}

	return cardToPlay, trickWinner, int(responseData["NextPlayer"].(float64))

}
