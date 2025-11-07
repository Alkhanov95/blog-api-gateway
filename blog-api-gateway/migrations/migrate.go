package migrations

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/mtvy/blog-api-gateway/internal/models"
	"github.com/pkg/errors"
)

type postsProvider interface {
	CreatePost(post models.PostDTO) (uint64, error)
}

//go:embed blog_data.json
var blogData []byte

type postsDTO struct {
	Posts []models.PostDTO `json:"posts"`
}

func Migrate(postsRepo postsProvider) error {
	posts := postsDTO{}
	if err := json.Unmarshal(blogData, &posts); err != nil {
		return errors.Wrap(err, "unmarshal blog data")
	}
	for _, post := range posts.Posts {
		if _, err := postsRepo.CreatePost(post); err != nil {
			return errors.Wrap(err, fmt.Sprintf("create post %s", post.Title))
		}
	}

	return nil
}
