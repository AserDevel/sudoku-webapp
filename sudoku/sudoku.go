package sudoku

import (
	"fmt"
	"math/rand"
	"slices"
)

type Clue struct {
	Row   int
	Col   int
	Value int
}

type Sudoku struct {
	Board [9][9]int
	Clues []Clue
}

// Generates a random sudoku
func GenerateSudoku(diff string) Sudoku {
	sudoku := Sudoku{}
	sudoku.Solve()

	switch diff {
	case "easy":
		for !sudoku.removeCells(35) {
			sudoku = Sudoku{}
			sudoku.Solve()
		}
	case "medium":
		for !sudoku.removeCells(45) {
			sudoku = Sudoku{}
			sudoku.Solve()
		}
	case "hard":
		for !sudoku.removeCells(55) {
			sudoku = Sudoku{}
			sudoku.Solve()
		}
	}

	// Generate the clue cells
	for r := range sudoku.Board {
		for c := range sudoku.Board[r] {
			if sudoku.Board[r][c] != 0 {
				sudoku.Clues = append(sudoku.Clues, Clue{
					Row:   r,
					Col:   c,
					Value: sudoku.Board[r][c],
				})
			}
		}
	}

	return sudoku
}

// Attempts to insert num at (r, c) returns true if successful
func (sudoku *Sudoku) Insert(r, c, num int) bool {
	if !inRange(r, 0, 8) || !inRange(c, 0, 8) || !inRange(num, 0, 9) {
		return false
	}
	// Check if the placement is a clue (locked) cell
	if IsClue(r, c, sudoku.Clues) {
		return false
	}
	sudoku.Board[r][c] = num
	return true
}

func (sudoku *Sudoku) IsSolved() bool {
	for r := range sudoku.Board {
		for c := range sudoku.Board[r] {
			if sudoku.Board[r][c] == 0 || !validIndex(sudoku, r, c) {
				return false
			}
		}
	}
	return true
}

func IsClue(r, c int, clues []Clue) bool {
	for _, clue := range clues {
		if clue.Row == r && clue.Col == c {
			return true
		}
	}
	return false
}

func (sudoku *Sudoku) Reset() {
	for r := range sudoku.Board {
		for c := range sudoku.Board[r] {
			sudoku.Board[r][c] = 0
		}
	}
	for _, clue := range sudoku.Clues {
		sudoku.Board[clue.Row][clue.Col] = clue.Value
	}
}

func (sudoku *Sudoku) Print() {
	line := "+---------+---------+---------+"
	for r := range sudoku.Board {
		fmt.Println("")
		if r%3 == 0 {
			fmt.Println(line)
		}
		for c := range sudoku.Board[r] {
			if c%3 == 0 {
				fmt.Print("|")
			}
			fmt.Printf(" %v ", sudoku.Board[r][c])
		}
		fmt.Print("|")
	}
	fmt.Printf("\n%v\n", line)
}

// Fills the sudoku with valid numbers in a random order
// The function is also used to generate random sudokus
func (sudoku *Sudoku) Solve() {
	var dfs func(r, c int) bool

	dfs = func(r, c int) bool {
		if r == 9 {
			return true // Grid is complete
		}
		nextR, nextC := (r*9+c+1)/9, (r*9+c+1)%9

		if sudoku.Board[r][c] != 0 {
			return dfs(nextR, nextC)
		}

		nums := rand.Perm(9)
		for _, n := range nums {
			sudoku.Board[r][c] = n + 1
			if validIndex(sudoku, r, c) {
				if dfs(nextR, nextC) {
					return true
				}
			}
			sudoku.Board[r][c] = 0
		}
		return false
	}

	dfs(0, 0)
}

// Boundary checker
func inRange(n, min, max int) bool {
	if n < min || n > max {
		return false
	}
	return true
}

// Randomly removes the given amount of cells from the sudoku,
// while keeping it valid. Returns true if succesful.
func (sudoku *Sudoku) removeCells(amount int) bool {
	if amount > 55 { // Cap to avoid infinite loops and long generation time
		amount = 55
	}

	removed := 0
	nums := rand.Perm(81)
	for i := 0; removed < amount && i < 81; i++ {
		r, c := nums[i]/9, nums[i]%9

		tmp := sudoku.Board[r][c]
		sudoku.Board[r][c] = 0

		// Check if the puzzle remains valid with a unique solution
		if valid(*sudoku) {
			removed++
		} else {
			sudoku.Board[r][c] = tmp
		}
	}

	return removed == amount
}

// Takes a soduko and returns true if it has exactly 1 solution
func valid(sudoku Sudoku) bool {
	var dfs func(r, c int) bool
	solutions := 0

	dfs = func(r, c int) bool {
		if r == 9 {
			solutions++
			return solutions > 1 // Stop searching if more than 1 solution
		}
		nextR, nextC := (r*9+c+1)/9, (r*9+c+1)%9 // Moves to the next cell

		if sudoku.Board[r][c] != 0 {
			return dfs(nextR, nextC)
		}

		for n := 1; n <= 9; n++ {
			sudoku.Board[r][c] = n
			if validIndex(&sudoku, r, c) && dfs(nextR, nextC) {
				return true
			}
			sudoku.Board[r][c] = 0
		}
		return false
	}

	dfs(0, 0)
	return solutions == 1
}

// Returns false if the number on the given cell violates the sudoku rules
func validIndex(sudoku *Sudoku, r, c int) bool {
	return (validRow(sudoku, r) && validCol(sudoku, c) && validSquare(sudoku, r, c))
}

func validRow(sudoku *Sudoku, r int) bool {
	return validNums(sudoku.Board[r])
}

func validCol(sudoku *Sudoku, c int) bool {
	var nums [9]int
	for r := range sudoku.Board {
		nums[r] = sudoku.Board[r][c]
	}
	return validNums(nums)
}

func validSquare(sudoku *Sudoku, r, c int) bool {
	var nums [9]int
	r -= r % 3 // Round to the corner of square
	c -= c % 3

	count := 0
	for i := r; i < r+3; i++ {
		for j := c; j < c+3; j++ {
			nums[count] = sudoku.Board[i][j]
			count++
		}
	}
	return validNums(nums)
}

// Returns true if all numbers in the array are unique or 0
func validNums(nums [9]int) bool {
	count := slices.Repeat([]int{1}, 9)
	for _, num := range nums {
		if num == 0 {
			continue
		}
		count[num-1]--
		if count[num-1] < 0 {
			return false
		}
	}
	return true
}
