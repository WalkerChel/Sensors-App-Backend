package repoPostgres

import (
	"context"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
	"sensors-app/internal/repository/repoPostgres/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReadingsRepo struct {
	db *sqlx.DB
}

func NewReadingsRepo(db *sqlx.DB) ReadingsRepo {
	return ReadingsRepo{
		db: db,
	}
}

func (r *ReadingsRepo) FindReadingsBySensorID(ctx context.Context, sensorId int64) ([]entities.Reading, error) {
	var readings []entities.Reading

	query := fmt.Sprintf(`
	SELECT * FROM %s AS r
	WHERE r.sensor_id = $1
	ORDER BY r.created_at`, readingsTable)

	if err := r.db.SelectContext(ctx, &readings, query, sensorId); err != nil {
		log.Printf("ReadingsRepo FindReadingsBySensorID error: %s", err)
		return nil, err
	}

	if len(readings) == 0 {
		log.Printf("ReadingsRepo FindReadingsBySensorID no recors in table %s for sensor_id=%d", regionsTable, sensorId)
		return nil, fmt.Errorf("%w: table %s, sensor_id=%d", repoErrors.ErrNoRecords, readingsTable, sensorId)
	}

	return readings, nil
}

func (r *ReadingsRepo) FindReadingsBySensorIdAndTimeInterval(ctx context.Context, sensorId int64, time1, time2 time.Time) ([]entities.Reading, error) {
	var readings []entities.Reading

	query := fmt.Sprintf(`
	SELECT * FROM %s AS r
	WHERE r.sensor_id = $1
	AND r.created_at BETWEEN $2 AND $3
	ORDER BY r.created_at`, readingsTable)

	if err := r.db.SelectContext(ctx, &readings, query, sensorId, time1, time2); err != nil {
		log.Printf("ReadingsRepo FindReadingsBySensorIdAndTimeInterval error: %s", err)
		return nil, err
	}

	if len(readings) == 0 {
		log.Printf("ReadingsRepo FindReadingsBySensorIdAndTimeInterval no recors in table %s for sensor_id=%d, time1: %s, time2: %s", regionsTable, sensorId, time1, time2)
		return nil, fmt.Errorf("%w: table %s, sensor_id=%d, time1: %s, time2: %s", repoErrors.ErrNoRecords, readingsTable, sensorId, time1, time2)
	}

	return readings, nil
}

func (r *ReadingsRepo) AltFindReadingsBySensorID(ctx context.Context, sensorId int64) ([]entities.Reading, float64, float64, float64, error) {
	var readings []entities.Reading
	var altReadings []models.ReadingsWthStats

	altQuery := fmt.Sprintf(`
	WITH GetTemperatureStats AS (
		Select 
  			min(temperature) AS temperature_min,
  			avg(temperature) AS temperature_avg,
  			max(temperature) AS temperature_max
 		FROM %s 
 		WHERE sensor_id = $1
		)
	SELECT r.*, g.* FROM %s AS r
	CROSS JOIN GetTemperatureStats AS g
	WHERE r.sensor_id = $1
	ORDER BY r.created_at `, readingsTable, readingsTable)

	if err := r.db.SelectContext(ctx, &altReadings, altQuery, sensorId); err != nil {
		log.Printf("ReadingsRepo AltFindReadingsBySensorID error: %s", err)
		return nil, 0, 0, 0, err
	}

	if len(altReadings) == 0 {
		log.Printf("ReadingsRepo AltFindReadingsBySensorID no recors in table %s for sensor_id=%d", regionsTable, sensorId)
		return nil, 0, 0, 0, fmt.Errorf("%w: table %s, sensor_id=%d", repoErrors.ErrNoRecords, readingsTable, sensorId)
	}

	for _, r := range altReadings {
		readings = append(readings, entities.Reading{
			Id:          r.Id,
			SensorId:    r.SensorId,
			Temperature: r.Temperature,
			CreatedAt:   r.CreatedAt,
		})
	}

	minTemperature := altReadings[0].MinTemperature
	avgTemperature := altReadings[0].AvgTemperature
	maxTemperature := altReadings[0].MaxTemperature

	return readings, minTemperature, avgTemperature, maxTemperature, nil
}

func (r *ReadingsRepo) AltFindReadingsBySensorIdAndTimeInterval(ctx context.Context, sensorId int64, time1, time2 time.Time) ([]entities.Reading, float64, float64, float64, error) {
	var readings []entities.Reading
	var altReadings []models.ReadingsWthStats

	altQuery := fmt.Sprintf(`
	WITH GetTemperatureStats AS (
		Select 
  			min(temperature) AS temperature_min,
  			avg(temperature) AS temperature_avg,
  			max(temperature) AS temperature_max
 		FROM %s 
 		WHERE sensor_id = $1
		AND created_at BETWEEN $2 AND $3
		)
	SELECT r.*, g.* FROM %s AS r
	CROSS JOIN GetTemperatureStats AS g
	WHERE r.sensor_id = $1
	AND r.created_at BETWEEN $2 AND $3
	ORDER BY r.created_at `, readingsTable, readingsTable)

	if err := r.db.SelectContext(ctx, &altReadings, altQuery, sensorId, time1, time2); err != nil {
		log.Printf("ReadingsRepo AltFindReadingsBySensorID error: %s", err)
		return nil, 0, 0, 0, err
	}

	if len(altReadings) == 0 {
		log.Printf("ReadingsRepo AltFindReadingsBySensorID no recors in table %s for sensor_id=%d, start-date: [%s], end-date: [%s]", regionsTable, sensorId, time1, time2)
		return nil, 0, 0, 0, fmt.Errorf("%w: table %s, sensor_id=%d, start-time: [%s], end-time: [%s]", repoErrors.ErrNoRecords, readingsTable, sensorId, time1, time2)
	}

	for _, r := range altReadings {
		readings = append(readings, entities.Reading{
			Id:          r.Id,
			SensorId:    r.SensorId,
			Temperature: r.Temperature,
			CreatedAt:   r.CreatedAt,
		})
	}

	minTemperature := altReadings[0].MinTemperature
	avgTemperature := altReadings[0].AvgTemperature
	maxTemperature := altReadings[0].MaxTemperature

	return readings, minTemperature, avgTemperature, maxTemperature, nil
}
