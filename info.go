package main

type Info struct {
	Hand          []Card
	Table         []Card
	NextPlayer    int
	Players       []*Player
	PlayableCards []Card
	TrickWinner   int
	Scores        map[int]int
}

