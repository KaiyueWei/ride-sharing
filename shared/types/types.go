package types

type Route struct {
	Distance float64     `json:"distance"`
	Duration float64     `json:"duration"`
	Geometry []*Geometry `json:"geometry"`
}

type Geometry struct {
	Coordinates []*Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}


type OsrmApiResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

// ToRoute converts the raw OSRM response (coordinates as [longitude, latitude]
// pairs) into the Route shape the frontend consumes.
func (r *OsrmApiResponse) ToRoute() *Route {
	if len(r.Routes) == 0 {
		return nil
	}

	osrmRoute := r.Routes[0]
	coordinates := make([]*Coordinate, 0, len(osrmRoute.Geometry.Coordinates))
	for _, coord := range osrmRoute.Geometry.Coordinates {
		if len(coord) < 2 {
			continue
		}
		coordinates = append(coordinates, &Coordinate{
			Longitude: coord[0],
			Latitude:  coord[1],
		})
	}

	return &Route{
		Distance: osrmRoute.Distance,
		Duration: osrmRoute.Duration,
		Geometry: []*Geometry{{Coordinates: coordinates}},
	}
}