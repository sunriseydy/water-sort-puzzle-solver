package main

import (
	"flag"
	"fmt"
	watersortpuzzle "github.com/pkositsyn/water-sort-puzzle-solver"
	"html/template"
	"net/http"
)

var algorithmType = flag.String("algorithm", "astar",
	`Algorithm to solve with. Choices: [astar, idastar, dijkstra]`)

type Result struct {
	Success bool
	Message string
	Step    int
	Steps   []watersortpuzzle.Step
}

func main() {
	flag.Parse()
	tmpl := template.Must(template.ParseFiles("form.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
		s := r.FormValue("str")
		if s == "" {
			result := Result{
				Success: false,
				Message: "s query param is empty",
			}
			tmpl.Execute(w, result)
		} else {
			result := solve(s)
			tmpl.Execute(w, result)
		}
	})
	http.ListenAndServe(":9280", nil)
}

func solve(s string) Result {
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
			Success: false,
			Message: message,
		}
		return result
	}

	steps, err := solver.Solve(initialState)
	if err != nil {
		message = fmt.Sprintf("Cannot solve puzzle: %s", err.Error())
		result := Result{
			Success: false,
			Message: message,
		}
		return result

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
		Success: true,
	}
	return result
}
