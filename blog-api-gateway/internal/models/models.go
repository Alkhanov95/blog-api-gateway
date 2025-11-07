package models

type PostDTO struct {
	ID      uint64
	Title   string
	Author  string
	Content string
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Author  string `json:"author" validate:"required,max=255"`
	Content string `json:"content"`
}

func (c CreatePostRequest) ToDTO() PostDTO {
	return PostDTO{
		Title:   c.Title,
		Author:  c.Author,
		Content: c.Content,
	}
}

type UpdatePostRequest struct {
	ID      uint64 `json:"id" validate:"required,gte=0"`
	Title   string `json:"title" validate:"required,max=255"`
	Author  string `json:"author" validate:"required,max=255"`
	Content string `json:"content"`
}

func (c UpdatePostRequest) ToDTO() PostDTO {
	return PostDTO{
		ID:      c.ID,
		Title:   c.Title,
		Author:  c.Author,
		Content: c.Content,
	}
}
