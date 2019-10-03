package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func moveAPI(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("got move command")
	decoder := json.NewDecoder(r.Body)
	var t MoveCommand
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect data structure ")
	}

	//check if it's the coorect players turn to play

	var playerColor string

	if t.PlayerID == gameManager.player1.PlayerID {
		playerColor = gameManager.player1.Color
	}

	if t.PlayerID == gameManager.player2.PlayerID {
		playerColor = gameManager.player2.Color
	}

	if playerColor == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong playerID")
		return
	}

	if playerColor == "White" && !state.WhiteTurn {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not your turn, it's black's turn")
		return
	}

	if playerColor == "Black" && state.WhiteTurn {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not your turn, it's white's turn")
		return
	}

	//check if the move is correct
	if !isLegalMove(&state, t.From, t.To) {
		// w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(w, "move is illegal, a random move has been played")

		//complete automatic move

		state.Status = "ILLEGAL_MOVE"
		randomMove()
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(state)
	} else {
		move(&state, t.From, t.To)
	}
	if canMove() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state)
	} else {
		sendStateWithError(w)
	}

	// fmt.Println(t)
}