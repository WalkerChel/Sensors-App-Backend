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
