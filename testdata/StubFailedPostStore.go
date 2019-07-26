package testdata

import (
	"errors"
	. "github.com/dsphub/go-simple-crud-sample/model"
)

type StubFailedPostStore struct{}

func (s *StubFailedPostStore) GetAllPosts() ([]Post, error) {
	return []Post{}, ErrorPostsAreNotFound
}

func (s *StubFailedPostStore) GetPostByID(id int) (Post, error) {
	return Post{}, ErrorPostDoesNotExist
}

func (s *StubFailedPostStore) CreatePost(title, text string) error {
	return ErrorPostIsNotCreated
}

func (s *StubFailedPostStore) UpdatePost(id int, title, text string) error {
	return ErrorPostDoesNotExist
}

func (s *StubFailedPostStore) DeletePost(id int) error {
	return ErrorPostDoesNotExist
}

func (i *StubFailedPostStore) Close() error {
	return errors.New("failed to close store")
}
