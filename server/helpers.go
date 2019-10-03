package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func minMax(array []int) (int, int) {
	var max = array[0]
	var min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func switchTurn() {
	state.WhiteTurn = !state.WhiteTurn
	state.Dice = rollDice()
}


func printBoard(board Board) {
	fmt.Printf("%+v\n", board)
}

func rollDice() []int {
	var dice []int

	dice = append(dice, rand.Intn(6)+1)
	dice = append(dice, rand.Intn(6)+1)

	if dice[0] == dice[1] {
		dice = append(dice, dice[0], dice[0])
	} else {
		dice = append(dice, 0, 0)
	}

	// fmt.Println(dice)
	return dice
}

func sendStateWithError(w http.ResponseWriter) {
	if state.Status != "" {
		state.Status += "|"
	}
	state.Status += "NO_MOVES"

	json.NewEncoder(w).Encode(state)
	time.AfterFunc(2*time.Second, switchTurn)
}