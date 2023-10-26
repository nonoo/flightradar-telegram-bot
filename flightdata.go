package main

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	geocoder "github.com/codingsince1985/geo-golang"
	"github.com/go-telegram/bot"
	"golang.org/x/exp/slices"
)

type FlightDataRaw struct {
	Aircraft [][]interface{} `json:"aircraft"`
}

func (f *FlightDataRaw) Parse(s string) error {
	if err := json.Unmarshal([]byte(s), f); err != nil {
		return err
	}
	return nil
}

// Example: ["32903253","4400EC",47.7039,17.1782,323,36025,415,"","F-LROD3","A320","OE-IVU",1698151578,"SKG","BER","U25030",0,-128,"EJU62QN",0]
type FlightDataAircraft struct {
	ID           string
	FlightID     string
	Latitude     float64
	Longitude    float64
	Heading      int
	Altitude     int
	Speed        int
	Type         string
	RegNumber    string
	Origin       string
	Destination  string
	FlightNumber string
}

type FlightDataAircrafts []FlightDataAircraft

func (f *FlightDataAircrafts) Contains(a FlightDataAircraft) bool {
	for i := range *f {
		if (*f)[i].ID == a.ID {
			return true
		}
	}
	return false
}

func (f *FlightDataAircrafts) Parse(fdr FlightDataRaw) {
	*f = FlightDataAircrafts{}

	for i := range fdr.Aircraft {
		if len(fdr.Aircraft[i]) < 15 {
			fmt.Print("  error: aircraft #", i, " has less than 15 fields: ", fdr.Aircraft[i], "\n")
			continue
		}

		ID, ok := fdr.Aircraft[i][0].(string)
		if !ok {
			fmt.Print("  error parsing ID of aircraft #", i, ": ", fdr.Aircraft[i][0], "\n")
			continue
		}
		FlightID, ok := fdr.Aircraft[i][1].(string)
		if !ok {
			fmt.Print("  error parsing FlightID of aircraft #", i, ": ", fdr.Aircraft[i][1], "\n")
			continue
		}
		Latitude, ok := fdr.Aircraft[i][2].(float64)
		if !ok {
			fmt.Print("  error parsing Latitude of aircraft #", i, ": ", fdr.Aircraft[i][2], "\n")
			continue
		}
		Longitude, ok := fdr.Aircraft[i][3].(float64)
		if !ok {
			fmt.Print("  error parsing Longitude of aircraft #", i, ": ", fdr.Aircraft[i][3], "\n")
			continue
		}
		Heading, ok := fdr.Aircraft[i][4].(float64)
		if !ok {
			fmt.Print("  error parsing Heading of aircraft #", i, ": ", fdr.Aircraft[i][4], "\n")
			continue
		}
		Altitude, ok := fdr.Aircraft[i][5].(float64)
		if !ok {
			fmt.Print("  error parsing Altitude of aircraft #", i, ": ", fdr.Aircraft[i][5], "\n")
			continue
		}
		Speed, ok := fdr.Aircraft[i][6].(float64)
		if !ok {
			fmt.Print("  error parsing Speed of aircraft #", i, ": ", fdr.Aircraft[i][6], "\n")
			continue
		}
		Type, ok := fdr.Aircraft[i][9].(string)
		if !ok {
			fmt.Print("  error parsing Type of aircraft #", i, ": ", fdr.Aircraft[i][9], "\n")
			continue
		}
		RegNumber, ok := fdr.Aircraft[i][10].(string)
		if !ok {
			fmt.Print("  error parsing RegNumber of aircraft #", i, ": ", fdr.Aircraft[i][10], "\n")
			continue
		}
		Origin, ok := fdr.Aircraft[i][12].(string)
		if !ok {
			fmt.Print("  error parsing Origin of aircraft #", i, ": ", fdr.Aircraft[i][12], "\n")
			continue
		}
		Destination, ok := fdr.Aircraft[i][13].(string)
		if !ok {
			fmt.Print("  error parsing Destination of aircraft #", i, ": ", fdr.Aircraft[i][13], "\n")
			continue
		}
		FlightNumber, ok := fdr.Aircraft[i][14].(string)
		if !ok {
			fmt.Print("  error parsing FlightNumber of aircraft #", i, ": ", fdr.Aircraft[i][14], "\n")
			continue
		}

		*f = append(*f, FlightDataAircraft{
			ID:           ID,
			FlightID:     FlightID,
			Latitude:     Latitude,
			Longitude:    Longitude,
			Heading:      int(Heading),
			Altitude:     int(Altitude),
			Speed:        int(Speed),
			Type:         Type,
			RegNumber:    RegNumber,
			Origin:       Origin,
			Destination:  Destination,
			FlightNumber: FlightNumber,
		})
	}
}

type FlightDataLocation struct {
	ChatIDs     []int64
	TopLeft     geocoder.Location
	BottomRight geocoder.Location

	Aircrafts FlightDataAircrafts
}

type FlightData struct {
	Location map[string]FlightDataLocation
}

const flightDataUpdateInterval = time.Minute
const flightDataUpdateTimeout = time.Second * 30

var flightData FlightData

func (f *FlightData) Updater(ctx context.Context) {
	f.Location = map[string]FlightDataLocation{}

	for {
		chatIDs := settings.GetChatIDs()
		for _, id := range chatIDs {
			_, err := telegramBot.GetChat(ctx, &bot.GetChatParams{
				ChatID: id,
			})
			if err != nil { // We are not in this chat anymore?
				_ = settings.RemoveChatID(id)
				continue
			}

			name := settings.GetString(id, "LocationName")
			rangeKm := settings.GetInt(id, "LocationRangeKm")
			loc := geocoder.Location{Lat: settings.GetFloat64(id, "LocationLat"), Lng: settings.GetFloat64(id, "LocationLng")}
			if name == "" || rangeKm == 0 || (loc.Lat == 0 && loc.Lng == 0) { // No location is set in this chat?
				continue
			}

			if _, ok := f.Location[name]; ok { // We already have info about this location?
				if !slices.Contains(f.Location[name].ChatIDs, id) { // This chat ID is not in the list of this location?
					l := f.Location[name]
					l.ChatIDs = append(f.Location[name].ChatIDs, id)
					f.Location[name] = l
				}
				continue
			}

			// Storing info about this location.
			p1, p2 := GetRectCoordinatesFromLocation(&loc, rangeKm)
			l := FlightDataLocation{
				ChatIDs:     []int64{id},
				TopLeft:     p1,
				BottomRight: p2,
			}
			f.Location[name] = l
		}

		// Iterating through stored locations.
		for i := range f.Location {
			fmt.Println("getting info for location:", i)
			reqCtx, cancel := context.WithTimeout(ctx, flightDataUpdateTimeout)
			url := fmt.Sprintf("https://data-cloud.flightradar24.com/zones/fcgi/feed.js?faa=1&satellite=1&mlat=1&flarm=1&adsb=1"+
				"&gnd=1&air=1&vehicles=1&estimated=1&gliders=1&stats=1&limit=5000&array=1&bounds=%.2f%%2C%.2f%%2C%.2f%%2C%.2f",
				f.Location[i].TopLeft.Lat, f.Location[i].BottomRight.Lat, f.Location[i].TopLeft.Lng, f.Location[i].BottomRight.Lng)
			fmt.Println("  url:", url)
			l := f.Location[i]
			rawServerReply, err := httpReq(reqCtx, url, nil)
			cancel()
			if err != nil {
				fmt.Println("  error:", err)
				continue
			}
			// fmt.Println("  got data:", rawServerReply)
			var fdr FlightDataRaw
			err = fdr.Parse(rawServerReply)
			if err != nil {
				fmt.Println("  error:", err)
				continue
			}

			newAircrafts := FlightDataAircrafts{}
			newAircrafts.Parse(fdr)

			if l.Aircrafts != nil {
				for _, newAircraft := range newAircrafts {
					if l.Aircrafts.Contains(newAircraft) {
						continue
					}

					var flightNr string
					if newAircraft.FlightNumber != "" {
						flightNr = newAircraft.FlightNumber
						airline := airlines.Find(newAircraft.FlightNumber)
						if airline != nil {
							flightNr += " (" + airline.Name + ")"
						}
					} else {
						flightNr = "N/A"
					}

					origin := "N/A"
					if newAircraft.Origin != "" {
						origin = newAircraft.Origin
						if a := airports.Find(newAircraft.Origin); a != nil {
							origin += " (" + a.Country + ")"
						}
					}
					dest := "N/A"
					if newAircraft.Destination != "" {
						dest = newAircraft.Destination
						if a := airports.Find(newAircraft.Destination); a != nil {
							dest += " (" + a.Country + ")"
						}
					}

					if newAircraft.FlightNumber == "" && origin == "N/A" && dest == "N/A" {
						continue
					}

					var origDest string
					if origin != "N/A" && dest != "N/A" {
						origDest = "🗺 " + origin + " → " + dest + "\n"
					}

					var aircraftType string
					if newAircraft.Type != "" {
						aircraftType = "⚙ " + newAircraft.Type + "\n"
					}

					msg := fmt.Sprintf("🛩️ " + flightNr + "\n" +
						origDest +
						aircraftType +
						path.Join("📍 https://www.flightradar24.com/", newAircraft.FlightNumber, newAircraft.ID))
					fmt.Println("  new aircraft:", msg)

					for _, chatID := range f.Location[i].ChatIDs {
						fmt.Println("    sending to chat:", chatID)
						_, err = telegramBot.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: chatID,
							Text:   msg,
						})
						if err != nil {
							fmt.Println("      error:", err)
						}
					}
				}
			}
			l.Aircrafts = newAircrafts
			f.Location[i] = l
		}

		time.Sleep(flightDataUpdateInterval)
	}
}

func (f *FlightData) Init(ctx context.Context) {
	go f.Updater(ctx)
}
