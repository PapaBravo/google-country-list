package main

import (
	"github.com/twpayne/go-geom"
	"testing"
)

func TestIsInLineString(t *testing.T) {
	interiorPoint := geom.NewPointFlat(geom.XY, []float64{0.5, 0.5})
	exteriorPoint := geom.NewPointFlat(geom.XY, []float64{1.5, 1.5})
	ring := geom.NewLinearRingFlat(geom.XY, []float64{0, 0, 0, 1, 1, 1, 1, 0})
	if !isInLinearRing(ring, interiorPoint) {
		t.Fatalf("Expected interiorPoint %v to be in ring %v", interiorPoint, ring)
	}
	if isInLinearRing(ring, exteriorPoint) {
		t.Fatalf("Expected exteriorPoint %v to be outside ring %v", exteriorPoint, ring)
	}
}

func TestIsInHoledPolygon(t *testing.T) {
	interiorPoint := geom.NewPointFlat(geom.XY, []float64{3, 3})
	holePoint := geom.NewPointFlat(geom.XY, []float64{1.5, 1.5})
	exteriorPoint := geom.NewPointFlat(geom.XY, []float64{5, 5})

	outerRing := geom.NewLinearRingFlat(geom.XY, []float64{0, 0, 0, 4, 4, 4, 4, 0})
	innerRing := geom.NewLinearRingFlat(geom.XY, []float64{1, 1, 1, 2, 2, 2, 2, 1})

	polygon := geom.NewPolygon(geom.XY)
	polygon.Push(outerRing)
	polygon.Push(innerRing)

	if !isInPolygon(polygon, interiorPoint) {
		t.Fatalf("Expected interiorPoint %v to be in polygon %v", interiorPoint, polygon)
	}
	if isInPolygon(polygon, exteriorPoint) {
		t.Fatalf("Expected exteriorPoint %v to be outside polygon %v", exteriorPoint, polygon)
	}
	if !isInLinearRing(innerRing, holePoint) {
		t.Fatalf("Expected holePoint %v to be in ring %v", interiorPoint, innerRing)
	}
	if isInPolygon(polygon, holePoint) {
		t.Fatalf("Expected holePoint %v to be outside polygon %v", holePoint, polygon)
	}
}

func TestThirdOrderPolygon(t *testing.T) {
	nr := 3
	points := make([]*geom.Point, nr+1)
	polygon := geom.NewPolygon(geom.XY)

	for i := 0; i < nr; i++ {
		l := float64(i)
		r := float64(2*nr - i - 1)
		polygon.Push(geom.NewLinearRingFlat(geom.XY, []float64{
			l, l, l, r,
			r, r, r, l,
		}))
		points[i] = geom.NewPointFlat(geom.XY, []float64{l + 0.5, l + 0.5})
	}
	points[nr] = geom.NewPointFlat(geom.XY, []float64{float64(nr) + 0.5, float64(nr) + 0.5})

	for i, p := range points {
		if i%2 == 0 {
			if !isInPolygon(polygon, p) {
				t.Fatalf("Expected point %d %v to be inside", i, p)
			}
		} else {
			if isInPolygon(polygon, p) {
				t.Fatalf("Expected point %d %v to be outside", i, p)
			}
		}
	}
}

func TestMultiPolygon(t *testing.T) {
	multiPolygon := geom.NewMultiPolygonFlat(geom.XY,
		[]float64{
			0, 0, 0, 1, 1, 1, 1, 0,
			2, 2, 2, 3, 3, 3, 3, 2,
		},
		[][]int{{8}, {16}},
	)

	interiorPoint := geom.NewPointFlat(geom.XY, []float64{0.5, 0.5})
	interiorPoint2 := geom.NewPointFlat(geom.XY, []float64{2.5, 2.5})
	exteriorPoint := geom.NewPointFlat(geom.XY, []float64{5, 5})

	if !isInMultiPolygon(multiPolygon, interiorPoint) {
		t.Fatalf("Expected interiorPoint %v to be in multiPolygon %v", interiorPoint, multiPolygon)
	}
	if !isInMultiPolygon(multiPolygon, interiorPoint2) {
		t.Fatalf("Expected interiorPoint2 %v to be in multiPolygon %v", interiorPoint2, multiPolygon)
	}
	if isInMultiPolygon(multiPolygon, exteriorPoint) {
		t.Fatalf("Expected exteriorPoint %v to be outside multiPolygon %v", exteriorPoint, multiPolygon)
	}

}
