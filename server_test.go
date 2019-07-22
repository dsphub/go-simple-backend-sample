package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPostStore struct {
	posts map[int]Post
}

func (s *StubPostStore) GetAllPosts() []Post {
	values := make([]Post, 0, len(s.posts))
	for _, v := range s.posts {
		values = append(values, v)
	}
	return values
}

func (s *StubPostStore) GetPostByID(id int) (Post, error) {
	post, ok := s.posts[id]
	if !ok {
		return Post{}, ErrorPostDoesNotExists //FIXIT post as null
	}
	return post, nil
}

func TestGETPosts(t *testing.T) {
	t.Run("return single post", func(t *testing.T) {
		request := newGetAllPostsRequest()
		response := httptest.NewRecorder()
		store := StubPostStore{
			map[int]Post{
				1: Post{1, "title", "text"},
			},
		}
		server := &PostServer{&store}

		server.ServeHTTP(response, request)

		assertStatus(t, http.StatusOK, response.Code)
		assertResponseBody(t, "[{1 title text}]", response.Body.String())
	})

	t.Run("return empty posts list", func(t *testing.T) {
		request := newGetAllPostsRequest()
		response := httptest.NewRecorder()
		store := StubPostStore{map[int]Post{}}
		server := &PostServer{&store}

		server.ServeHTTP(response, request)

		assertStatus(t, http.StatusOK, response.Code)
		assertResponseBody(t, "[]", response.Body.String())
	})
}

func newGetAllPostsRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/posts/", nil)
	return request
}

func TestGETPostByID(t *testing.T) {
	const fakeID = 1
	t.Run("return post by id", func(t *testing.T) {
		request := newGetPostByIDRequest(fakeID)
		response := httptest.NewRecorder()
		store := StubPostStore{
			map[int]Post{
				fakeID: Post{fakeID, "title", "text"},
			},
		}
		server := &PostServer{&store}

		server.ServeHTTP(response, request)

		assertStatus(t, http.StatusOK, response.Code)
		assertResponseBody(t, "{1 title text}", response.Body.String())
	})

	t.Run("return 404 on missing post", func(t *testing.T) {
		request := newGetPostByIDRequest(2)
		response := httptest.NewRecorder()
		store := StubPostStore{
			map[int]Post{
				fakeID: Post{fakeID, "title", "text"},
			},
		}
		server := &PostServer{&store}

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newGetPostByIDRequest(id int) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/posts/%d", id), nil)
	return request
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
