package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func registerAPI(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("got move command")
	decoder := json.NewDecoder(r.Body)
	var p Player
	err := decoder.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect data structure")
	} else {
		playerID := rand.Intn(100000000)
		p.PlayerID = playerID

		if gameManager.player1.TeamName == "" {
			p.Color = "White"
			gameManager.player1 = p
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
		} else if gameManager.player2.TeamName == "" {
			p.Color = "Black"
			gameManager.player2 = p
			state.state = GAME_RUNNING
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "2 players have already joined the match")
		}
	}

	//check if it's the coorect players turn to play

	// var playerColor string

	// if t.PlayerID == gameManager.player1.playerID {
	// 	playerColor = gameManager.player1.color
	// }

	// if t.PlayerID == gameManager.player2.playerID {
	// 	playerColor = gameManager.player2.color
	// }

	// if playerColor == "White" && !state.WhiteTurn {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "not your turn,, it's black's turn ")
	// }

	// if playerColor == "Black" && state.WhiteTurn {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "not your turn, it's white's turn")
	// }

	// //check if the move is correct
	// if !isLegalMove(&state, t.From, t.To) {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "move is illegal")
	// 	return
	// 	//complete automatic move
	// }

	// move(&state, t.From, t.To)

	// if !canMove() {
	// 	sendStateWithError(w)
	// }

	// // fmt.Println(t)
}
