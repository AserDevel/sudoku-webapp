package sudoku

import "testing"

func TestValidation(t *testing.T) {
	sudoku := Sudoku{}
	if valid(sudoku) {
		sudoku.Print()
		t.Fatalf("Expected valid to be false, but returned true")
	}
	sudoku = GenerateSudoku("hard")
	if !valid(sudoku) {
		sudoku.Print()
		t.Fatalf("Expected valid to be true, but returned false")
	}
	sudoku.Solve()
	if !valid(sudoku) {
		sudoku.Print()
		t.Fatalf("Expected valid to be true, but returned false")
	}
}

func TestGeneration(t *testing.T) {
	var sudoku Sudoku

	// Test easy
	count := 81 - 35
	for i := 0; i < 15; i++ {
		sudoku = GenerateSudoku("easy")
		testGenerationHelper(sudoku, count, t)
	}

	// Test medium
	count = 81 - 45
	for i := 0; i < 15; i++ {
		sudoku = GenerateSudoku("medium")
		testGenerationHelper(sudoku, count, t)
	}

	// Test hard
	count = 81 - 55
	for i := 0; i < 15; i++ {
		sudoku = GenerateSudoku("hard")
		testGenerationHelper(sudoku, count, t)
	}
}

func testGenerationHelper(sudoku Sudoku, count int, t *testing.T) {
	if !valid(sudoku) {
		sudoku.Print()
		t.Fatalf("Sudoku does not have exactly 1 solution")
	}
	if len(sudoku.Clues) != count {
		sudoku.Print()
		t.Fatalf("Expected sudoku to have: %v clues, but found: %v clues", count, len(sudoku.Clues))
	}
}

func TestSolving(t *testing.T) {
	for i := 0; i < 15; i++ {
		sudoku := GenerateSudoku("hard")
		if sudoku.IsSolved() {
			sudoku.Print()
			t.Fatalf("IsSolved was true, but expected it to be false")
		}
		sudoku.Solve()
		if !sudoku.IsSolved() {
			sudoku.Print()
			t.Fatalf("IsSolved was false, but expected it to be true")
		}
	}
}
