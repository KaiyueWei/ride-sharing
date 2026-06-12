package types

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

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
func (r *OsrmApiResponse) ToRoute() *types.Route {
	if len(r.Routes) == 0 {
		return nil
	}

	osrmRoute := r.Routes[0]
	coordinates := make([]*types.Coordinate, 0, len(osrmRoute.Geometry.Coordinates))
	for _, coord := range osrmRoute.Geometry.Coordinates {
		if len(coord) < 2 {
			continue
		}
		coordinates = append(coordinates, &types.Coordinate{
			Longitude: coord[0],
			Latitude:  coord[1],
		})
	}

	return &types.Route{
		Distance: osrmRoute.Distance,
		Duration: osrmRoute.Duration,
		Geometry: []*types.Geometry{{Coordinates: coordinates}},
	}
}

// ToProto converts the raw OSRM response (coordinates as [longitude, latitude]
// pairs) into the gRPC Route message.
func (r *OsrmApiResponse) ToProto() *pb.Route {
	if len(r.Routes) == 0 {
		return nil
	}

	osrmRoute := r.Routes[0]
	coordinates := make([]*pb.Coordinate, 0, len(osrmRoute.Geometry.Coordinates))
	for _, coord := range osrmRoute.Geometry.Coordinates {
		if len(coord) < 2 {
			continue
		}
		coordinates = append(coordinates, &pb.Coordinate{
			Longitude: coord[0],
			Latitude:  coord[1],
		})
	}

	return &pb.Route{
		Distance: osrmRoute.Distance,
		Duration: osrmRoute.Duration,
		Geometry: []*pb.Geometry{{Coordinates: coordinates}},
	}
}


type PricingConfig struct {
	PricePerUnitOfDistance float64
	PricePerMinute float64
}

func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		PricePerUnitOfDistance: 1.5,
		PricePerMinute: 0.25,
	}
}