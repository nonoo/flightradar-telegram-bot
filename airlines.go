package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const airlineUpdateTimeout = 3 * time.Minute

type Airline struct {
	Name string
	Code string
	ICAO string
}

type Airlines struct {
	Airlines []Airline `json:"rows"`
}

var airlines Airlines

func (a *Airlines) Find(flightNr string) *Airline {
	for i := range a.Airlines {
		if a.Airlines[i].Code != "" && strings.HasPrefix(flightNr, a.Airlines[i].Code) {
			return &a.Airlines[i]
		}
	}
	return nil
}

func (a *Airlines) Load(ctx context.Context) error {
	fmt.Println("getting airline info")
	reqCtx, cancel := context.WithTimeout(ctx, airlineUpdateTimeout)
	rawServerReply, err := httpReq(reqCtx, "https://www.flightradar24.com/_json/airlines.php", nil)
	cancel()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(rawServerReply), a)
}
