package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type StubPostStore struct {
	counter int
	posts   map[int]Post
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

func (s *StubPostStore) CreatePost(title, text string) {
	s.counter++
	s.posts[s.counter] = Post{s.counter, title, text}
}

func (s *StubPostStore) UpdatePost(id int, title, text string) error {
	_, err := s.GetPostByID(id)
	if err != nil {
		return err
	}
	s.posts[id] = Post{id, title, text}
	return nil
}

func (s *StubPostStore) DeletePost(id int) error {
	_, err := s.GetPostByID(id)
	if err != nil {
		return err
	}
	delete(s.posts, id)
	return nil
}

func TestGetPosts(t *testing.T) {
	t.Run("return all posts", func(t *testing.T) {
		const fakeID = 1
		const actualPostCount = 1
		store := StubPostStore{
			actualPostCount,
			map[int]Post{
				fakeID: Post{fakeID, "title", "text"},
			},
		}
		server := NewPostServer(&store)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetAllPostsRequest())

		got := getPostsFromResponse(t, response.Body)
		assertStatus(t, http.StatusOK, response.Code)
		assertContentType(t, response)
		assertPosts(t, got, store.GetAllPosts())
	})

	t.Run("return empty posts list", func(t *testing.T) {
		request := newGetAllPostsRequest()
		response := httptest.NewRecorder()
		store := StubPostStore{0, map[int]Post{}}
		server := NewPostServer(&store)

		server.ServeHTTP(response, request)

		got := getPostsFromResponse(t, response.Body)
		assertStatus(t, http.StatusOK, response.Code)
		assertPosts(t, got, store.GetAllPosts())
	})
}

func newGetAllPostsRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/posts/", nil)
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

func TestGetPostByID(t *testing.T) {
	const fakeID = 1
	const actualPostCount = 1
	t.Run("return post by id", func(t *testing.T) {
		request := newGetPostByIDRequest(fakeID)
		response := httptest.NewRecorder()
		store := StubPostStore{
			actualPostCount,
			map[int]Post{
				fakeID: Post{fakeID, "title", "text"},
			},
		}
		server := NewPostServer(&store)

		server.ServeHTTP(response, request)

		got := getSinglePostFromResponse(t, response.Body)
		assertStatus(t, http.StatusOK, response.Code)
		assertContentType(t, response)
		assertPosts(t, []Post{got}, store.GetAllPosts())
	})

	t.Run("return 404 on missing post", func(t *testing.T) {
		const actualPostCount = 1
		request := newGetPostByIDRequest(2)
		response := httptest.NewRecorder()
		store := StubPostStore{
			actualPostCount,
			map[int]Post{
				fakeID: Post{fakeID, "title", "text"},
			},
		}
		server := NewPostServer(&store)

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newGetPostByIDRequest(id int) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/posts/%d", id), nil)
	return request
}

func TestCreatePost(t *testing.T) {
	const actualPostCount = 0
	const expectedPostCount = 1
	store := StubPostStore{
		actualPostCount,
		map[int]Post{},
	}
	server := NewPostServer(&store)

	t.Run("create a new post)", func(t *testing.T) {
		request := newCreatePostRequest("title", "text")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, http.StatusCreated, response.Code)
		assertPostCount(t, expectedPostCount, len(store.posts))
	})
}

func newCreatePostRequest(title, text string) *http.Request {
	data := url.Values{"title": {title}, "text": {text}}
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/posts/new?%s", data.Encode()), nil)
	return request
}

func TestUpdatePost(t *testing.T) {
	const fakeID = 1
	const actualPostCount = 1
	const expectedPostCount = 1
	store := StubPostStore{
		1,
		map[int]Post{
			fakeID: Post{fakeID, "title", "text"},
		},
	}
	server := NewPostServer(&store)

	t.Run("update all the post details", func(t *testing.T) {
		request := newUpdatePostRequest(fakeID, "new title", "new text")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, http.StatusOK, response.Code)
		assertPostCount(t, expectedPostCount, len(store.posts))
	})

	t.Run("return 404 on missing post", func(t *testing.T) {
		request := newUpdatePostRequest(2, "dummy title", "dummy text")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertPostCount(t, expectedPostCount, len(store.posts))
	})
}

func newUpdatePostRequest(id int, title, text string) *http.Request {
	data := url.Values{"title": {"new title"}, "text": {"new text"}}
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/posts/%d?%s", id, data.Encode()), nil)
	return request
}

func assertPostCount(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("wrong post count %d, want %d", got, want)
	}
}

func TestDeletePost(t *testing.T) {
	t.Run("delet the post by id", func(t *testing.T) {
		const fakeID = 1
		const actualPostCount = 1
		const expectedPostCount = 0
		store := StubPostStore{
			actualPostCount,
			map[int]Post{
				fakeID: Post{fakeID, "title", "text"},
			},
		}
		server := NewPostServer(&store)
		request := newDeletePostRequest(fakeID)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, http.StatusNoContent, response.Code)
		assertPostCount(t, expectedPostCount, len(store.posts))
	})

	t.Run("return 404 on missing post", func(t *testing.T) {
		const fakeID = 1
		const unknownID = 2
		const actualPostCount = 1
		const expectedPostCount = 1
		store := StubPostStore{
			actualPostCount,
			map[int]Post{
				fakeID: Post{fakeID, "title", "text"},
			},
		}
		server := NewPostServer(&store)
		request := newDeletePostRequest(unknownID)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertPostCount(t, expectedPostCount, len(store.posts))
	})
}

func newDeletePostRequest(id int) *http.Request {
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/posts/%d", id), nil)
	return request
}

func getPostsFromResponse(t *testing.T, body io.Reader) (posts []Post) {
	err := json.NewDecoder(body).Decode(&posts)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into list of Post, '%v'", body, err)
	}
	return
}

func getSinglePostFromResponse(t *testing.T, body io.Reader) (post Post) {
	err := json.NewDecoder(body).Decode(&post)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into Post, '%v'", body, err)
	}
	return
}

func assertPosts(t *testing.T, want, got []Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	want := jsonContentType
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
