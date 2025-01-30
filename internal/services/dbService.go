package services

import (
	"database/sql"
	"fmt"
	"restoportGo/internal/utils"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func EstablishConnection() (*sql.DB, error) {
	env := utils.NewEnv()
	db, err := sql.Open("mysql", env.DSN)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}

	return db, nil
}

func GetFlightsToday(db *sql.DB, minGate, maxGate int) ([]utils.Flight, error) {
	querry := "SELECT * FROM todays_flights WHERE gate BETWEEN ? AND ?"

	rows, err := db.Query(querry, minGate, maxGate)
	if err != nil {
		return nil, fmt.Errorf("failed to query flights: %w", err)
	}
	defer rows.Close()

	var flights []utils.Flight

	for rows.Next() {
		var flight utils.Flight
		err := rows.Scan(
			&flight.AirlineName,
			&flight.Gate,
			&flight.Time,
			&flight.UpdatedTime,
			&flight.Destination,
			&flight.Status,
			&flight.UniqueDisplayNo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		flights = append(flights, flight)
	}
	return flights, nil
}
