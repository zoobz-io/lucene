package lucene

import "testing"

type geoTestDoc struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

func TestGeoDistanceQuery(t *testing.T) {
	b, err := New[geoTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic geo_distance", func(t *testing.T) {
		q := b.GeoDistance("location", 40.73, -73.93)

		if q.Op() != OpGeoDistance {
			t.Errorf("Op() = %v, want OpGeoDistance", q.Op())
		}
		if q.Field() != "location" {
			t.Errorf("Field() = %v, want location", q.Field())
		}
		if q.Lat() != 40.73 {
			t.Errorf("Lat() = %v, want 40.73", q.Lat())
		}
		if q.Lon() != -73.93 {
			t.Errorf("Lon() = %v, want -73.93", q.Lon())
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.GeoDistance("location", 40.73, -73.93).
			Distance("10km").
			DistanceType("arc").
			Boost(1.5)

		if q.DistanceValue() == nil || *q.DistanceValue() != "10km" {
			t.Errorf("DistanceValue() = %v, want 10km", q.DistanceValue())
		}
		if q.DistanceTypeValue() == nil || *q.DistanceTypeValue() != "arc" {
			t.Errorf("DistanceTypeValue() = %v, want arc", q.DistanceTypeValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.GeoDistance("invalid_field", 40.73, -73.93)

		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestGeoBoundingBoxQuery(t *testing.T) {
	b, err := New[geoTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic geo_bounding_box", func(t *testing.T) {
		q := b.GeoBoundingBox("location")

		if q.Op() != OpGeoBoundingBox {
			t.Errorf("Op() = %v, want OpGeoBoundingBox", q.Op())
		}
		if q.Field() != "location" {
			t.Errorf("Field() = %v, want location", q.Field())
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("with corners", func(t *testing.T) {
		q := b.GeoBoundingBox("location").
			TopLeft(40.73, -74.1).
			BottomRight(40.01, -71.12)

		if q.TopLeftLat() == nil || *q.TopLeftLat() != 40.73 {
			t.Errorf("TopLeftLat() = %v, want 40.73", q.TopLeftLat())
		}
		if q.TopLeftLon() == nil || *q.TopLeftLon() != -74.1 {
			t.Errorf("TopLeftLon() = %v, want -74.1", q.TopLeftLon())
		}
		if q.BottomRightLat() == nil || *q.BottomRightLat() != 40.01 {
			t.Errorf("BottomRightLat() = %v, want 40.01", q.BottomRightLat())
		}
		if q.BottomRightLon() == nil || *q.BottomRightLon() != -71.12 {
			t.Errorf("BottomRightLon() = %v, want -71.12", q.BottomRightLon())
		}
	})

	t.Run("with boost", func(t *testing.T) {
		q := b.GeoBoundingBox("location").Boost(1.5)

		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.GeoBoundingBox("invalid_field")

		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}
