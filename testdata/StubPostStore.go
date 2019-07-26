package testdata

import . "github.com/dsphub/go-simple-crud-sample/model"

type StubPostStore struct {
	Counter int
	Posts   map[int]Post
}

func (s *StubPostStore) GetAllPosts() ([]Post, error) {
	values := make([]Post, 0, len(s.Posts))
	for _, v := range s.Posts {
		values = append(values, v)
	}
	return values, nil
}

func (s *StubPostStore) GetPostByID(id int) (Post, error) {
	post, ok := s.Posts[id]
	if !ok {
		var p Post
		return p, ErrorPostDoesNotExist
	}
	return post, nil
}

func (s *StubPostStore) CreatePost(title, text string) error {
	if title == "" || text == "" {
		return ErrorPostIsNotCreated
	}
	s.Counter++
	s.Posts[s.Counter] = Post{s.Counter, title, text}
	return nil
}

func (s *StubPostStore) UpdatePost(id int, title, text string) error {
	_, err := s.GetPostByID(id)
	if err != nil {
		return err
	}
	s.Posts[id] = Post{id, title, text}
	return nil
}

func (s *StubPostStore) DeletePost(id int) error {
	_, err := s.GetPostByID(id)
	if err != nil {
		return err
	}
	delete(s.Posts, id)
	return nil
}

func (i *StubPostStore) Close() error {
	return nil
}
