package repository

import (
	"sync"

	"github.com/mtvy/blog-api-gateway/internal/apperr"
	"github.com/mtvy/blog-api-gateway/internal/models"
)

type PostRepo struct {
	mu     sync.RWMutex
	posts  map[uint64]models.PostDTO
	lastID uint64
}

func NewPostProvider() *PostRepo {
	return &PostRepo{
		posts: make(map[uint64]models.PostDTO),
	}
}

func (b *PostRepo) get(id uint64) (*models.PostDTO, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	post, ok := b.posts[id]
	return &post, ok
}

func (b *PostRepo) set(post models.PostDTO) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.posts[post.ID] = post
}

func (b *PostRepo) incID() uint64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lastID += 1
	return b.lastID
}

func (b *PostRepo) delete(id uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.posts, id)
}

func (b *PostRepo) ListPost() ([]models.PostDTO, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	posts := make([]models.PostDTO, 0, len(b.posts))
	for _, post := range b.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (b *PostRepo) GetPost(id uint64) (*models.PostDTO, error) {
	if post, ok := b.get(id); ok {
		return post, nil
	}
	return nil, apperr.ErrNotFound
}

func (b *PostRepo) CreatePost(post models.PostDTO) (uint64, error) {
	post.ID = b.incID()
	b.set(post)
	return post.ID, nil
}

func (b *PostRepo) UpdatePost(post models.PostDTO) error {
	if _, ok := b.get(post.ID); !ok {
		return apperr.ErrNotFound
	}
	b.set(post)
	return nil
}

func (b *PostRepo) DeletePost(id uint64) error {
	if _, ok := b.get(id); !ok {
		return apperr.ErrNotFound
	}
	b.delete(id)
	return nil
}
