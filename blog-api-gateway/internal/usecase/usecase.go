package usecase

import "github.com/mtvy/blog-api-gateway/internal/models"

type postProvider interface {
	ListPost() ([]models.PostDTO, error)
	GetPost(id uint64) (*models.PostDTO, error)
	CreatePost(post models.PostDTO) (uint64, error)
	DeletePost(id uint64) error
	UpdatePost(post models.PostDTO) error
}

type Usecase struct {
	postRepo postProvider
}

func NewPostProvider(postRepo postProvider) *Usecase {
	return &Usecase{
		postRepo: postRepo,
	}
}

func (u *Usecase) ListPost() ([]models.PostDTO, error) {
	return u.postRepo.ListPost()
}

func (u *Usecase) GetPost(id uint64) (*models.PostDTO, error) {
	return u.postRepo.GetPost(id)
}

func (u *Usecase) CreatePost(post models.PostDTO) (uint64, error) {
	return u.postRepo.CreatePost(post)
}

func (u *Usecase) UpdatePost(post models.PostDTO) error {
	return u.postRepo.UpdatePost(post)
}

func (u *Usecase) DeletePost(id uint64) error {
	return u.postRepo.DeletePost(id)
}
