package model

const (
	ErrorPostsAreNotFound = PostError("could not find posts")
	ErrorPostDoesNotExist = PostError("could not find the post by id")
	ErrorPostIsNotCreated = PostError("could not create the post")
)

type PostError string

func (e PostError) Error() string {
	return string(e)
}
