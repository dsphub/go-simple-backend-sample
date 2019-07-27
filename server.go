package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	. "github.com/dsphub/go-simple-crud-sample/model"
	. "github.com/dsphub/go-simple-crud-sample/store"
)

const jsonContentType = "application/json"

type PostServer struct {
	store PostStore
	http.Handler
	log *log.Logger
}

func NewPostServer(store PostStore, log *log.Logger) *PostServer {
	p := new(PostServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/posts/", http.HandlerFunc(p.postsHandler))

	p.Handler = router
	return p
}

func (p *PostServer) postsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Path[len("/posts/"):]
	switch r.Method {
	case http.MethodGet:
		if postID == "" {
			p.getAllPosts(w)
		} else {
			id, err := strconv.Atoi(postID)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			p.getPostByID(w, id)
		}
	case http.MethodPost:
		if postID == "new" {
			r.ParseForm()
			p.CreatePost(w, r.Form["title"][0], r.Form["text"][0]) //FIXIT title, text
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
	case http.MethodPut:
		if postID == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			r.ParseForm()
			id, err := strconv.Atoi(postID)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			p.UpdatePost(w, id, r.Form["title"][0], r.Form["text"][0]) //FIXIT title, text
		}
	case http.MethodDelete:
		if postID == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			id, err := strconv.Atoi(postID)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			p.DeletePost(w, id)
		}
	}
}

func (p *PostServer) getAllPosts(w http.ResponseWriter) {
	posts, err := p.store.GetAllPosts()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	setResponseContentTypeAsJSON(w)
	json.NewEncoder(w).Encode(posts)
}

func setResponseContentTypeAsJSON(w http.ResponseWriter) {
	w.Header().Set("content-type", jsonContentType)
}

func (p *PostServer) getPostByID(w http.ResponseWriter, id int) {
	post, err := p.store.GetPostByID(id)
	switch err {
	case ErrorPostDoesNotExist:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		setResponseContentTypeAsJSON(w)
		json.NewEncoder(w).Encode(post)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (p *PostServer) CreatePost(w http.ResponseWriter, title, text string) {
	err := p.store.CreatePost(title, text)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) //FIXIT status
	}
	w.WriteHeader(http.StatusCreated)
}

func (p *PostServer) UpdatePost(w http.ResponseWriter, id int, title, text string) {
	err := p.store.UpdatePost(id, title, text)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (p *PostServer) DeletePost(w http.ResponseWriter, id int) {
	err := p.store.DeletePost(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)
}
