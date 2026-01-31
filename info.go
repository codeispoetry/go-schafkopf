package main

type Info struct {
	Hand          []*Card
	Table         []Card
	NextPlayer    int
	Players       []*Player
	TrickWinner   int
	Scores        map[int]int
}

	