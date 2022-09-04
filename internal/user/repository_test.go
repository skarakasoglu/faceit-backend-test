package user

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var e = Entity{
	Id:        "50bc1cf6-4c96-4b31-83e4-0be8c4c93414",
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

func TestCreate(t *testing.T) {
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
	assert.EqualValues(t, e, actual)
}
