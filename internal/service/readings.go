package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
	"sensors-app/internal/service/serviceErrors"
	"time"
)

type ReadingsRepo interface {
	FindReadingsBySensorID(ctx context.Context, sensorId int64) ([]entities.Reading, error)
	FindReadingsBySensorIdAndTimeInterval(ctx context.Context, sensorId int64, time1, time2 time.Time) ([]entities.Reading, error)
}

type ReadingsService struct {
	readingsRepo ReadingsRepo
}

func NewReadingsService(readingsRepo ReadingsRepo) ReadingsService {
	return ReadingsService{
		readingsRepo: readingsRepo,
	}
}

func (s *ReadingsService) FindReadingsBySensorID(ctx context.Context, sensorId int64, time1, time2 time.Time) ([]entities.Reading, error) {
	var (
		readings    []entities.Reading
		err         error
		currentTime = time.Now()
	)

	if time2.After(currentTime) {
		return nil, fmt.Errorf("%w, endDate:\"%s\" currentDate:\"%s\"", serviceErrors.ErrEndDateAfterCurrDate, time2, currentTime)
	}

	if time1.After(time2) {
		return nil, fmt.Errorf("%w, time1:\"%s\" time2:\"%s\"", serviceErrors.ErrIncorrectDates, time1, time2)
	}

	if time1.IsZero() && time2.IsZero() {
		log.Printf("Start searching for all readings with sensor_id: %d", sensorId)
		readings, err = s.readingsRepo.FindReadingsBySensorID(ctx, sensorId)
	} else {
		log.Printf("Start searching for readings with sensor_id: %d between time1:\"%s\" and time2:\"%s\"", sensorId, time1, time2)
		readings, err = s.readingsRepo.FindReadingsBySensorIdAndTimeInterval(ctx, sensorId, time1, time2)
	}

	if err != nil {
		if errors.Is(err, repoErrors.ErrNoRecords) {
			return nil, fmt.Errorf("%w, sensorID=%d", serviceErrors.ErrNoReadingsData, sensorId)
		}
		return nil, err
	}

	return readings, nil
}
