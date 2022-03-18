package main

import (
	"errors"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"time"
)

type CountryVisit struct {
	Country string
	Start   time.Time
	End     time.Time
}

func AnalyzeHistory(fc geojson.FeatureCollection) []CountryVisit {
	// reading files from disk is much faster than analyzing them so increasing the buffer size does not help
	visits := make(chan PlaceVisit, 500)

	go WriteLocationHistory("./data/history/", visits)

	var countryVisits []CountryVisit

	var point *geom.Point
	for visit := range visits {
		point = geom.NewPointFlat(geom.XY,
			[]float64{
				float64(visit.Location.Lon) / 1e7,
				float64(visit.Location.Lat) / 1e7,
			})
		feature, err := FindSurroundingFeature(&fc, point)
		if err == nil {
			countryVisits = updateVisits(countryVisits, feature, visit)
		}
	}
	return countryVisits
}

func updateVisits(countryVisits []CountryVisit, feature *geojson.Feature, visit PlaceVisit) []CountryVisit {
	country := fmt.Sprintf("%v", feature.Properties["ADMIN"])
	newVisit := CountryVisit{country, visit.Duration.Start, visit.Duration.End}

	if len(countryVisits) == 0 {
		countryVisits = append(countryVisits, newVisit)
	} else if lastVisit := &countryVisits[len(countryVisits)-1]; lastVisit.Country != country {
		lastVisit.End = newVisit.Start
		countryVisits = append(countryVisits, newVisit)
	}
	return countryVisits
}

func FindSurroundingFeature(fc *geojson.FeatureCollection, point *geom.Point) (*geojson.Feature, error) {
	idx := -1
	for i, feature := range fc.Features {
		if b, err := IsInGeometry(feature.Geometry, point); b && err == nil {
			idx = i
		}
	}

	if idx == -1 {
		return nil, errors.New("could not find matching point")
	}

	// small optimization hack. Put the matching country to the beginning of the slice
	// as some kind of primitive cache
	newFeatureSlice := make([]*geojson.Feature, 0, len(fc.Features))

	newFeatureSlice = append(newFeatureSlice, fc.Features[idx])
	newFeatureSlice = append(newFeatureSlice, fc.Features[:idx]...)
	newFeatureSlice = append(newFeatureSlice, fc.Features[idx+1:]...)

	fc.Features = newFeatureSlice

	return newFeatureSlice[0], nil
}

func IsInGeometry(geometry geom.T, point *geom.Point) (bool, error) {
	switch g := geometry.(type) {
	case *geom.Polygon:
		return isInPolygon(g, point), nil
	case *geom.MultiPolygon:
		return isInMultiPolygon(g, point), nil
	default:
		return false, errors.New("unknown geometry type")
	}
}

func isInMultiPolygon(multiPolygon *geom.MultiPolygon, point *geom.Point) bool {
	for i := 0; i < multiPolygon.NumPolygons(); i++ {
		if isInPolygon(multiPolygon.Polygon(i), point) {
			return true
		}
	}
	return false
}

func isInPolygon(polygon *geom.Polygon, point *geom.Point) bool {
	// a polygon consists of multiple rings and
	// every other ring describes a hole in the polygon
	for i := 0; i < polygon.NumLinearRings(); i++ {
		isInRing := isInLinearRing(polygon.LinearRing(i), point)
		if i%2 == 0 && !isInRing {
			return false
		} else if i%2 == 1 && !isInRing {
			return true
		}
	}
	return polygon.NumLinearRings()%2 == 1
}

//isInLinearRing follows the algorithm laid out in "Numerical Recipes" Chapter 21.4.3 Polygons. See also http://numerical.recipes/book/book.html
func isInLinearRing(ring *geom.LinearRing, point *geom.Point) bool {
	wind := 0
	p := ring.Coord(ring.NumCoords() - 1) // last point is previous to first

	for _, v := range ring.Coords() {
		if p.Y() < point.Y() {
			if v.Y() > point.Y() && (p.X()-point.X())*(v.Y()-point.Y())-(p.Y()-point.Y())*(v.X()-point.X()) > 0 {
				wind++
			}
		} else {
			if v.Y() <= point.Y() && (p.X()-point.X())*(v.Y()-point.Y())-(p.Y()-point.Y())*(v.X()-point.X()) < 0 {
				wind--
			}
		}
		p = v
	}
	return wind != 0
}
