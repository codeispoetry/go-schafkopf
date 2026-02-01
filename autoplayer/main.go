package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var url = "http://localhost:9010/"

func main() {
	player := 0
	if len(os.Args) > 1 {
		if arg, err := strconv.Atoi(os.Args[1]); err == nil {
			player = arg
		}
	}

	fmt.Print("Autoplayer for player ", player, "\n")

	payload := map[string]int{"player": player}
	jsonData, _ := json.Marshal(payload)
	
	route := "start"
	resp, err := http.Post(url+route, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)
}
