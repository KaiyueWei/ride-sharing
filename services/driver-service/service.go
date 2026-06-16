package main

import (
	"math/rand/v2"
	pb "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/util"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type Service struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

type driverInMap struct {
	Driver *pb.Driver
	// index int
	// todo: route
}

func NewService() *Service {
	return &Service{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *Service) RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := rand.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]
	randomPlate := GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)
	// we can ignore this property for now, but it must be sent to the frontend.
	geoHash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		Geohash:        geoHash,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:           "Lando Norris",
		Id:             driverId,
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlate,
	}

	// TODO: Add driver to list
	s.drivers = append(s.drivers, &driverInMap{
		Driver: driver,
	})

	return driver, nil
}

func (s *Service) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: Filter driver from list

	for i, driver := range s.drivers {
		if driver.Driver.Id == driverId{
			s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
		}
	}
}


