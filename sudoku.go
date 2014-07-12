/* This application takes a json sudoku board as input (stdin), and returns a
 * sudoku board in json as output (stdout).
 * If an error occurs (ie board invalid, input not valid) an error string is
 * written to stderr and no stdout is supplied.
 */
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Board []int

// Returns true/false, and an error if the board is not valid.
func (b Board) IsValid() (bool, error) {
	// Validate the length of the board.
	if len(b) != 81 {
		return false, errors.New("Board is not 9x9.")
	}

	// Validate that the numbers are 0-9.
	for i, val := range b {
		if val < 0 || val > 9 {
			error := fmt.Sprintf(
				"Internal number is not between 0 and 9 at position: %s",
				i)
			return false, errors.New(error)
		}
	}

	return true, nil
}

// A pretty string repressenting the board.
func (b Board) String() string {
	buffer := bytes.NewBufferString("")
	for y := 0; y < 9; y++ {
		if y > 0 && y%3 == 0 {
			buffer.WriteString("---+---+---\n")
		}
		for x := 0; x < 9; x++ {
			if x > 0 && x%3 == 0 {
				buffer.WriteString("|")
			}
			i := b[y*9+x]
			if i == 0 {
				buffer.WriteString(".")
			} else {
				buffer.WriteString(fmt.Sprintf("%d", i))
			}
		}
		if y < 8 {
			buffer.WriteString("\n")
		}
	}
	output, _ := ioutil.ReadAll(buffer)
	return string(output)
}

// Solves the board, returns a solved board, or nil if the board cannot be solved.
func (b Board) Solve() Board {
	// Validate the board.
	_, err := b.IsValid()
	if err != nil {
		return nil
	}

	// Solve using backtrack
	return b.backtrack(b, 0, 0)
}

func (b Board) deepcopy(board Board) Board {
	result := make(Board, 81)
	copy(result, board)
	return result
}

func (b Board) backtrack(board Board, x int, y int) Board {
	board = b.deepcopy(board)

	// Skip positions with existing data.
	if board[y*9+x] != 0 {
		return b.next(board, x, y)
	}

	// Iterate on possible solutions.
	for i := 1; i <= 9; i++ {
		if !b.check(board, i, x, y) {
			continue
		}
		board[y*9+x] = i
		result := b.next(board, x, y)
		if result != nil {
			return result
		}
	}

	// No solution found.
	return nil
}

func (b Board) next(board Board, x int, y int) Board {
	if x == 8 {
		if y == 8 {
			return board
		}
		return b.backtrack(board, 0, y+1)
	} else {
		return b.backtrack(board, x+1, y)
	}
}

//
func (b Board) check(board Board, val int, x int, y int) bool {
	// Validate horizontal.
	for _x := 0; _x < 9; _x++ {
		if _x != x {
			if board[y*9+_x] == val {
				return false
			}
		}
	}

	// Validate vertical.
	for _y := 0; _y < 9; _y++ {
		if _y != y {
			if board[_y*9+x] == val {
				return false
			}
		}
	}

	// check the current box.
	ybox := (y / 3) * 3
	xbox := (x / 3) * 3
	for _x := xbox; _x < xbox+3; _x++ {
		for _y := ybox; _y < ybox+3; _y++ {
			if _y != y || _x != x {
				if board[_y*9+_x] == val {
					return false
				}
			}
		}
	}

	return true
}

// Read from stdin. Write to stdout, or stderr and return non-0 return code.
func main() {
	// Read stdin.
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	if len(bytes) == 0 {
		fmt.Fprintln(os.Stderr, "No input")
		os.Exit(1)
	}

	// Parse json.
	board := Board{}
	err = json.Unmarshal(bytes, &board)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Validate that board is valid.
	_, err = board.IsValid()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// solve, or fail.
	board = board.Solve()

	// write the result.
	result, err := json.Marshal(board)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", result)
}
