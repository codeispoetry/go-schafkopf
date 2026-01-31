package main

type Card struct {
	Id         int
	Suit       string // e.g., "Eichel", "Gras", "Herz", "Schellen"
	Rank       string // e.g., "Ober", "Unter", "König", "As", "10", "9", "8", "7"
	Value      int    // point value of the card, 0, 2, 3, 4, 10, or 11
	FightOrder int    // order in which cards win tricks
	SortOrder  int    // order for sorting in hand
	Player     int    // ID of the player who currently holds the card, or has won it
	Place      string // "Deck", "Hand", "Table", "Trick"
	Position   int    // position on the table when played
}




func (c Card) isTrump() bool {
	return c.Suit == game.TrumpSuit || c.Rank == "Ober" || c.Rank == "Unter"
}

func PlayableCards() []Card {
	if game.NextPlayer == -1 {
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
