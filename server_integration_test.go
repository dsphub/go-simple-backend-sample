package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/dsphub/go-simple-crud-sample/model"
	. "github.com/dsphub/go-simple-crud-sample/testdata"
)

func NewInMemoryPostStore() *StubPostStore {
	return &StubPostStore{
		0,
		map[int]Post{
			1: Post{1, "title", "text"},
		},
	}
}

func EmptyInMemoryPostStore() *StubPostStore {
	return &StubPostStore{
		0,
		map[int]Post{},
	}
}

func TestCreatingPostsAndRetrievingThem(t *testing.T) {
	store := EmptyInMemoryPostStore()
	server := NewPostServer(store)
	title, text := "title", "text"
	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(title, text))
	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(title, text))
	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(title, text))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, newGetAllPostsRequest())

	got := getPostsFromResponse(t, response.Body)
	want, _ := store.GetAllPosts()
	assertStatus(t, response.Code, http.StatusOK)
	assertPosts(t, got, want)
}

func TestUpdatingThePostAndRetrievingIt(t *testing.T) {
	store := NewInMemoryPostStore()
	server := NewPostServer(store)
	server.ServeHTTP(httptest.NewRecorder(), newUpdatePostRequest(1, "new title", "new text"))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, newGetAllPostsRequest())

	got := getPostsFromResponse(t, response.Body)
	want, _ := store.GetAllPosts()
	assertStatus(t, response.Code, http.StatusOK)
	assertPosts(t, got, want)
}

func TestDeletingThePostAndRetrievingOtherOnes(t *testing.T) {
	store := NewInMemoryPostStore()
	server := NewPostServer(store)
	server.ServeHTTP(httptest.NewRecorder(), newDeletePostRequest(1))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, newGetAllPostsRequest())

	got := getPostsFromResponse(t, response.Body)
	want, _ := store.GetAllPosts()
	assertStatus(t, response.Code, http.StatusOK)
	assertPosts(t, got, want)
}
