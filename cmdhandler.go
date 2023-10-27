package main

import (
	"context"
	"fmt"
	"strconv"

	geocoder "github.com/codingsince1985/geo-golang"
	"github.com/go-telegram/bot/models"
	"golang.org/x/exp/slices"
)

type cmdHandlerType struct{}

func (c *cmdHandlerType) getLocationDescription(chatID int64) string {
	name := settings.GetString(chatID, "LocationName")
	rangeKm := settings.GetInt(chatID, "LocationRangeKm")
	loc := geocoder.Location{Lat: settings.GetFloat64(chatID, "LocationLat"), Lng: settings.GetFloat64(chatID, "LocationLng")}
	if name == "" || (loc.Lat == 0 && loc.Lng == 0) {
		return "📌 No location set"
	}
	if rangeKm == 0 {
		return "📏 No range set"
	}

	p1, p2 := GetRectCoordinatesFromLocation(&loc, rangeKm)
	return "📌 Current location is " + name + " https://www.google.com/maps/place/" + fmt.Sprint(loc.Lat, ",", loc.Lng) + "\n" +
		"📏 Range: " + fmt.Sprint(rangeKm) + " km\n" +
		"🗺 Area top left: https://www.google.com/maps/place/" + fmt.Sprint(p1.Lat, ",", p1.Lng) + "\n" +
		"🗺 Area bottom right: https://www.google.com/maps/place/" + fmt.Sprint(p2.Lat, ",", p2.Lng)
}

func (c *cmdHandlerType) Location(ctx context.Context, msg *models.Message) {
	if msg.Text == "" {
		sendReplyToMessage(ctx, msg, c.getLocationDescription(msg.Chat.ID))
		return
	}

	if !slices.Contains(params.AllowedUserIDs, msg.From.ID) {
		fmt.Println("  user not allowed")
		return
	}

	loc, err := GetCoordinatesFromAddress(msg.Text)
	if err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}
	if loc == nil {
		sendReplyToMessage(ctx, msg, errorStr+": unknown address")
		return
	}

	if err := settings.Set(msg.Chat.ID, "LocationName", msg.Text); err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}
	if err := settings.Set(msg.Chat.ID, "LocationLat", loc.Lat); err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}
	if err := settings.Set(msg.Chat.ID, "LocationLng", loc.Lng); err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}

	sendReplyToMessage(ctx, msg, "✅ Location set!\n"+c.getLocationDescription(msg.Chat.ID))
}

func (c *cmdHandlerType) Range(ctx context.Context, msg *models.Message) {
	if msg.Text == "" {
		sendReplyToMessage(ctx, msg, c.getLocationDescription(msg.Chat.ID))
		return
	}

	if !slices.Contains(params.AllowedUserIDs, msg.From.ID) {
		fmt.Println("  user not allowed")
		return
	}

	r, err := strconv.Atoi(msg.Text)
	if err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}

	if err := settings.Set(msg.Chat.ID, "LocationRangeKm", r); err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}

	sendReplyToMessage(ctx, msg, "📏 Range set!\n"+c.getLocationDescription(msg.Chat.ID))
}

func (c *cmdHandlerType) MinRange(ctx context.Context, msg *models.Message) {
	if msg.Text == "" {
		minRangeKm := settings.GetInt(msg.Chat.ID, "MinimumRangeKm")
		sendReplyToMessage(ctx, msg, "📏 Minimum range: "+fmt.Sprint(minRangeKm)+" km\n")
		return
	}

	if !slices.Contains(params.AllowedUserIDs, msg.From.ID) {
		fmt.Println("  user not allowed")
		return
	}

	r, err := strconv.Atoi(msg.Text)
	if err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}

	if err := settings.Set(msg.Chat.ID, "MinimumRangeKm", r); err != nil {
		sendReplyToMessage(ctx, msg, errorStr+": "+err.Error())
		return
	}

	sendReplyToMessage(ctx, msg, "📏 Minimum range set to "+fmt.Sprint(r)+" km!")
}

func (c *cmdHandlerType) Help(ctx context.Context, msg *models.Message, cmdChar string) {
	sendReplyToMessage(ctx, msg, "🤖 Flightradar Telegram Bot\n\n"+
		"Available commands:\n\n"+
		cmdChar+"frloc (location) - set or show current location\n"+
		cmdChar+"frrange (range) - set or show current range\n"+
		cmdChar+"frminrange (range) - set or show current minimum range\n"+
		cmdChar+"frhelp - show this help\n\n"+
		"For more information see https://github.com/nonoo/flightradar-telegram-bot")
}
