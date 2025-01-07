package service

import (
	"context"
	"errors"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
	"sensors-app/internal/service/serviceErrors"
)

type RegionRepo interface {
	GetAllRegions(ctx context.Context) ([]entities.Region, error)
}

type RegionService struct {
	regionRepo RegionRepo
}

func NewRegionService(regionRepo RegionRepo) RegionService {
	return RegionService{
		regionRepo: regionRepo,
	}
}

func (s *RegionService) GetAllRegions(ctx context.Context) ([]entities.Region, error) {
	regions, err := s.regionRepo.GetAllRegions(ctx)
	if err != nil {
		if errors.Is(err, repoErrors.ErrNoRecords) {
			return nil, serviceErrors.ErrNoRegionsData
		}
		return nil, err
	}

	return regions, nil
}
