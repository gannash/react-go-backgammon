package main

// "math/rand"

func initState(state *GameState) {
	state.WhiteTurn = false //rand.Intn(2) == 0
	state.WhiteWon = false
	state.BlackWon = false
	state.Dice = rollDice()

	initBoard(&state.Board)
	state.State = WAITING_PLAYERS
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
