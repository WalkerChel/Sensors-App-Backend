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
