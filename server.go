package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//////////////////////// Types //////////////////////////////

//GameState ..
type GameState struct {
	Board     Board
	WhiteTurn bool
	Dice      []int
	WhiteWon  bool
	BlackWon  bool
	state     string //waiting for players, game on, game ended
}

//GameManager ..
type GameManager struct {
	player1  Player
	player2  Player
	state    GameState
	gameMode string
}

//Player ..
type Player struct {
	color    string
	teamName string
	playerID int
}

//Column ..
type Column struct {
	WhiteCheckers int `json:"whiteCheckers"`
	BlackCheckers int `json:"BlackCheckers"`
}

type GameStateWithError struct {
	Error string    `json:"error"`
	State GameState `json:"state"`
}

//Board ..
type Board struct {
	Columns    [24]Column `json:"columns"`
	WhiteEaten int        `json:"whiteEaten"`
	BlackEaten int        `json:"blackEaten"`

	WhiteOutCheckers int `json:"whiteOutCheckers"`
	BlackOutCheckers int `json:"blackOutCheckers"`
}

//MoveCommand ..
type MoveCommand struct {
	PlayerID int `json:"playerID"`
	From     int `json:"from"`
	To       int `json:"to"`
}

//////////////////////// end of Types //////////////////////////////

//////////////////// Global Variables /////////////////////////////

var state GameState

var gameManager GameManager

//////////////////// end of Global Variables /////////////////////////////

//////////////////// functions //////////////////////////////

//////////// APIS ////////////
func getBoardAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(state.Board)
	//fmt.Fprintf(w, state.board)
}

// func sendCantMove(w http.ResponseWriter) {
// 	w.WriteHeader(http.StatusLocked)
// 	fmt.Fprintf(w, "no available moves")

// 	time.AfterFunc(2*time.Second, switchTurn)
// }

func sendStateWithError(w http.ResponseWriter) {
	var gameStateWithError GameStateWithError
	gameStateWithError.Error = "NO_MOVES"
	gameStateWithError.State = state
	json.NewEncoder(w).Encode(gameStateWithError)

	time.AfterFunc(2*time.Second, switchTurn)
}

func getStateAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !canMove() {
		sendStateWithError(w)
	} else {
		json.NewEncoder(w).Encode(state)
	}
	// fmt.Println(gameManager.player1)
	// fmt.Println(gameManager.player2)
}

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

	if t.PlayerID == gameManager.player1.playerID {
		playerColor = gameManager.player1.color
	}

	if t.PlayerID == gameManager.player2.playerID {
		playerColor = gameManager.player2.color
	}

	if playerColor == "White" && !state.WhiteTurn {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not your turn, it's black's turn ")
	}

	if playerColor == "Black" && state.WhiteTurn {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not your turn, it's white's turn")
		fmt.Println("cant move")
	}

	//check if the move is correct
	if !isLegalMove(&state, t.From, t.To) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "move is illegal")

		//complete automatic move
	} else {
		move(&state, t.From, t.To)

		if canMove() {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(state)
		} else {
			sendStateWithError(w)
		}
	}

	// fmt.Println(t)
}

func registerAPI(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("got move command")
	decoder := json.NewDecoder(r.Body)
	var t MoveCommand
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect data structure")
	}

	//check if it's the coorect players turn to play

	var playerColor string

	if t.PlayerID == gameManager.player1.playerID {
		playerColor = gameManager.player1.color
	}

	if t.PlayerID == gameManager.player2.playerID {
		playerColor = gameManager.player2.color
	}

	if playerColor == "White" && !state.WhiteTurn {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not your turn,, it's black's turn ")
	}

	if playerColor == "Black" && state.WhiteTurn {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not your turn, it's white's turn")
	}

	//check if the move is correct
	if !isLegalMove(&state, t.From, t.To) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "move is illegal")

		//complete automatic move
	} else {
		move(&state, t.From, t.To)

		if !canMove() {
			sendStateWithError(w)
		}
	}

	// fmt.Println(t)
}

//////////// end of APIS ////////////

func getState() GameState {
	return state
}

func initState(state *GameState) {
	state.WhiteTurn = true
	state.WhiteWon = false
	state.BlackWon = false
	initBoard(&state.Board)
	state.Dice = rollDice()
	// For tests
	// state.Dice = []int{6, 6, 6, 6}
	// state.Board.WhiteEaten = 1
	state.state = "waitingForPlayers"

	gameManager.player1.playerID = 0
	gameManager.player1.color = "White"
	gameManager.player1.teamName = "Eminem"

	gameManager.player2.playerID = 1
	gameManager.player2.color = "Black"
	gameManager.player2.teamName = "Tupac"
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

func isLegalMove(state *GameState, from int, to int) bool {

	if state.WhiteWon || state.BlackWon {
		return false
	}

	//check correct checkers are moving accoridng to player turn
	if state.WhiteTurn && from != -1 && state.Board.Columns[from].WhiteCheckers == 0 {
		return false
	}

	if !state.WhiteTurn && from != -1 && state.Board.Columns[from].BlackCheckers == 0 {
		return false
	}

	if state.WhiteTurn {

		if from == -1 {

			//no checkers out
			if state.Board.WhiteEaten < 1 {
				return false
			}

			//target outside of home
			if to > 6 {
				return false
			}

			//trying to land on a house
			if state.Board.Columns[to].BlackCheckers >= 2 {
				return false
			}

			// not according to dice
			possibleDiceMoves := false
			for _, die := range state.Dice {
				if (to + 1) == die {
					possibleDiceMoves = true
				}
			}

			return possibleDiceMoves
		} else if to == -1 {
			//check if all checkers are at home

			allCheckersAtHome := true

			for i := 0; i < 18; i++ {
				if state.Board.Columns[i].WhiteCheckers > 0 {
					allCheckersAtHome = false
				}
			}

			if !allCheckersAtHome {
				return false
			}

			//check if dice can take out a checker

			distance := 24 - from
			diceMatchFound := false
			for _, itemCopy := range state.Dice {
				if distance <= itemCopy {
					diceMatchFound = true
					break
				}
			}

			if !diceMatchFound {
				return false
			}

		} else { //regular play

			if to > 23 || to < 0 || from < 0 || from > 23 {
				return false
			}

			//check if move is in range and in dice
			distanceToMove := to - from
			diceMatch := false

			for _, itemCopy := range state.Dice {
				if distanceToMove == itemCopy {
					diceMatch = true
					break
				}
				// fmt.Printf("%v: %v\n", index, itemCopy)
			}

			if !diceMatch {
				return false
			}

			//check that the to does not have a house in it
			if state.Board.Columns[to].BlackCheckers >= 2 {
				return false
			}

		}

	}

	//in case it's black's turn
	if !state.WhiteTurn {

		if from == -1 {

			//no checkers out
			if state.Board.BlackEaten < 1 {
				return false
			}

			//target outside of home
			if (24 - to) > 6 {
				return false
			}

			//trying to land on a house
			if state.Board.Columns[to].WhiteCheckers >= 2 {
				return false
			}

			// not according to dice
			possibleDiceMoves := false
			for _, die := range state.Dice {
				if (24 - to) == die {
					possibleDiceMoves = true
				}
			}

			return possibleDiceMoves
		} else if to == -1 {
			//check if all checkers are at home

			allCheckersAtHome := true

			for i := 23; i > 6; i-- {
				if state.Board.Columns[i].BlackCheckers > 0 {
					allCheckersAtHome = false
				}
			}

			if !allCheckersAtHome {
				return false
			}

			//check if dice can take out a checker

			distance := from
			diceMatchFound := false
			for _, itemCopy := range state.Dice {
				if distance <= itemCopy {
					diceMatchFound = true
					break
				}
			}

			if !diceMatchFound {
				return false
			}

		} else { //regular play

			if to > 23 || to < 0 || from < 0 || from > 23 {
				return false
			}

			//check if move is in range and in dice
			distanceToMove := from - to
			diceMatch := false

			for _, itemCopy := range state.Dice {
				if distanceToMove == itemCopy {
					diceMatch = true
					break
				}
				// fmt.Printf("%v: %v\n", index, itemCopy)
			}

			if !diceMatch {
				return false
			}

			//check that the to does not have a house in it
			if state.Board.Columns[to].WhiteCheckers >= 2 {
				return false
			}

		}

	}

	return true
}

func isCurrentPlayerEaten() bool {
	return (state.WhiteTurn && state.Board.WhiteEaten > 0) || (!state.WhiteTurn && state.Board.BlackEaten > 0)
}

func anyReturnToBoardDestination() bool {
	for _, die := range state.Dice {
		if die > 0 {
			if state.WhiteTurn {
				destination := die - 1
				if state.Board.Columns[destination].BlackCheckers < 2 {
					return true
				}
			} else {
				destination := 24 - die
				if state.Board.Columns[destination].WhiteCheckers < 2 {
					return true
				}
			}
		}
	}

	return false
}

func allPlayersCheckersOnHome() bool {
	// assumes that no checkers are on grey board
	if state.WhiteTurn {
		for idx, column := range state.Board.Columns {
			if idx < 19 && column.WhiteCheckers > 0 {
				return false
			}
		}
	} else { // black turn
		for idx, column := range state.Board.Columns {
			if idx > 5 && column.BlackCheckers > 0 {
				return false
			}
		}
	}

	return true
}

func anyCheckerToBearOff() bool {
	for _, die := range state.Dice {
		if die > 0 {
			if state.WhiteTurn {
				for i := 0; i <= die-1; i++ {
					if state.Board.Columns[i].WhiteCheckers > 0 {
						return true
					}
				}
			} else {
				for i := 23; i >= 24-die; i++ {
					if state.Board.Columns[i].BlackCheckers > 0 {
						return true
					}
				}
			}
		}
	}

	return false
}

func anyCheckerCanMoveWithDice() bool {
	for idx, column := range state.Board.Columns {
		if (state.WhiteTurn && column.WhiteCheckers > 0) || (!state.WhiteTurn && column.BlackCheckers > 0) {
			for _, die := range state.Dice {
				if die > 0 {
					if state.WhiteTurn {
						destination := die + idx
						if destination >= 0 && destination < 24 && state.Board.Columns[destination].BlackCheckers == 0 {
							return true
						}
					} else {
						destination := idx - die
						if destination >= 0 && destination < 24 && state.Board.Columns[destination].WhiteCheckers == 0 {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func canMove() bool {
	if isCurrentPlayerEaten() {
		// fmt.Println("inside isCurrentPlayerEaten condition")
		return anyReturnToBoardDestination()
	} else if allPlayersCheckersOnHome() {
		// fmt.Println("inside allPlayersCheckersOnHome condition")
		return anyCheckerToBearOff() || anyCheckerCanMoveWithDice()
	}

	// fmt.Println("checking anyCheckerCanMoveWithDice")
	return anyCheckerCanMoveWithDice()
}

func move(state *GameState, from int, to int) int {
	// // check for possible moves
	// if !canMove() {
	// 	switchTurn()
	// }

	//player bringing an eaten checker back into game
	if from == -1 {

		if state.WhiteTurn {
			state.Board.WhiteEaten--
			state.Board.Columns[to].WhiteCheckers++

			//check if eating a black piece
			if state.Board.Columns[to].BlackCheckers > 0 {
				state.Board.Columns[to].BlackCheckers--
				state.Board.BlackEaten++
			}

		} else { //black turn
			state.Board.BlackEaten--
			state.Board.Columns[to].BlackCheckers++

			//check if eating a white piece
			if state.Board.Columns[to].WhiteCheckers > 0 {
				state.Board.Columns[to].WhiteCheckers--
				state.Board.WhiteEaten++
			}
		}

	}

	//player taking a checker out of the game
	if to == -1 {
		if state.WhiteTurn {
			state.Board.Columns[from].WhiteCheckers--
			state.Board.WhiteOutCheckers++
		} else {
			state.Board.Columns[from].BlackCheckers--
			state.Board.BlackOutCheckers++
		}
	}

	//regular play in board
	if from >= 0 && from <= 23 && to >= 0 && to <= 23 {
		if state.WhiteTurn {
			state.Board.Columns[from].WhiteCheckers--

			state.Board.Columns[to].WhiteCheckers++

			//check if eating a black piece
			if state.Board.Columns[to].BlackCheckers > 0 {
				state.Board.Columns[to].BlackCheckers--
				state.Board.BlackEaten++
			}
		} else {
			state.Board.Columns[from].BlackCheckers--

			state.Board.Columns[to].BlackCheckers++

			//check if eating a black piece
			if state.Board.Columns[to].WhiteCheckers > 0 {
				state.Board.Columns[to].WhiteCheckers--
				state.Board.WhiteEaten++
			}
		}

	}

	//remove the dice from remaining dice
	dicePlayed := -50
	if from == -1 {
		if state.WhiteTurn {
			dicePlayed = to + 1
		} else {
			dicePlayed = 24 - to
		}
	} else if to == -1 {
		if state.WhiteTurn {
			dicePlayed = 24 - from
		} else {
			dicePlayed = from + 1
		}
	} else {
		if state.WhiteTurn {
			dicePlayed = to - from
		} else {
			dicePlayed = from - to
		}
	}

	// fmt.Println(dicePlayed)

	for index, die := range state.Dice {
		if dicePlayed == die {
			state.Dice[index] = 0
			break
		}
		// fmt.Printf("deleted die %v: %v\n", index, die)
	}

	//check if all dice were played and if so swith turns

	shouldSwitchTurn := true
	for _, die := range state.Dice {
		if die != 0 {
			shouldSwitchTurn = false
		}
	}

	if shouldSwitchTurn {
		switchTurn()
	}

	return 0
}

func switchTurn() {
	state.WhiteTurn = !state.WhiteTurn
	state.Dice = rollDice()
}

func initBoard(board *Board) {

	for i := 0; i < 24; i++ {
		board.Columns[i] = Column{WhiteCheckers: 0, BlackCheckers: 0}
	}

	board.Columns[0] = Column{WhiteCheckers: 2, BlackCheckers: 0}
	board.Columns[11] = Column{WhiteCheckers: 5, BlackCheckers: 0}
	board.Columns[16] = Column{WhiteCheckers: 3, BlackCheckers: 0}
	board.Columns[18] = Column{WhiteCheckers: 5, BlackCheckers: 0}

	board.Columns[23] = Column{WhiteCheckers: 0, BlackCheckers: 2}
	board.Columns[12] = Column{WhiteCheckers: 0, BlackCheckers: 5}
	board.Columns[7] = Column{WhiteCheckers: 0, BlackCheckers: 3}
	board.Columns[5] = Column{WhiteCheckers: 0, BlackCheckers: 5}
}

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
	printBoard(state.Board)

	log.Fatal(http.ListenAndServe(":7861", handler))
}
