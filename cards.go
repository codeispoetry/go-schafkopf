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
	Playable   bool   // whether the card is currently playable
}


func (c Card) isTrump() bool {
	return c.Suit == game.TrumpSuit || c.Rank == "Ober" || c.Rank == "Unter"
}

func (c Card) isPlayable() bool {
	tableCards := getTable()
	if len(tableCards) == 0 {
		return true // I am the lead player and can play any card
	}
	
	leadCard := tableCards[0]
	
	if( leadCard.isTrump() ) {
		if( !players[game.NextPlayer].HasTrump ) {
			return true
		} else{
			if( c.isTrump() ) {
				return true
			}
		}
	}else{
		leadSuit := leadCard.Suit
		
		if( !players[game.NextPlayer].HasSuit ) {
			return true
		}else{
			if( c.Suit == leadSuit && !c.isTrump() ) {
				return true
			}
		}
	}
	
		
	return false;
}


func getCardById(id int) *Card {
	for i := range Deck {
		if Deck[i].Id == id {
			return &Deck[i]
		}
	}
	return nil
}
