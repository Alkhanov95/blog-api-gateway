package usecase

import (
	"testing"

	"github.com/mtvy/blog-api-gateway/internal/apperr"
	"github.com/mtvy/blog-api-gateway/internal/models"
	"github.com/mtvy/blog-api-gateway/internal/repository"
	"github.com/mtvy/blog-api-gateway/migrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) *repository.PostRepo {
	repo := repository.NewPostProvider()
	err := migrations.Migrate(repo)
	require.NoError(t, err)
	return repo
}

func TestUsecase_GetPost(t *testing.T) {
	type want struct {
		post *models.PostDTO
		err  error
	}

	testCases := []struct {
		name string
		id   uint64
		want want
	}{
		{
			name: "valid",
			id:   22,
			want: want{
				post: &models.PostDTO{
					ID:      22,
					Author:  "Author 22",
					Title:   "Title 22",
					Content: "Labore quiquia tempora modi. Dolore ut amet modi sed porro. Dolorem velit porro non adipisci. Etincidunt tempora labore dolore dolorem consectetur. Labore labore quaerat magnam ut. Quaerat ut labore ut modi quaerat. Ipsum ut sit sed ut porro non.",
				},
			},
		},
		{
			name: "post_id_not_found",
			id:   222,
			want: want{
				err: apperr.ErrNotFound,
			},
		},
	}

	repo := newTestRepo(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := NewPostProvider(repo)

			post, err := uc.GetPost(tc.id)
			if tc.want.err != nil {
				require.ErrorContains(t, err, tc.want.err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.want.post, post)
		})
	}
}

func TestUsecase_CreatePost(t *testing.T) {
	type want struct {
		post models.PostDTO
		err  error
	}

	testCases := []struct {
		name string
		post models.PostDTO
		want want
	}{
		{
			name: "create post",
			post: models.PostDTO{
				Author:  "testA",
				Title:   "testB",
				Content: "testC",
			},
			want: want{
				post: models.PostDTO{
					Author:  "testA",
					Title:   "testB",
					Content: "testC",
				},
			},
		},
	}

	repo := newTestRepo(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := NewPostProvider(repo)
			id, err := uc.CreatePost(tc.post)
			if tc.want.err != nil {
				require.ErrorContains(t, err, tc.want.err.Error())
			} else {
				require.NoError(t, err)
			}
			post, err := uc.GetPost(id)
			require.NoError(t, err)

			assert.Equal(t, tc.want.post.Content, post.Content)
			assert.Equal(t, tc.want.post.Author, post.Author)
			assert.Equal(t, tc.want.post.Title, post.Title)
		})
	}
}

func TestUsecase_DeletePost(t *testing.T) {
	type want struct {
		err error
	}

	testCases := []struct {
		name string
		id   uint64
		want want
	}{
		{
			name: "delete_ok",
			id:   22,
			want: want{},
		},
		{
			name: "post_id_not_found",
			id:   222,
			want: want{
				err: apperr.ErrNotFound,
			},
		},
	}

	repo := newTestRepo(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := NewPostProvider(repo)

			err := uc.DeletePost(tc.id)
			if tc.want.err != nil {
				require.ErrorContains(t, err, tc.want.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUsecase_UpdatePost(t *testing.T) {
	type want struct {
		post *models.PostDTO
		err  error
	}

	testCases := []struct {
		name string
		post *models.PostDTO
		want want
	}{
		{
			name: "update_ok",
			post: &models.PostDTO{
				ID:      22,
				Author:  "testA",
				Title:   "testB",
				Content: "testC",
			},
			want: want{
				post: &models.PostDTO{
					ID:      22,
					Author:  "testA",
					Title:   "testB",
					Content: "testC",
				},
			},
		},
		{
			name: "post_not_found",
			post: &models.PostDTO{
				ID:      222,
				Author:  "testA",
				Title:   "testB",
				Content: "testC",
			},
			want: want{
				err: apperr.ErrNotFound,
			},
		},
	}

	repo := newTestRepo(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := NewPostProvider(repo)

			err := uc.UpdatePost(*tc.post)
			if tc.want.err != nil {
				require.ErrorContains(t, err, tc.want.err.Error())
			} else {
				require.NoError(t, err)
				post, err := uc.GetPost(tc.post.ID)
				require.NoError(t, err)

				assert.Equal(t, tc.want.post, post)
			}
		})
	}
}
