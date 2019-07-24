package store

import . "github.com/dsphub/go-simple-crud-sample/model"

type InMemoryPostStore struct {
	Counter int
	Store   map[int]Post
}

func NewInMemoryPostStore() *InMemoryPostStore {
	return &InMemoryPostStore{
		1,
		map[int]Post{
			1: Post{1, "title", "text"},
		},
	}
}

func (i *InMemoryPostStore) GetAllPosts() ([]Post, error) {
	posts := make([]Post, 0, len(i.Store))
	for _, post := range i.Store {
		posts = append(posts, post)
	}
	return posts, nil
}

func (i *InMemoryPostStore) GetPostByID(id int) (Post, error) {
	return i.Store[id], nil
}

func (i *InMemoryPostStore) CreatePost(title, text string) error {
	i.Counter++
	i.Store[i.Counter] = Post{i.Counter, title, text}
	return nil
}

func (i *InMemoryPostStore) UpdatePost(id int, title, text string) error {
	i.Store[id] = Post{id, title, text}
	return nil
}

func (i *InMemoryPostStore) DeletePost(id int) error {
	delete(i.Store, id)
	return nil
}
