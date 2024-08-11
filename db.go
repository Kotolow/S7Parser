package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

var db *sql.DB

func initDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {

	}
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS flights (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			carrier TEXT NOT NULL,
			observationDate TEXT NOT NULL,
			observationTime TEXT NOT NULL,
			origin TEXT NOT NULL,
			destination TEXT NOT NULL,
			isOneWay TEXT NOT NULL,
			outboundFlightNo TEXT NOT NULL,
			outboundDepartureDate TEXT NOT NULL,
			outboundArrivalDate TEXT NOT NULL,
			priceExc REAL NOT NULL,
			tax REAL NOT NULL,
			currency TEXT NOT NULL,
			aircraftCode TEXT NOT NULL,
			aircraftName TEXT NOT NULL,
			transferIata TEXT NOT NULL,
			transferDuration INTEGER NOT NULL
		)`,
	)
	if err != nil {
		return err
	}
	return nil
}

func addFlight(info *FlightInfo) (int64, error) {
	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO flights
   			(carrier, observationDate, observationTime, origin, destination, isOneWay,
   			 outboundFlightNo, outboundDepartureDate, outboundArrivalDate, priceExc,
   			 tax, currency, aircraftCode, aircraftName, transferIata, transferDuration) VALUES
   			(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`, info.Carrier, info.ObservationDate, info.ObservationTime,
		info.Origin, info.Destination, info.IsOneWay, info.OutboundFlightNo, info.OutboundDepartureDate,
		info.OutboundArrivalDate, info.PriceExc, info.Tax, info.Currency, info.AircraftCode, info.AircraftName,
		info.TransferIata, info.TransferDuration)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func closeDatabase() {
	err := db.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func exportToCSV(filePath string) error {
	rows, err := db.QueryContext(context.Background(), `SELECT 
		carrier, observationDate, observationTime, origin, destination, isOneWay, 
		outboundFlightNo, outboundDepartureDate, outboundArrivalDate, priceExc, 
		tax, currency, aircraftCode, aircraftName, transferIata, transferDuration 
		FROM flights`)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Carrier", "Observation Date", "Observation Time", "Origin", "Destination", "Is One Way",
		"Outbound Flight No", "Outbound Departure Date", "Outbound Arrival Date", "Price Excluding Tax",
		"Tax", "Currency", "Aircraft Code", "Aircraft Name", "Transfer IATA", "Transfer Duration",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header to CSV: %v", err)
	}

	for rows.Next() {
		var flight FlightInfo
		if err := rows.Scan(&flight.Carrier, &flight.ObservationDate, &flight.ObservationTime,
			&flight.Origin, &flight.Destination, &flight.IsOneWay, &flight.OutboundFlightNo,
			&flight.OutboundDepartureDate, &flight.OutboundArrivalDate, &flight.PriceExc,
			&flight.Tax, &flight.Currency, &flight.AircraftCode, &flight.AircraftName,
			&flight.TransferIata, &flight.TransferDuration); err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		record := []string{
			flight.Carrier, flight.ObservationDate, flight.ObservationTime, flight.Origin,
			flight.Destination, flight.IsOneWay, flight.OutboundFlightNo, flight.OutboundDepartureDate,
			flight.OutboundArrivalDate, strconv.FormatFloat(flight.PriceExc, 'f', 2, 64),
			strconv.FormatFloat(flight.Tax, 'f', 2, 64), flight.Currency, flight.AircraftCode,
			flight.AircraftName, flight.TransferIata, strconv.FormatInt(flight.TransferDuration, 10),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error with rows: %v", err)
	}

	return nil
}

func exportFilteredFlightsToCSV(filePath string, origin, destination string) error {
	rows, err := db.QueryContext(context.Background(), `SELECT 
		carrier, observationDate, observationTime, origin, destination, isOneWay, 
		outboundFlightNo, outboundDepartureDate, outboundArrivalDate, priceExc, 
		tax, currency, aircraftCode, aircraftName, transferIata, transferDuration 
		FROM flights WHERE origin = ? AND destination = ?`, origin, destination)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Carrier", "Observation Date", "Observation Time", "Origin", "Destination", "Is One Way",
		"Outbound Flight No", "Outbound Departure Date", "Outbound Arrival Date", "Price Excluding Tax",
		"Tax", "Currency", "Aircraft Code", "Aircraft Name", "Transfer IATA", "Transfer Duration",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header to CSV: %v", err)
	}

	for rows.Next() {
		var flight FlightInfo
		if err := rows.Scan(&flight.Carrier, &flight.ObservationDate, &flight.ObservationTime,
			&flight.Origin, &flight.Destination, &flight.IsOneWay, &flight.OutboundFlightNo,
			&flight.OutboundDepartureDate, &flight.OutboundArrivalDate, &flight.PriceExc,
			&flight.Tax, &flight.Currency, &flight.AircraftCode, &flight.AircraftName,
			&flight.TransferIata, &flight.TransferDuration); err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		record := []string{
			flight.Carrier, flight.ObservationDate, flight.ObservationTime, flight.Origin,
			flight.Destination, flight.IsOneWay, flight.OutboundFlightNo, flight.OutboundDepartureDate,
			flight.OutboundArrivalDate, strconv.FormatFloat(flight.PriceExc, 'f', 2, 64),
			strconv.FormatFloat(flight.Tax, 'f', 2, 64), flight.Currency, flight.AircraftCode,
			flight.AircraftName, flight.TransferIata, strconv.FormatInt(flight.TransferDuration, 10),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error with rows: %v", err)
	}

	return nil
}
