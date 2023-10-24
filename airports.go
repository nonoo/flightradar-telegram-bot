package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const airportUpdateTimeout = 3 * time.Minute

type Airport struct {
	Name      string      `json:"name"`
	IATA      string      `json:"iata"`
	ICAO      string      `json:"icao"`
	Latitude  float64     `json:"lat"`
	Longitude float64     `json:"lon"`
	Country   string      `json:"country"`
	Altitude  interface{} `json:"alt"`
}

type Airports struct {
	Airports []Airport `json:"rows"`
}

var airports Airports

func (a *Airports) Find(iata string) *Airport {
	for i := range a.Airports {
		if a.Airports[i].IATA == iata {
			return &a.Airports[i]
		}
	}
	return nil
}

func (a *Airports) Load(ctx context.Context) error {
	fmt.Println("getting airport info")
	reqCtx, cancel := context.WithTimeout(ctx, airportUpdateTimeout)
	rawServerReply, err := httpReq(reqCtx, "https://www.flightradar24.com/_json/airports.php", nil)
	cancel()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(rawServerReply), a)
}
