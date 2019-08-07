package store

import (
	"database/sql"
	"github.com/pkg/errors"
	"log"

	. "github.com/dsphub/go-simple-crud-sample/model"
)

type PostStore interface {
	Connect() error
	Disconnect() error
	GetAllPosts() ([]Post, error)
	GetPostByID(id int) (Post, error)
	CreatePost(title, text string) error
	UpdatePost(id int, title, text string) error
	DeletePost(id int) error
}

type PostgresPostStore struct {
	db *sql.DB
}

func NewPostgresPostStore(connInfo string) (*PostgresPostStore, error) {
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}
	return &PostgresPostStore{db}, nil
}

func (p *PostgresPostStore) Connect() error {
	return p.db.Ping()
}

func (p *PostgresPostStore) Disconnect() error {
	return p.db.Close()
}

func (p *PostgresPostStore) GetAllPosts() ([]Post, error) {
	rows, err := p.db.Query("SELECT * FROM posts;")
	if err != nil {
		return nil, errors.Wrap(err, "can't get all posts")
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return posts, nil
}

func (p *PostgresPostStore) GetPostByID(id int) (Post, error) {
	var post Post
	err := p.db.QueryRow("SELECT * FROM posts WHERE id = $1;", id).
		Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		log.Fatal(err)
	}
	return post, err
}

func (p *PostgresPostStore) CreatePost(title, content string) error {
	q := "INSERT INTO posts(title, content) VALUES $1, $2);"
	_, err := p.db.Exec(q, title, content)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (p *PostgresPostStore) UpdatePost(id int, title, content string) error {
	q := "UPDATE posts SET title = $2 content = $3 WHERE id = $1;"
	_, err := p.db.Exec(q, id, title, content)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (p *PostgresPostStore) DeletePost(id int) error {
	q := "DELETE FROM posts WHERE id = $1;"
	_, err := p.db.Exec(q, id)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
