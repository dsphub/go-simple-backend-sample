package main

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	ErrorPostDoesNotExists = PostError("could not find the post by id")
)

type PostError string

func (e PostError) Error() string {
	return string(e)
}

type PostServer struct {
	store PostStore
}

func (p *PostServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Path[len("/posts/"):]
	if r.Method == http.MethodGet {
		if postID == "" {
			p.getAllPosts(w)
		} else {
			id, _ := strconv.Atoi(postID) //FIXIT error
			p.getPostByID(w, id)
		}
	}
}

func (p *PostServer) getAllPosts(w http.ResponseWriter) {
	posts := p.store.GetAllPosts()
	fmt.Fprint(w, posts)
}

func (p *PostServer) getPostByID(w http.ResponseWriter, id int) {
	post, err := p.store.GetPostByID(id)
	switch err {
	case ErrorPostDoesNotExists:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		fmt.Fprint(w, post)
	default:
		w.WriteHeader(http.StatusInternalServerError) //FIXIT status
	}
}

type PostStore interface {
	GetAllPosts() []Post
	GetPostByID(id int) (Post, error)
}

type Post struct {
	id    int
	title string
	text  string
}
