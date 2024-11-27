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
var validPath = regexp.MustCompile("^/(easy|medium|hard)/(gen|reset|check|insert)?$")

var sudokuBoards = map[string]*sudoku.Sudoku{}

// Initializes everything and starts the server
func StartServer() {
	easy := sudoku.GenerateSudoku("easy")
	medium := sudoku.GenerateSudoku("medium")
	hard := sudoku.GenerateSudoku("hard")

	sudokuBoards["easy"] = &easy
	sudokuBoards["medium"] = &medium
	sudokuBoards["hard"] = &hard

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/easy/", makeHandler("easy"))
	http.HandleFunc("/medium/", makeHandler("medium"))
	http.HandleFunc("/hard/", makeHandler("hard"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/medium", http.StatusFound)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil)) // Lan server
	//log.Fatal(http.ListenAndServe(":8080", nil)) // Local host
}

// Returns a handler for the given difficulty
func makeHandler(diff string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		switch m[2] {
		case "gen":
			*sudokuBoards[diff] = sudoku.GenerateSudoku(diff)
			http.Redirect(w, r, "/"+diff, http.StatusFound)
		case "reset":
			sudokuBoards[diff].Reset()
			http.Redirect(w, r, "/"+diff, http.StatusFound)
		case "check":
			checkResponse(w, sudokuBoards[diff].IsSolved())
		case "insert":
			insertRequest(w, r, sudokuBoards[diff])
		default:
			renderTemplate(w, sudokuBoards[diff])
		}
	}
}

func renderTemplate(w http.ResponseWriter, s *sudoku.Sudoku) {
	err := tmpl.Execute(w, s)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// Updates the sudoku based on the request
func insertRequest(w http.ResponseWriter, r *http.Request, s *sudoku.Sudoku) {
	var updatedCell struct {
		Row   int `json:"row"`
		Col   int `json:"col"`
		Value int `json:"value"`
	}

	// Decode the request body to get updated cell data
	err := json.NewDecoder(r.Body).Decode(&updatedCell)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update the specified cell in the appropriate Sudoku board
	s.Insert(updatedCell.Row, updatedCell.Col, updatedCell.Value)
}

// Responds to the check request
func checkResponse(w http.ResponseWriter, correct bool) {
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
