package user

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var entities = []Entity{
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "UK",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "NL",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "NL",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "UK",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "UK",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "UK",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "TR",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "DE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "UK",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Id:        uuid.New().String(),
		FirstName: "Roberto",
		LastName:  "Firmino",
		Nickname:  "bobby.firmino",
		Password:  "liverpool321",
		Email:     "robertofirmino@lfc.co.uk",
		Country:   "TR",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
var e = Entity{
	Id:        uuid.New().String(),
	FirstName: "Roberto",
	LastName:  "Firmino",
	Nickname:  "bobby.firmino",
	Password:  "liverpool321",
	Email:     "robertofirmino@lfc.co.uk",
	Country:   "UK",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func NewMock() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.WithField("error", err).Fatalf("unexpected error while opening a stub database")
	}

	sqlxDb := sqlx.NewDb(db, "sqlmock")
	return sqlxDb, mock
}

func TestRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := NewMock()
		repo := NewRepository(
			WithDb(db),
		)

		rows := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at"}).AddRow(e.Id, e.CreatedAt, e.UpdatedAt)

		insertQuery := "INSERT INTO users"
		prep := mock.ExpectPrepare(insertQuery)
		prep.ExpectQuery().
			WithArgs(e.FirstName, e.LastName, e.Nickname, e.Password, e.Email, e.Country).
			WillReturnRows(rows)

		actual, err := repo.Create(context.Background(), e)
		assert.NoError(t, err)
		assert.EqualValues(t, e, actual, "should return id, created_at and updated_at fields matching with the database on create query")
	})

	t.Run("query missing argument error", func(t *testing.T) {
		db, mock := NewMock()
		repo := NewRepository(
			WithDb(db),
		)

		insertQuery := "INSERT INTO users"
		prep := mock.ExpectPrepare(insertQuery)
		prep.ExpectQuery().
			WithArgs(e.FirstName, e.LastName, e.Nickname, e.Password, e.Email)

		_, err := repo.Create(context.Background(), e)
		assert.Error(t, err, "should return error when there is a missing argument")
	})
}

func TestRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := NewMock()
		repo := NewRepository(
			WithDb(db),
		)

		rows := sqlmock.
			NewRows([]string{"updated_at"}).AddRow(e.UpdatedAt)

		query := "UPDATE users"
		prep := mock.ExpectPrepare(query)
		prep.ExpectQuery().
			WithArgs(e.FirstName, e.LastName, e.Nickname, e.Password, e.Email, e.Country, e.Id).
			WillReturnRows(rows)

		actual, err := repo.Update(context.Background(), e)
		assert.NoError(t, err)
		assert.EqualValues(t, e, actual, "should return updated_at field matching with the database on update query")
	})
}

func TestRepository_DeleteById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := NewMock()
		repo := NewRepository(
			WithDb(db),
		)

		query := "DELETE FROM users"
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().
			WithArgs(e.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteById(context.Background(), e.Id)
		assert.NoError(t, err)
	})
}

func TestRepository_GetMany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := NewMock()
		repo := NewRepository(
			WithDb(db),
		)

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"})
		for _, e := range entities {
			rows.AddRow(e.Id, e.FirstName, e.LastName, e.Nickname, e.Password, e.Email, e.Country, e.CreatedAt, e.UpdatedAt)
		}

		params := GetManyParameters{
			Page:    1,
			PerPage: 3,
			Filter:  Entity{},
		}

		query := "SELECT (.+) FROM users WHERE"
		prep := mock.ExpectPrepare(query)
		prep.ExpectQuery().WillReturnRows(rows)

		actual, err := repo.GetMany(context.Background(), params)
		assert.NoError(t, err)
		assert.EqualValues(t, entities, actual)
	})

	t.Run("argument mismatch", func(t *testing.T) {
		db, mock := NewMock()
		repo := NewRepository(
			WithDb(db),
		)

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
			AddRow(e.Id, e.FirstName, e.LastName, e.Nickname, e.Password, e.Email, e.Country, e.CreatedAt, e.UpdatedAt)

		params := GetManyParameters{
			Page:    1,
			PerPage: 3,
			Filter:  Entity{},
		}

		query := "SELECT (.+) FROM users WHERE"
		prep := mock.ExpectPrepare(query)
		prep.ExpectQuery().WithArgs(e.Country).WillReturnRows(rows)

		_, err := repo.GetMany(context.Background(), params)
		assert.Error(t, err)
	})
}
