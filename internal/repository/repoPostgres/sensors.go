package repoPostgres

import (
	"context"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"

	"github.com/jmoiron/sqlx"
)

type SensorsRepo struct {
	db *sqlx.DB
}

func NewSensorsRepo(db *sqlx.DB) SensorsRepo {
	return SensorsRepo{
		db: db,
	}
}

func (r *SensorsRepo) GetSensorsByRegionID(ctx context.Context, regionId int64) ([]entities.Sensor, error) {
	var sensors []entities.Sensor

	query := fmt.Sprintf(`
	SELECT * FROM %s AS s
	WHERE s.region_id = $1`, sensorsTable)

	if err := r.db.SelectContext(ctx, &sensors, query, regionId); err != nil {
		log.Printf("SensorsRepo GetSensorsByRegionID error: %s", err)
		return nil, err
	}

	if len(sensors) == 0 {
		log.Printf("SensorsRepo GetSensorsByRegionID no recors in table %s", sensorsTable)
		return nil, fmt.Errorf("%w: table %s", repoErrors.ErrNoRecords, sensorsTable)
	}

	return sensors, nil
}

func (r *SensorsRepo) GetPaginatedSensors(ctx context.Context, limit, offset int64) ([]entities.Sensor, error) {
	var (
		sensors []entities.Sensor
		query   string
		err     error
	)

	if limit == -1 {
		query = fmt.Sprintf(`
		SELECT * FROM %s
		ORDER BY id DESC`, sensorsTable)

		err = r.db.SelectContext(ctx, &sensors, query)
	} else {
		query = fmt.Sprintf(`
	SELECT * FROM %s
	ORDER BY id DESC
	LIMIT $1 OFFSET $2`, sensorsTable)

		err = r.db.SelectContext(ctx, &sensors, query, limit, offset)
	}

	if err != nil {
		log.Printf("SensorsRepo GetPaginatedSensors error: %s", err)
		return nil, err
	}

	if len(sensors) == 0 {
		log.Printf("SensorsRepo GetPaginatedSensors no recors in table %s for limit: %d, offset: %d", sensorsTable, limit, offset)
		return nil, fmt.Errorf("%w: table %s, limit: %d, offset: %d", repoErrors.ErrNoRecords, sensorsTable, limit, offset)
	}

	return sensors, nil
}
