package lucene

// GeoDistanceQuery matches documents within a radius from a point.
type GeoDistanceQuery struct {
	query
	lat          float64
	lon          float64
	distance     *string
	distanceType *string
	boost        *float64
}

// Distance sets the radius (e.g., "10km", "5mi").
func (q *GeoDistanceQuery) Distance(d string) *GeoDistanceQuery { q.distance = &d; return q }

// DistanceType sets the distance calculation type ("arc" or "plane").
func (q *GeoDistanceQuery) DistanceType(t string) *GeoDistanceQuery { q.distanceType = &t; return q }

// Boost sets the relevance score multiplier.
func (q *GeoDistanceQuery) Boost(b float64) *GeoDistanceQuery { q.boost = &b; return q }

// Lat returns the latitude.
func (q *GeoDistanceQuery) Lat() float64 { return q.lat }

// Lon returns the longitude.
func (q *GeoDistanceQuery) Lon() float64 { return q.lon }

// DistanceValue returns the distance value if set.
func (q *GeoDistanceQuery) DistanceValue() *string { return q.distance }

// DistanceTypeValue returns the distance_type value if set.
func (q *GeoDistanceQuery) DistanceTypeValue() *string { return q.distanceType }

// BoostValue returns the boost value if set.
func (q *GeoDistanceQuery) BoostValue() *float64 { return q.boost }

// GeoDistance creates a geo_distance query.
func (b *Builder[T]) GeoDistance(field string, lat, lon float64) *GeoDistanceQuery {
	spec, errQ := b.validateField(OpGeoDistance, field)
	if errQ != nil {
		return &GeoDistanceQuery{query: *errQ}
	}
	return &GeoDistanceQuery{
		query: query{op: OpGeoDistance, field: spec.Name},
		lat:   lat,
		lon:   lon,
	}
}

// GeoBoundingBoxQuery matches documents within a bounding box.
type GeoBoundingBoxQuery struct {
	query
	topLeftLat     *float64
	topLeftLon     *float64
	bottomRightLat *float64
	bottomRightLon *float64
	boost          *float64
}

// TopLeft sets the top-left corner of the bounding box.
func (q *GeoBoundingBoxQuery) TopLeft(lat, lon float64) *GeoBoundingBoxQuery {
	q.topLeftLat = &lat
	q.topLeftLon = &lon
	return q
}

// BottomRight sets the bottom-right corner of the bounding box.
func (q *GeoBoundingBoxQuery) BottomRight(lat, lon float64) *GeoBoundingBoxQuery {
	q.bottomRightLat = &lat
	q.bottomRightLon = &lon
	return q
}

// Boost sets the relevance score multiplier.
func (q *GeoBoundingBoxQuery) Boost(b float64) *GeoBoundingBoxQuery { q.boost = &b; return q }

// TopLeftLat returns the top-left latitude if set.
func (q *GeoBoundingBoxQuery) TopLeftLat() *float64 { return q.topLeftLat }

// TopLeftLon returns the top-left longitude if set.
func (q *GeoBoundingBoxQuery) TopLeftLon() *float64 { return q.topLeftLon }

// BottomRightLat returns the bottom-right latitude if set.
func (q *GeoBoundingBoxQuery) BottomRightLat() *float64 { return q.bottomRightLat }

// BottomRightLon returns the bottom-right longitude if set.
func (q *GeoBoundingBoxQuery) BottomRightLon() *float64 { return q.bottomRightLon }

// BoostValue returns the boost value if set.
func (q *GeoBoundingBoxQuery) BoostValue() *float64 { return q.boost }

// GeoBoundingBox creates a geo_bounding_box query.
func (b *Builder[T]) GeoBoundingBox(field string) *GeoBoundingBoxQuery {
	spec, errQ := b.validateField(OpGeoBoundingBox, field)
	if errQ != nil {
		return &GeoBoundingBoxQuery{query: *errQ}
	}
	return &GeoBoundingBoxQuery{
		query: query{op: OpGeoBoundingBox, field: spec.Name},
	}
}
