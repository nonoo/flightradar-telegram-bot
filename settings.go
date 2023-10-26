package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

type ChatSettings struct {
	Location struct {
		Name           string  `json:"name,omitempty"`
		Lat            float64 `json:"lat,omitempty"`
		Lng            float64 `json:"lng,omitempty"`
		RangeKm        int     `json:"rangeKm,omitempty"`
		MinimumRangeKm int     `json:"minimumRangeKm,omitempty"`
	} `json:"location,omitempty"`
}

type Settings struct {
	Chat map[int64]ChatSettings `json:"ids"`
}

var configPath = configdir.LocalConfig("flightradar-telegram-bot")
var configFilePath = filepath.Join(configPath, "settings.json")
var settings Settings

func (s *Settings) Load() error {
	data, err := os.ReadFile(configFilePath)
	if err == nil {
		err = json.Unmarshal(data, &settings)
		if err != nil {
			return err
		}
	}

	// Initializing ChatSettings ID map if needed.
	if m := settings.Chat; m == nil {
		m = map[int64]ChatSettings{}
		for _, id := range params.AllowedGroupIDs {
			m[id] = ChatSettings{}
		}
		settings.Chat = m
	}
	return settings.save()
}

func (s *Settings) save() error {
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}

	err = configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *Settings) Set(chatID int64, key string, value interface{}) error {
	cs := settings.Chat[chatID]

	switch key {
	case "LocationName":
		cs.Location.Name = value.(string)
	case "LocationLat":
		cs.Location.Lat = value.(float64)
	case "LocationLng":
		cs.Location.Lng = value.(float64)
	case "LocationRangeKm":
		cs.Location.RangeKm = value.(int)
	case "MinimumRangeKm":
		cs.Location.MinimumRangeKm = value.(int)
	default:
		return fmt.Errorf("unknown setting key: %s", key)
	}

	settings.Chat[chatID] = cs
	return settings.save()
}

func (s *Settings) GetString(chatID int64, key string) string {
	cs := settings.Chat[chatID]

	switch key {
	case "LocationName":
		return cs.Location.Name
	default:
		panic("invalid setting key")
	}
}

func (s *Settings) GetFloat64(chatID int64, key string) float64 {
	cs := settings.Chat[chatID]

	switch key {
	case "LocationLat":
		return cs.Location.Lat
	case "LocationLng":
		return cs.Location.Lng
	default:
		panic("invalid setting key")
	}
}

func (s *Settings) GetInt(chatID int64, key string) int {
	cs := settings.Chat[chatID]

	switch key {
	case "LocationRangeKm":
		return cs.Location.RangeKm
	case "MinimumRangeKm":
		return cs.Location.MinimumRangeKm
	default:
		panic("invalid setting key")
	}
}

func (s *Settings) GetChatIDs() []int64 {
	chatIDs := []int64{}
	for id := range settings.Chat {
		chatIDs = append(chatIDs, id)
	}
	return chatIDs
}

func (s *Settings) RemoveChatID(chatID int64) error {
	delete(settings.Chat, chatID)
	return settings.save()
}
