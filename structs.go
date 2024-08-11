package main

type searchParams struct {
	Currency                 string           `json:"currency"`
	DirectOnly               bool             `json:"directOnly"`
	OwnAirlineOnly           bool             `json:"ownAirlineOnly"`
	PassengersAmount         passengersAmount `json:"passengersAmount"`
	PromoCode                string           `json:"promoCode"`
	Redemption               bool             `json:"redemption"`
	Routes                   []route          `json:"routes"`
	SearchType               string           `json:"searchType"`
	SubsidizedPassengerTypes []string         `json:"subsidizedPassengerTypes"`
	TripType                 string           `json:"tripType"`
}

type route struct {
	DepartureDate string `json:"departureDate"`
	Destination   string `json:"destination"`
	Origin        string `json:"origin"`
}

type passengersAmount struct {
	Adults   int `json:"adults"`
	Children int `json:"children"`
	Infants  int `json:"infants"`
}

type RequestBody struct {
	SearchParams searchParams `json:"searchParams"`
}

type FlightInfo struct {
	Carrier               string
	ObservationDate       string
	ObservationTime       string
	Origin                string
	Destination           string
	IsOneWay              string
	OutboundFlightNo      string
	OutboundDepartureDate string
	OutboundArrivalDate   string
	PriceExc              float64
	Tax                   float64
	Currency              string
	AircraftCode          string
	AircraftName          string
	TransferIata          string
	TransferDuration      int64
}
