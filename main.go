package main

import (
	"fmt"
	"math/rand"
	"slices"
)

type player struct {
	hand   []int
	id     int
	name   string
	points int
}

func main() {
	var players []player
	var table [][]int
	// [1] Get players' names and ids
	players = getPlayers()
	nRounds := 11
	if len(players) == 10 {
		nRounds = 10
	}
	// [2] Populate player hands
	players, table = getHands(players)
	fmt.Println("SETUP:")
	for _, x := range players {
		fmt.Println(x.name, x.id, x.hand)
	}
	for _, row := range table {
		fmt.Println(row)
	}
	fmt.Println()
	fmt.Println("========================================================")

	// [3] playing card randomly
	// simulation:
	var played map[int]int
	for round := 0; round < nRounds; round++ {
		// [4] play card
		played, players = playCards(players)

		// [5] process played card
		sortedCards := getSortedCards(played)

		for _, card := range sortedCards {
			// Get nearest card
			nearestPos, isSmallest := getNearest(card, table)
			// fmt.Printf("Nearest: %v\n", nearestPos)
			fmt.Printf("Played: ")
			for _, k := range sortedCards {
				fmt.Printf("%v: %v, ", string(played[k]+65), k)
			}
			fmt.Printf("Action for: %v ", string(played[card]+65))
			if isSmallest {
				// Choose Row
				var choice int
				choice = rand.Intn(4)
				// fmt.Printf("CHOICE IS :%v\n", choice)
				fmt.Printf("Cannot stack, takes row %v!\n", choice)
				players[played[card]].points += getPoints(table[choice])
				table[choice] = []int{card}
			} else {
				fmt.Printf("Stack on %v. ", nearestPos)
				if len(table[nearestPos]) == 5 {
					fmt.Println("Busted!")
					players[played[card]].points += getPoints(table[nearestPos])
					table[nearestPos] = []int{}
				}
				table[nearestPos] = append(table[nearestPos], card)
				fmt.Println()
			}

			fmt.Println()
			for _, row := range table {
				fmt.Println(row)
			}
			fmt.Println()
			fmt.Println("========================================================")
		}

		fmt.Println("AFTER ROUND ", round)
		for _, x := range players {
			fmt.Println(x.name, x.id, x.hand, x.points)
		}
	}
}

func getPoints(row []int) int {
	var total int
	for _, x := range row {
		total += pointsLookUp(x)
	}
	return total
}

func pointsLookUp(x int) int {
	if x == 55 {
		return 7
	}
	return 1
}

func getNearest(card int, table [][]int) (int, bool) {
	// Find closest cards
	min := 1000
	var pos int
	isSmallest := true
	for i, row := range table {
		tail := row[len(row)-1]
		if card < tail {
			continue
		}
		if (card - tail) < min {
			min = card - tail
			pos = i
		}
		isSmallest = false
	}
	return pos, isSmallest
}

func playCards(players []player) (map[int]int, []player) {
	played := make(map[int]int)

	for i := 0; i < len(players); i++ {
		id := players[i].id
		hand := players[i].hand
		choice := rand.Intn(len(hand))
		played[hand[choice]] = id
		players[i].hand = append(hand[:choice], hand[choice+1:]...)
	}
	return played, players
}

func getPlayers() []player {
	// sample players
	numPlayers := 4
	var players []player
	for i := 0; i < numPlayers; i++ {
		players = append(players, player{id: i, name: string(i + 65), points: 0, hand: []int{}})
	}
	return players
}

func getHands(players []player) ([]player, [][]int) {
	fullDeck := getFullDeck()
	handLimit := 11
	if len(players) == 10 {
		handLimit = 10
	}
	start := 0
	for i := 0; i < len(players); i++ {
		dealtHand := fullDeck[start : start+handLimit]
		slices.Sort(dealtHand)
		players[i].hand = dealtHand
		start += handLimit
	}
	var table [][]int
	for i := start; i < start+4; i++ {
		row := []int{fullDeck[i]}
		table = append(table, row)
	}
	return players, table
}

func getFullDeck() []int {
	var deck []int
	for i := 1; i <= 104; i++ {
		deck = append(deck, i)
	}

	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck
}

func getSortedCards(played map[int]int) []int {
	var keys []int
	for k := range played {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
