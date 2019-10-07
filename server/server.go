package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//////////////////// end of functions //////////////////////////////

// main function
func main() {
	rand.Seed(time.Now().UnixNano())

	// fmt.Println("Running BG 1.0")
	router := mux.NewRouter()
	router.HandleFunc("/getBoard", getBoardAPI).Methods("GET")
	router.HandleFunc("/move", moveAPI).Methods("POST")
	router.HandleFunc("/register", registerAPI).Methods("POST")
	router.HandleFunc("/getState", getStateAPI).Methods("GET")

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := corsConfig.Handler(router)

	initState(&state)
	history = append(history, state)
	diceCopy := []int{state.Dice[0], state.Dice[1], state.Dice[2], state.Dice[3]}
	history[len(history)-1].Dice = diceCopy

	// printBoard(state.Board)

	fmt.Println("running on port 7861...")

	log.Fatal(http.ListenAndServe(":7861", handler))
}
