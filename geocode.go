package main

import (
	geocoder "github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	geo "github.com/kellydunn/golang-geo"
)

func GetCoordinatesFromAddress(address string) (location *geocoder.Location, err error) {
	g := openstreetmap.Geocoder()
	return g.Geocode(address)
}

func GetRectCoordinatesFromLocation(location *geocoder.Location, rangeKm int) (p1 geocoder.Location, p2 geocoder.Location) {
	p := geo.NewPoint(location.Lat, location.Lng)
	rp1 := p.PointAtDistanceAndBearing(float64(rangeKm), -45) // Top left point of the rectangle.
	p1 = geocoder.Location{Lat: rp1.Lat(), Lng: rp1.Lng()}
	rp2 := p.PointAtDistanceAndBearing(float64(rangeKm), 135) // Bottom right point of the rectangle.
	p2 = geocoder.Location{Lat: rp2.Lat(), Lng: rp2.Lng()}
	return
}

func GetDistanceInKm(l1 *geocoder.Location, l2 *geocoder.Location) int {
	p1 := geo.NewPoint(l1.Lat, l1.Lng)
	p2 := geo.NewPoint(l2.Lat, l2.Lng)
	return int(p1.GreatCircleDistance(p2) / 1000)
}
