package user

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"reflect"
	"strings"
)

const insertUserQuery = `INSERT INTO users(first_name, last_name, nickname, password, email, country) 
							values(:first_name, :last_name, :nickname, :password, :email, :country)
							RETURNING id, created_at, updated_at;`
const updateUserQuery = `UPDATE users SET first_name=:first_name, last_name=:last_name, 
						nickname=:nickname, password=:password, email=:email, country=:country 
						where id=:id
						RETURNING updated_at;`
const deleteUserByIdQuery = `DELETE FROM users WHERE id=:id;`
const selectUsersQuery = `SELECT id, first_name, last_name, nickname, 
							password, email, country, created_at, updated_at FROM users WHERE 1 = 1`

type repository struct {
	db *sqlx.DB
}

var _ Repository = (*repository)(nil)

type RepositoryOpts func(*repository)

func NewRepository(opts ...RepositoryOpts) *repository {
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

	err = stmt.QueryRowContext(ctx, entity).Scan(&entity.UpdatedAt)
	return entity, err
}

func (r *repository) DeleteById(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareNamedContext(ctx, deleteUserByIdQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, map[string]interface{}{"id": id})
	defer stmt.Close()
	return err
}

func (r *repository) GetMany(ctx context.Context, parameters GetManyParameters) ([]Entity, error) {
	var entities []Entity

	query := r.createGetManyQuery(parameters, selectUsersQuery)
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return entities, err
	}

	rows, err := stmt.QueryxContext(ctx, parameters.Filter)
	if err != nil {
		return entities, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity Entity
		err = rows.StructScan(&entity)
		if err != nil {
			return entities, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *repository) createGetManyQuery(parameters GetManyParameters, query string) string {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(query)

	filterQuery := r.createFilterQuery(parameters.Filter)
	paginationQuery := r.createPaginationQuery(parameters.Page, parameters.PerPage)

	queryBuilder.WriteString(filterQuery)
	queryBuilder.WriteString(" ORDER BY created_at ASC")
	queryBuilder.WriteString(paginationQuery)

	return queryBuilder.String()
}

func (r *repository) createFilterQuery(filter Entity) string {
	filterQuery := strings.Builder{}

	t := reflect.TypeOf(&filter).Elem()
	v := reflect.ValueOf(filter)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.IsZero() {
			typeField := t.Field(i)
			fieldColumnName := typeField.Tag.Get("db")
			filterQuery.WriteString(fmt.Sprintf(" AND %[1]s=:%[1]s", fieldColumnName))
		}
	}

	return filterQuery.String()
}

func (r *repository) createPaginationQuery(page int, perPage int) string {
	offset := perPage * (page - 1)
	limit := perPage
	return fmt.Sprintf("LIMIT %d OFFSET %d;", limit, offset)
}
