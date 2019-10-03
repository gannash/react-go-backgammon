package main

import (
	"encoding/json"
	"net/http"
)

func getStateAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	state.Status = ""

	if !canMove() {
		sendStateWithError(w)
	} else {
		json.NewEncoder(w).Encode(state)
	}
	// fmt.Println(gameManager.player1)
	// fmt.Println(gameManager.player2)
}