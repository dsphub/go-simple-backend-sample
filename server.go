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

type PostStore interface {
	GetAllPosts() []Post
	GetPostByID(id int) (Post, error)
	CreatePost(title, text string)
	UpdatePost(id int, title, text string) error
}

type Post struct {
	id    int
	title string
	text  string
}

func (p *PostServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Path[len("/posts/"):]
	switch r.Method {
	case http.MethodGet:
		if postID == "" {
			p.getAllPosts(w)
		} else {
			id, _ := strconv.Atoi(postID) //FIXIT error
			p.getPostByID(w, id)
		}
	case http.MethodPost:
		p.CreatePost(w, postID, "test") //FIXIT title, text
	case http.MethodPut:
		if postID == "" {
			w.WriteHeader(http.StatusInternalServerError) //FIXIT status
		} else {
			r.ParseForm()
			id, _ := strconv.Atoi(postID)                              //FIXIT error
			p.UpdatePost(w, id, r.Form["title"][0], r.Form["text"][0]) //FIXIT title, text
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

func (p *PostServer) CreatePost(w http.ResponseWriter, title, text string) {
	p.store.CreatePost(title, text)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PostServer) UpdatePost(w http.ResponseWriter, id int, title, text string) {
	err := p.store.UpdatePost(id, title, text)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}
