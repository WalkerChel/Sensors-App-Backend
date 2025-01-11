package repoPostgres

import (
	"context"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
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
