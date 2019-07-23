package main

import (
	"encoding/json"
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
	http.Handler
}

func NewPostServer(store PostStore) *PostServer {
	p := new(PostServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/posts/", http.HandlerFunc(p.postsHandler))

	p.Handler = router
	return p
}

type PostStore interface {
	GetAllPosts() []Post
	GetPostByID(id int) (Post, error)
	CreatePost(title, text string)
	UpdatePost(id int, title, text string) error
	DeletePost(id int) error
}

type Post struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (p *PostServer) postsHandler(w http.ResponseWriter, r *http.Request) {
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
		if postID == "new" {
			r.ParseForm()
			p.CreatePost(w, r.Form["title"][0], r.Form["text"][0]) //FIXIT title, text
		} else {
			w.WriteHeader(http.StatusInternalServerError) //FIXIT status
		}
	case http.MethodPut:
		if postID == "" {
			w.WriteHeader(http.StatusInternalServerError) //FIXIT status
		} else {
			r.ParseForm()
			id, _ := strconv.Atoi(postID)                              //FIXIT error
			p.UpdatePost(w, id, r.Form["title"][0], r.Form["text"][0]) //FIXIT title, text
		}
	case http.MethodDelete:
		if postID == "" {
			w.WriteHeader(http.StatusInternalServerError) //FIXIT status
		} else {
			id, _ := strconv.Atoi(postID) //FIXIT error
			p.DeletePost(w, id)
		}
	}
}

func (p *PostServer) getAllPosts(w http.ResponseWriter) {
	posts := p.store.GetAllPosts()
	setResponseContentTypeAsJSON(w)
	json.NewEncoder(w).Encode(posts)
}

func setResponseContentTypeAsJSON(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
}

func (p *PostServer) getPostByID(w http.ResponseWriter, id int) {
	post, err := p.store.GetPostByID(id)
	switch err {
	case ErrorPostDoesNotExists:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		setResponseContentTypeAsJSON(w)
		json.NewEncoder(w).Encode(post)
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

func (p *PostServer) DeletePost(w http.ResponseWriter, id int) {
	err := p.store.DeletePost(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}
