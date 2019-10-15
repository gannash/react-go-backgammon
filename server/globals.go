package main

import (
	"time"
)

var state GameState

var gameManager GameManager

var history []GameState

var lastPlayTimestamp time.Time