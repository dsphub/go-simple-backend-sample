package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/dsphub/go-simple-crud-sample/model"
	"github.com/stretchr/testify/assert"
)

func NewTestPostgresPostStore(db *sql.DB) *PostgresPostStore {
	return &PostgresPostStore{db}
}

func TestShouldGetAllPosts(t *testing.T) {
	want := []Post{
		Post{ID: 1, Title: "title1", Content: "text1"},
		Post{ID: 2, Title: "title2", Content: "text2"},
	}
	db, mock, err := dbMock(t)
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(1, "title1", "text1").
		AddRow(2, "title2", "text2")
	mock.ExpectQuery("SELECT \\* FROM (.+)").WillReturnRows(rows)

	store := NewTestPostgresPostStore(db)
	got, err := store.GetAllPosts()

	if assert.NoError(t, err, "Error was not expected while getting all posts") {
		assert.ElementsMatch(t, want, got, "Unexpected posts")
	}
	assert.NoError(t, mock.ExpectationsWereMet(), "Failed read all behaviour")
}

func TestShouldGetPostByID(t *testing.T) {
	want := Post{ID: 1, Title: "title1", Content: "text1"}
	db, mock, err := dbMock(t)
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(want.ID, want.Title, want.Content)
	mock.ExpectQuery("SELECT \\* FROM (.+) WHERE").WillReturnRows(rows)

	store := NewTestPostgresPostStore(db)
	got, err := store.GetPostByID(1)

	if assert.NoError(t, err, "Error was not expected while getting post") {
		assert.Equal(t, want, got, "Unexpected post")
	}
	assert.NoError(t, mock.ExpectationsWereMet(), "Failed read behaviour")
}

func TestShouldCreatePost(t *testing.T) {
	want := Post{1, "title", "new text"}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error on stub database connection: %s", err)
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO (.+) VALUES (.+)").
		WithArgs(want.Title, want.Content).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewTestPostgresPostStore(db)

	err = store.CreatePost(want.Title, want.Content)

	assert.NoError(t, err, "Error was not expected while creating post")
	assert.NoError(t, mock.ExpectationsWereMet(), "Failed create behaviour")
}

func TestShouldUpdatePost(t *testing.T) {
	want := Post{1, "new title", "new text"}
	db, mock, err := dbMock(t)

	defer db.Close()
	mock.ExpectExec("UPDATE (.+) SET (.+) WHERE").
		WithArgs(want.ID, want.Title, want.Content).
		WillReturnResult(sqlmock.NewResult(0, 1))

	store := NewTestPostgresPostStore(db)
	err = store.UpdatePost(want.ID, want.Title, want.Content)

	assert.NoError(t, err, "Error was not expected while updating post")
	assert.NoError(t, mock.ExpectationsWereMet(), "Failed update behaviour")
}

func TestShouldDeletPost(t *testing.T) {
	want := Post{1, "", ""}
	db, mock, err := dbMock(t)
	defer db.Close()
	mock.ExpectExec("DELETE FROM (.+) WHERE").
		WithArgs(want.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewTestPostgresPostStore(db)
	err = store.DeletePost(want.ID)

	assert.NoError(t, err, "Error was not expected while deleting post")
	assert.NoError(t, mock.ExpectationsWereMet(), "Failed delete behaviour")
}

func dbMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error on stub database connection: %s", err)
	}
	return db, mock, err
}
