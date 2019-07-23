package main

import (
	"log"
	"net/http"
)

func NewInMemoryPostStore() *InMemoryPostStore {
	return &InMemoryPostStore{
		1,
		map[int]Post{
			1: Post{1, "title", "text"},
		},
	}
}

type InMemoryPostStore struct {
	counter int
	store   map[int]Post
}

func (i *InMemoryPostStore) GetAllPosts() []Post {
	posts := make([]Post, 0, len(i.store))
	for _, post := range i.store {
		posts = append(posts, post)
	}
	return posts
}

func (i *InMemoryPostStore) GetPostByID(id int) (Post, error) {
	return i.store[id], nil
}

func (i *InMemoryPostStore) CreatePost(title, text string) {
	i.counter++
	i.store[i.counter] = Post{i.counter, title, text}
}

func (i *InMemoryPostStore) UpdatePost(id int, title, text string) error {
	i.store[id] = Post{id, title, text}
	return nil
}

func (i *InMemoryPostStore) DeletePost(id int) error {
	delete(i.store, id)
	return nil
}

func main() {
	server := NewPostServer(NewInMemoryPostStore())

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
