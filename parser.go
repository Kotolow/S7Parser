package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"time"
)

type reqT interface {
	map[string]string | RequestBody
}

func fetchData[T reqT](url string, body T, client *http.Client, convId string) []byte {
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Impossible to marshall body:", err)
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalledBody))
	if err != nil {
		fmt.Println("Post error:", err)
		return nil
	}
	request.Header.Set("x-locale", "ru")
	request.Header.Set("content-type", "application/json; charset=UTF-8")
	request.Header.Set("x-application", "ibe")
	request.Header.Set("x-application-version", "5.3.5")
	request.Header.Set("x-platform", "android")
	request.Header.Set("authorization", "Basic YW5hcHA6U0VxNXBYVkFaenVWaHpIQWJVZ3VGM3ZZ")
	if convId != "" {
		request.Header.Set("x-conversation", convId)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error in Do():", err)
		return nil
	}
	defer response.Body.Close()
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Impossible to read all body of response:", err)
		return nil
	}
	return resBody
}

func updateConversationID(client *http.Client) string {
	url := "https://mproxy.api.s7.ru/3.0/conversations"

	body := map[string]string{"ip": "127.0.0.1"}

	resBody := fetchData(url, body, client, "")

	id := gjson.Get(string(resBody), "conversation.id")

	return id.String()
}

func getFlightInfo(convID, date, destination, origin string, client *http.Client) {
	url := "https://mproxy.api.s7.ru/3.0/flights/search"
	passengers := passengersAmount{
		Adults:   1,
		Children: 0,
		Infants:  0,
	}

	var routes []route

	routes = append(routes, route{
		DepartureDate: date,
		Destination:   destination,
		Origin:        origin,
	})

	params := searchParams{
		Currency:                 "RUB",
		DirectOnly:               false,
		OwnAirlineOnly:           false,
		PassengersAmount:         passengers,
		PromoCode:                "",
		Redemption:               false,
		Routes:                   routes,
		SearchType:               "EXACT",
		SubsidizedPassengerTypes: nil,
		TripType:                 "ONE_WAY",
	}

	body := RequestBody{SearchParams: params}

	resBody := fetchData(url, body, client, convID)

	optionSets := gjson.Get(string(resBody), "search.optionSets")
	for _, optionSet := range optionSets.Array() {
		options := gjson.Get(optionSet.String(), "options")
		for _, option := range options.Array() {
			tRoute := gjson.Get(option.String(), "routes").Array()[0]
			solution := gjson.Get(option.String(), "solutions.BASICECONOMY.pricing")
			var info FlightInfo

			segments := gjson.Get(tRoute.String(), "segments").Array()
			firstSegment := segments[0]
			lastSegment := segments[len(segments)-1]

			basePrice := gjson.Get(solution.String(), "base.price.amount").Float()
			fees := gjson.Get(solution.String(), "fees.price.amount").Float()
			taxes := gjson.Get(solution.String(), "taxes.price.amount").Float()

			flightCode := gjson.Get(firstSegment.String(), "operatingAirline.displayCode").String()
			flightNumber := gjson.Get(firstSegment.String(), "operatingAirline.flightNumber").String()

			currentTime := time.Now()

			info.Carrier = gjson.Get(firstSegment.String(), "displayAirlineCode").String()
			info.ObservationDate = currentTime.Format("2006-01-02")
			info.ObservationTime = currentTime.Format("15:04:05")
			info.Origin = gjson.Get(tRoute.String(), "origin").String()
			info.Destination = gjson.Get(tRoute.String(), "destination").String()
			info.IsOneWay = "true"
			info.OutboundFlightNo = flightCode + " " + flightNumber
			info.OutboundDepartureDate = gjson.Get(firstSegment.String(), "departureDate").String()
			info.OutboundArrivalDate = gjson.Get(lastSegment.String(), "arrivalDate").String()
			info.PriceExc = basePrice
			info.Tax = fees + taxes
			info.Currency = gjson.Get(solution.String(), "total.price.currency").String()
			info.AircraftCode = gjson.Get(firstSegment.String(), "aircraft.code").String()
			info.AircraftName = gjson.Get(firstSegment.String(), "aircraft.name").String()
			if len(segments) == 2 {
				stop := gjson.Get(tRoute.String(), "stops").Array()[0]
				info.TransferIata = gjson.Get(firstSegment.String(), "arrivalAirport.code").String()
				info.TransferDuration = gjson.Get(stop.String(), "duration.amount").Int()
			} else {
				info.TransferIata = "-"
				info.TransferDuration = 0
			}
			_, err := addFlight(&info)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
