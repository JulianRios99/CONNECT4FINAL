package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Game struct {
	Board         [][]int
	CurrentPlayer int
}

func NewGame(columns int, rows int) *Game {
	board := make([][]int, rows)
	for i := range board {
		board[i] = make([]int, columns)
	}
	return &Game{
		Board:         board,
		CurrentPlayer: 1,
	}
}

func (g *Game) MakeMove(col int) (int, error) {
	row := g.findLowestAvailableRow(col)
	if row == -1 {
		return -1, errors.New("column is full")
	}

	g.Board[row][col] = g.CurrentPlayer
	if g.checkWin(g.CurrentPlayer) {
		return g.CurrentPlayer, nil
	}

	g.CurrentPlayer = 3 - g.CurrentPlayer // Switch between 1 and 2
	return 0, nil
}

func (g *Game) MakeMoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	colStr := r.URL.Query().Get("col")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid column"})
		return
	}

	status, err := g.MakeMove(col)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var gameStatus string
	if status == 0 {
		gameStatus = "continue"
	} else {
		gameStatus = "win"
	}

	response := map[string]interface{}{
		"status": gameStatus,
		"row":    g.findLowestAvailableRow(col) + 1,
		"col":    col,
		"player": g.CurrentPlayer,
	}

	json.NewEncoder(w).Encode(response)
}

func (g *Game) findLowestAvailableRow(col int) int {
	for row := len(g.Board) - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			return row
		}
	}
	return -1
}

func (g *Game) checkWin(player int) bool {
	rows := len(g.Board)
	columns := len(g.Board[0])

	// Check rows
	for row := 0; row < rows; row++ {
		for col := 0; col < columns-3; col++ {
			if g.Board[row][col] == player &&
				g.Board[row][col+1] == player &&
				g.Board[row][col+2] == player &&
				g.Board[row][col+3] == player {
				return true
			}
		}
	}

	// Check columns
	for row := 0; row < rows-3; row++ {
		for col := 0; col < columns; col++ {
			if g.Board[row][col] == player &&
				g.Board[row+1][col] == player &&
				g.Board[row+2][col] == player &&
				g.Board[row+3][col] == player {
				return true
			}
		}
	}

	// Check ascending diagonals
	for row := 3; row < rows; row++ {
		for col := 0; col < columns-3; col++ {
			if g.Board[row][col] == player &&
				g.Board[row-1][col+1] == player &&
				g.Board[row-2][col+2] == player &&
				g.Board[row-3][col+3] == player {
				return true
			}
		}
	}

	// Check descending diagonals
	for row := 0; row < rows-3; row++ {
		for col := 0; col < columns-3; col++ {
			if g.Board[row][col] == player &&
				g.Board[row+1][col+1] == player &&
				g.Board[row+2][col+2] == player &&
				g.Board[row+3][col+3] == player {
				return true
			}
		}
	}

	return false
}
