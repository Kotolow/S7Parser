package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}
	err := initDatabase("/mnt/c/SQLite/s7.db")
	defer closeDatabase()
	if err != nil {
		fmt.Println(err)
	}
	convId := updateConversationID(client)

	var dates []string
	layout := "2006-01-02"
	startDate := time.Now().Format(layout)
	t, err := time.Parse(layout, startDate)

	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 90; i++ {
		dates = append(dates, t.Format(layout))
		t = t.AddDate(0, 0, 1)
	}
	directions := map[string]string{
		"MOW": "MMK",
		"AER": "GOJ",
	}
	for i, v := range directions {
		for _, date := range dates {
			getFlightInfo(convId, date, i, v, client)
			getFlightInfo(convId, date, v, i, client)
		}
	}

	err = exportToCSV("results/flights.csv")
	if err != nil {
		fmt.Printf("Error exporting to CSV: %v\n", err)
	} else {
		fmt.Println("Data successfully exported to flights.csv")
	}

	err = exportFilteredFlightsToCSV("results/GOJ-AER.csv", "GOJ", "AER")
	if err != nil {
		fmt.Printf("Error exporting filtered flights to CSV: %v\n", err)
	} else {
		fmt.Println("Filtered data successfully exported to filtered_flights.csv")
	}
	fmt.Println("Finished in ", time.Now())
}
