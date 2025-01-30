package routes

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"restoportGo/internal/services"
	"strconv"
	"text/template"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HomeHandler)

	mux.HandleFunc("/CarlosFlights", flightsHandler(db, 62, 69))

	mux.HandleFunc("/UbarFlights", flightsHandler(db, 51, 69))

	mux.HandleFunc("/CustomFlights", customFlightsHandler(db))

	// Static files
	fs := http.FileServer(http.Dir("../templates/styles"))

	mux.Handle("/styles/", http.StripPrefix("/styles", fs))
	return mux
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(getTemplatePath("home.html")))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// getTemplatePath returns the absolute path to a file in the templates directory
func getTemplatePath(filename string) string {
	return filepath.Join("..", "..", "cmd", "templates", filename)
}

func flightsHandler(db *sql.DB, gateMin, gateMax int) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(getTemplatePath("index.html")))

	return func(w http.ResponseWriter, r *http.Request) {
		// Get flights data
		flights, err := services.GetFlightsToday(db, gateMin, gateMax)
		if err != nil {
			http.Error(w, "Failed to retrieve flights", http.StatusInternalServerError)
			log.Printf("Failed to get flights: %v", err)
			return
		}

		// Execute template with data
		err = tmpl.Execute(w, flights)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Printf("Template execution error: %v", err)
		}
	}
}

func customFlightsHandler(db *sql.DB) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(getTemplatePath("index.html")))

	return func(w http.ResponseWriter, r *http.Request) {
		gate1 := r.URL.Query().Get("gate1")
		gate2 := r.URL.Query().Get("gate2")

		g1, err := strconv.Atoi(gate1)
		if err != nil {
			http.Error(w, "Invalid gate1 value", http.StatusBadRequest)
			return
		}

		g2, err := strconv.Atoi(gate2)
		if err != nil {
			http.Error(w, "Invalid gate2 value", http.StatusBadRequest)
		}

		// Validation: Ensure gate1 is lower than gate2
		if g1 >= g2 {
			http.Error(w, "Gate 1 must be lower than Gate 2. Please enter valid values.", http.StatusBadRequest)
			log.Println("Gate validation failed: gate1 must be lower than gate2")
			return
		}

		flights, err := services.GetFlightsToday(db, g1, g2)
		if err != nil {
			http.Error(w, "Failed to retrieve flights", http.StatusInternalServerError)
			log.Printf("Failed to get flights: %v", err)
			return
		}

		err = tmpl.Execute(w, flights)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Printf("Template execution error: %v", err)
		}
	}
}
