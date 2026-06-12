package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
)

type InmemRepository struct {
	trips map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}


func NewInmemRepository() *InmemRepository {
	return &InmemRepository{
		trips: make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}


func (r *InmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *InmemRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.rideFares[f.ID.Hex()] = f
	return nil
}

func (r *InmemRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error){
	fare, exist := r.rideFares[id]
	if !exist {
		return nil, fmt.Errorf("fare does not exist with id: %s", id)
	}

	return fare, nil
}