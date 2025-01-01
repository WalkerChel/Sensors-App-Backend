package repoPostgres

import (
	"context"
	"errors"
	"fmt"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const foreignKeyViolationErrorCode = "23505"

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return UserRepo{
		db: db,
	}
}

func (r *UserRepo) CreateUser(cxt context.Context, user entities.User) (int64, error) {
	var userId int64

	query := fmt.Sprintf(`
	INSERT INTO %s(name, password_hash, email) 
	VALUES($1, $2, $3)
	RETURNING id`,
		userTable)

	row := r.db.QueryRowContext(cxt, query, user.Name, user.Password, user.Email)

	if err := row.Scan(&userId); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == foreignKeyViolationErrorCode {
				return 0, repoErrors.ErrUserAlreadyExists
			}
		}
		return 0, err
	}

	return userId, nil
}

func (r *UserRepo) DeleteUser(cxt context.Context, userId int64) error {
	query := fmt.Sprintf(`
	DELETE FROM %s AS u
	WHERE u.id = $1`,
		userTable)

	res, err := r.db.ExecContext(cxt, query, userId)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("no user was deleted; check userId")
	}

	return err
}
