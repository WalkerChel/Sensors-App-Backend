package repoPostgres

import (
	"context"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"

	"github.com/jmoiron/sqlx"
)

type RegionsRepo struct {
	db *sqlx.DB
}

func NewRegionsRepo(db *sqlx.DB) RegionsRepo {
	return RegionsRepo{
		db: db,
	}
}

func (r *RegionsRepo) CreateRegion() {

}

func (r *RegionsRepo) GetRegionById() {

}

func (r *RegionsRepo) GetAllRegions(ctx context.Context) ([]entities.Region, error) {
	var regions []entities.Region

	query := fmt.Sprintf(`
	SELECT * FROM %s`, regionsTable)

	if err := r.db.SelectContext(ctx, &regions, query); err != nil {
		log.Printf("RegionsRepo GetAllRegions error: %s", err)
		return nil, err
	}

	if len(regions) == 0 {
		log.Printf("RegionsRepo GetAllRegions no recors in table %s", regionsTable)
		return nil, fmt.Errorf("%w: table %s", repoErrors.ErrNoRecords, regionsTable)
	}

	return regions, nil
}

func (r *RegionsRepo) DeleteRegion() {

}
