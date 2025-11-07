package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mtvy/blog-api-gateway/internal/apperr"
	"github.com/mtvy/blog-api-gateway/internal/models"
	"github.com/mtvy/blog-api-gateway/internal/validator"
	"github.com/pkg/errors"

	_ "github.com/mtvy/blog-api-gateway/docs"
)

type postsProvider interface {
	ListPost() ([]models.PostDTO, error)
	GetPost(id uint64) (*models.PostDTO, error)
	CreatePost(post models.PostDTO) (uint64, error)
	UpdatePost(post models.PostDTO) error
	DeletePost(id uint64) error
}

type Handle struct {
	postsUC postsProvider
}

func New(postsUC postsProvider) *Handle {
	return &Handle{
		postsUC: postsUC,
	}
}

// GetPost получает пост по ID.
//
//	@Summary		Получить пост
//	@Description	Получить пост по идентификатору
//	@Tags			posts
//	@Param			id	path		int	true	"ID поста"
//	@Success		200	{object}	models.PostDTO
//	@Failure		400	{object}	map[string]interface{}	"Некорректный ID"
//	@Failure		404	{object}	map[string]interface{}	"Пост не найден"
//	@Failure		500	{object}	map[string]interface{}	"Внутренняя ошибка сервера"
//	@Router			/posts/{id} [get]
func (h *Handle) GetPost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "id should be uint")
	}

	post, err := h.postsUC.GetPost(id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.NewError(http.StatusNotFound)
		}
		slog.Error("get post", slog.Any("error", err))
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(post)
}

// ListPost возвращает список всех постов.
//
//	@Summary		Список постов
//	@Description	Получить список всех постов
//	@Tags			posts
//	@Success		200	{object}	map[string][]models.PostDTO
//	@Failure		500	{object}	map[string]interface{}	"Внутренняя ошибка сервера"
//	@Router			/posts [get]
func (h *Handle) ListPost(c *fiber.Ctx) error {
	posts, err := h.postsUC.ListPost()
	if err != nil {
		slog.Error("get post", slog.Any("error", err))
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"posts": posts})
}

// CreatePost создает новый пост.
//
//	@Summary		Создать пост
//	@Description	Создать новый пост
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		models.CreatePostRequest	true	"Данные поста"
//	@Success		200		{object}	map[string]uint64			"ID созданного поста"
//	@Failure		400		{object}	map[string]interface{}		"Ошибка валидации"
//	@Failure		500		{object}	map[string]interface{}		"Внутренняя ошибка сервера"
//	@Router			/posts [post]
func (h *Handle) CreatePost(c *fiber.Ctx) error {
	req := &models.CreatePostRequest{}
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := h.postsUC.CreatePost(req.ToDTO())
	if err != nil {
		slog.Error("create post", slog.Any("error", err))
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": id})
}

// UpdatePost обновляет существующий пост.
//
//	@Summary		Обновить пост
//	@Description	Обновить существующий пост
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		models.UpdatePostRequest	true	"Данные для обновления"
//	@Success		200		{object}	map[string]uint64			"ID обновленного поста"
//	@Failure		400		{object}	map[string]interface{}		"Ошибка валидации"
//	@Failure		404		{object}	map[string]interface{}		"Пост не найден"
//	@Failure		500		{object}	map[string]interface{}		"Внутренняя ошибка сервера"
//	@Router			/posts [put]
func (h *Handle) UpdatePost(c *fiber.Ctx) error {
	req := &models.UpdatePostRequest{}
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	post := req.ToDTO()
	if err := h.postsUC.UpdatePost(post); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.NewError(http.StatusNotFound)
		}
		slog.Error("update post", slog.Any("error", err))
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": post.ID})
}

// DeletePost удаляет пост по ID.
//
//	@Summary		Удалить пост
//	@Description	Удалить пост по идентификатору
//	@Tags			posts
//	@Param			id	path	int	true	"ID поста"
//	@Success		200	"Пост удален"
//	@Failure		400	{object}	map[string]interface{}	"Некорректный ID"
//	@Failure		404	{object}	map[string]interface{}	"Пост не найден"
//	@Failure		500	{object}	map[string]interface{}	"Внутренняя ошибка сервера"
//	@Router			/posts/{id} [delete]
func (h *Handle) DeletePost(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "id should be uint")
	}

	if err := h.postsUC.DeletePost(id); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.NewError(http.StatusNotFound)
		}
		slog.Error("delete post", slog.Any("error", err))
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}
