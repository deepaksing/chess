package store

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func ConvertChessboardToString(board [][]int) string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, row := range board {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("{")
		for j, val := range row {
			if j > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(fmt.Sprintf("%d", val))
		}
		buffer.WriteString("}")
	}
	buffer.WriteString("}")
	return buffer.String()
}

func ConvertStringToChessboard(chessboardString string) [][]int {
	// Remove surrounding braces
	chessboardString = strings.Trim(chessboardString, "{}")

	// Split the string into rows
	rows := strings.Split(chessboardString, "},{")
	chessboard := make([][]int, len(rows))

	// Parse each row
	for i, row := range rows {
		// Split the row into columns
		cols := strings.Split(row, ",")
		chessboard[i] = make([]int, len(cols))

		// Parse each column
		for j, col := range cols {
			val, err := strconv.Atoi(col)
			if err != nil {
				// Handle conversion error
				return nil
			}
			chessboard[i][j] = val
		}
	}

	return chessboard
}

func NewChessBoard() string {

	board := [][]int{
		{-6, -5, -4, -3, -2, -4, -5, -6},
		{-1, -1, -1, -1, -1, -1, -1, -1},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{6, 5, 4, 3, 2, 4, 5, 6},
	}

	return ConvertChessboardToString(board)
}

func getCol(move_col byte) int {
	switch move_col {
	case 97:
		return 0
	case 98:
		return 1
	// Add cases for the remaining values
	case 99:
		return 2
	case 100:
		return 3
	case 101:
		return 4
	case 102:
		return 5
	case 103:
		return 6
	case 104:
		return 7
	default:
		// Handle other cases
		return -1
	}
}

func UpdateBoardState(board [][]int, move_from string, move_to string, move_type int) ([][]int, error) {
	// fmt.Println("b1 ", board)
	fromCol := getCol(move_from[0])           // Convert column letter to index (A=0, B=1, ..., H=7)
	fromRow := int(7 - int(move_from[1]-'1')) // Convert row number to index (1=0, 2=1, ..., 8=7)
	toCol := getCol(move_to[0])
	toRow := int(7 - int(move_to[1]-'1'))

	fmt.Println(fromRow, " ", fromCol)
	// fmt.Println(move_from[1], " ", fromCol)
	fmt.Println(toRow)
	fmt.Println(toCol)

	piece := board[fromRow][fromCol]
	// fmt.Println("piece ", piece)

	board[fromRow][fromCol] = 0
	// fmt.Println("b2 ", board)
	board[toRow][toCol] = piece
	// fmt.Println("b3 ", board)
	return board, nil
}
