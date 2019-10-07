package main

type GameState struct {
	Board     Board
	WhiteTurn bool
	Dice      []int
	WhiteWon  bool
	BlackWon  bool
	Status    string
	State     string //waiting for players, game on, game ended
	Players	  [2]Player
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
	Color    string `json:"color"`
	TeamName string `json:"teamName"`
	PlayerID int    `json:"playerID"`
}

//Column ..
type Column struct {
	WhiteCheckers int `json:"whiteCheckers"`
	BlackCheckers int `json:"BlackCheckers"`
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