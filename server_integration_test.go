package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func EmptyInMemoryPostStore() *InMemoryPostStore {
	return &InMemoryPostStore{
		0,
		map[int]Post{},
	}
}

func TestCreatingPostsAndRetrievingThem(t *testing.T) {
	store := EmptyInMemoryPostStore()
	server := PostServer{store}
	title, text := "title", "text"

	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(title, text))
	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(title, text))
	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(title, text))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetAllPostsRequest())
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, "[{1 title text} {2 title text} {3 title text}]", response.Body.String())
}

func TestUpdatingThePostAndRetrievingIt(t *testing.T) {
	store := NewInMemoryPostStore()
	server := PostServer{store}

	server.ServeHTTP(httptest.NewRecorder(), newUpdatePostRequest(1, "new title", "new text"))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetAllPostsRequest())
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, "[{1 new title new text}]", response.Body.String())
}
