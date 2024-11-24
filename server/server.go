package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sudoku-webapp/sudoku"
)

var funcMap = template.FuncMap{
	"mod":    func(a, b int) int { return a % b },
	"isClue": sudoku.IsClue,
}
var tmpl = template.Must(template.New("sudoku.html").Funcs(funcMap).ParseFiles("templ/sudoku.html"))
var validPath = regexp.MustCompile("^/(easy|medium|hard)/(gen|reset)?$")

var easySudoku = sudoku.GenerateSudoku(sudoku.Easy)
var mediumSudoku = sudoku.GenerateSudoku(sudoku.Medium)
var hardSudoku = sudoku.GenerateSudoku(sudoku.Hard)

var sudokuBoards = map[string]*sudoku.Sudoku{
	"easy":   &easySudoku,   // Initialize easy board
	"medium": &mediumSudoku, // Initialize medium board
	"hard":   &hardSudoku,   // Initialize hard board
}

func StartServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/easy/", easyHandler)
	http.HandleFunc("/medium/", mediumHandler)
	http.HandleFunc("/hard/", hardHandler)
	http.HandleFunc("/check-sudoku", checkSudokuHandler)
	http.HandleFunc("/update-sudoku", updateSudokuHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/medium", http.StatusFound)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil)) // Lan server
	//log.Fatal(http.ListenAndServe(":8000", nil)) // Local host
}

func easyHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	switch m[2] {
	case "gen":
		easySudoku = sudoku.GenerateSudoku(sudoku.Easy)
		http.Redirect(w, r, "/easy", http.StatusFound)
	case "reset":
		sudokuBoards["easy"].Reset()
		http.Redirect(w, r, "/easy", http.StatusFound)
	default:
		renderTemplate(w, sudokuBoards["easy"])
	}
}

func mediumHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	switch m[2] {
	case "gen":
		mediumSudoku = sudoku.GenerateSudoku(sudoku.Medium)
		http.Redirect(w, r, "/medium", http.StatusFound)
	case "reset":
		sudokuBoards["medium"].Reset()
		http.Redirect(w, r, "/medium", http.StatusFound)
	default:
		renderTemplate(w, sudokuBoards["medium"])
	}
}

func hardHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	switch m[2] {
	case "gen":
		hardSudoku = sudoku.GenerateSudoku(sudoku.Hard)
		http.Redirect(w, r, "/hard", http.StatusFound)
	case "reset":
		sudokuBoards["hard"].Reset()
		http.Redirect(w, r, "/hard", http.StatusFound)
	default:
		renderTemplate(w, sudokuBoards["hard"])
	}
}

func renderTemplate(w http.ResponseWriter, s *sudoku.Sudoku) {
	err := tmpl.Execute(w, s)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// Handler to receive updates for a specific cell and difficulty
func updateSudokuHandler(w http.ResponseWriter, r *http.Request) {
	var updatedCell struct {
		Difficulty string `json:"difficulty"`
		Row        int    `json:"row"`
		Col        int    `json:"col"`
		Value      int    `json:"value"`
	}

	// Decode the request body to get updated cell data
	err := json.NewDecoder(r.Body).Decode(&updatedCell)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get the correct Sudoku board based on the difficulty level
	sudoku, exists := sudokuBoards[updatedCell.Difficulty]
	if !exists {
		http.Error(w, "Invalid difficulty level", http.StatusBadRequest)
		return
	}

	// Update the specified cell in the appropriate Sudoku board
	sudoku.Insert(updatedCell.Row, updatedCell.Col, updatedCell.Value)
}

func checkSudokuHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Difficulty string `json:"difficulty"`
	}

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve the correct board based on difficulty
	sudoku, exists := sudokuBoards[requestData.Difficulty]
	if !exists {
		http.Error(w, "Invalid difficulty level", http.StatusBadRequest)
		return
	}

	// Validate the board
	correct := sudoku.IsSolved()
	message := "Incorrect solution"
	if correct {
		message = "Congratulations! You solved the puzzle."
	}

	// Send the result back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"correct": correct,
		"message": message,
	})
}
