package main

import . "github.com/dsphub/go-simple-crud-sample/model"

type Service interface {
	GetAllPosts() []Post
	GetPostByID(id int) (Post, error)
	CreatePost(title, text string)
	UpdatePost(id int, title, text string) error
	DeletePost(id int) error
}

/*type postService struct {
	store PostStore
}

/*func NewService(store PostStore) Service {
	return &postService{store}
}*/
