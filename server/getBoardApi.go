package main

import (
	"encoding/json"
	"net/http"
)

func getBoardAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(state.Board)
	//fmt.Fprintf(w, state.board)
}