package user

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const insertUserQuery = `INSERT INTO users(first_name, last_name, nickname, password, email, country) 
							values(:first_name, :last_name, :nickname, :password, :email, :country)
							RETURNING id, created_at, updated_at;`
const updateUserQuery = `UPDATE users SET first_name=:first_name, last_name=:last_name, 
						nickname=:nickname, password=:password, email=:email, country=:country 
						where id=:id
						RETURNING created_at, updated_at;`
const deleteUserByIdQuery = `DELETE FROM users WHERE id=:id;`
const selectUsersQuery = `SELECT * FROM users;`

type repository struct {
	db *sqlx.DB
}

var _ Repository = (*repository)(nil)

type RepositoryOpts func(*repository)

func NewRepository(opts ...RepositoryOpts) Repository {
	r := &repository{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithDb(db *sqlx.DB) RepositoryOpts {
	return func(r *repository) {
		r.db = db
	}
}

func (r *repository) Create(ctx context.Context, entity Entity) (Entity, error) {
	stmt, err := r.db.PrepareNamedContext(ctx, insertUserQuery)
	if err != nil {
		return Entity{}, err
	}

	err = stmt.QueryRowxContext(ctx, entity).Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		return Entity{}, err
	}

	return entity, nil
}

func (r *repository) Update(ctx context.Context, entity Entity) (Entity, error) {
	stmt, err := r.db.PrepareNamedContext(ctx, updateUserQuery)
	if err != nil {
		return Entity{}, err
	}

	err = stmt.QueryRowContext(ctx, entity).Scan(&entity.CreatedAt, &entity.UpdatedAt)
	return entity, err
}

func (r *repository) DeleteById(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareNamedContext(ctx, deleteUserByIdQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, map[string]interface{}{"id": id})
	return err
}

func (r *repository) GetMany(ctx context.Context) ([]Entity, error) {
	var entities []Entity
	err := r.db.Select(&entities, selectUsersQuery)
	if err != nil {
		return nil, err
	}

	return entities, nil
}
