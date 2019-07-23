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
	postTitle := "title"

	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(postTitle))
	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(postTitle))
	server.ServeHTTP(httptest.NewRecorder(), newCreatePostRequest(postTitle))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetAllPostsRequest())
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, "[{1 title test} {2 title test} {3 title test}]", response.Body.String())
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
