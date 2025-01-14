package service

import (
	"context"
	"errors"
	"fmt"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
	"sensors-app/internal/service/serviceErrors"
)

type SensorRepo interface {
	GetSensorsByRegionID(ctx context.Context, regionId int64) ([]entities.Sensor, error)
	GetPaginatedSensors(ctx context.Context, limit, offset int64) ([]entities.Sensor, int64, error)
}

type SensorService struct {
	sensorRepo SensorRepo
}

func NewSensorService(sensorRepo SensorRepo) SensorService {
	return SensorService{
		sensorRepo: sensorRepo,
	}
}

func (s *SensorService) GetSensorsByRegionID(ctx context.Context, regionId int64) ([]entities.Sensor, error) {
	sensors, err := s.sensorRepo.GetSensorsByRegionID(ctx, regionId)
	if err != nil {
		if errors.Is(err, repoErrors.ErrNoRecords) {
			return nil, fmt.Errorf("%w, regionID: %d", serviceErrors.ErrNoSensorsData, regionId)
		}
		return nil, err
	}

	return sensors, nil
}

func(s *SensorService) GetPaginatedSensors(ctx context.Context, limit, offset int64) ([]entities.Sensor, int64, error) {
	sensors, amount, err := s.sensorRepo.GetPaginatedSensors(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, repoErrors.ErrNoRecords) {
			return nil, 0, fmt.Errorf("%w, limit: %d, offset: %d", serviceErrors.ErrNoSensorsData, limit, offset)
		}
		return nil, 0, err
	}

	return sensors, amount, nil
}
