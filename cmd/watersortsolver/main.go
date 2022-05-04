package main

import (
	"encoding/json"
	"flag"
	"fmt"
	watersortpuzzle "github.com/pkositsyn/water-sort-puzzle-solver"
	"net/http"
)

var algorithmType = flag.String("algorithm", "astar",
	`Algorithm to solve with. Choices: [astar, idastar, dijkstra]`)

type Result struct {
	Message string                 `json:"message"`
	Step    int                    `json:"step"`
	Steps   []watersortpuzzle.Step `json:"steps"`
}

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("s")
		if s == "" {
			result := Result{
				Message: "s query param is empty",
			}
			json.NewEncoder(w).Encode(result)
		} else {
			solve(r.URL.Query().Get("s"), w)
		}
	})
	http.ListenAndServe(":9280", nil)
}

func solve(s string, w http.ResponseWriter) {
	var solver watersortpuzzle.Solver
	switch *algorithmType {
	case "astar":
		solver = watersortpuzzle.NewAStarSolver()
	case "idastar":
		solver = watersortpuzzle.NewIDAStarSolver()
	case "dijkstra":
		solver = watersortpuzzle.NewDijkstraSolver()
	}

	var initialState watersortpuzzle.State
	message := ""
	if err := initialState.FromString(s); err != nil {
		message = fmt.Sprintf("Invalid puzzle state provided: %s", err.Error())
		result := Result{
			Message: message,
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	steps, err := solver.Solve(initialState)
	if err != nil {
		message = fmt.Sprintf("Cannot solve puzzle: %s", err.Error())
		result := Result{
			Message: message,
		}
		json.NewEncoder(w).Encode(result)
		return

	}

	suffix := ""
	if statsSolver, ok := solver.(watersortpuzzle.SolverWithStats); ok {
		suffix = fmt.Sprintf(" Algorithm took %d iterations to find solution.", statsSolver.Stats().Steps)
	}

	message = fmt.Sprintf("Puzzle solved in %d steps!%s", len(steps), suffix)

	//for _, step := range steps {
	//	fmt.Println(step.From+1, step.To+1)
	//}

	result := Result{
		Message: message,
		Step:    len(steps),
		Steps:   steps,
	}
	json.NewEncoder(w).Encode(result)
}
