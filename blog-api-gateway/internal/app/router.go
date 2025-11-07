package app

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/mtvy/blog-api-gateway/internal/handler"
)

func getRouter(handle *handler.Handle) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	posts := app.Group("/posts")
	{
		posts.Get("", handle.ListPost)
		posts.Get("/:id", handle.GetPost)
		posts.Post("", handle.CreatePost)
		posts.Put("", handle.UpdatePost)
		posts.Delete("/:id", handle.DeletePost)
	}
	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	slog.Debug(
		fmt.Sprintf("resp uri=%s body=%s code=%d",
			c.Request().URI().RequestURI(),
			c.Response().Body(),
			c.Response().StatusCode()),
		slog.Any("error", err),
	)

	// check fiber error
	if e, ok := err.(*fiber.Error); ok {
		switch e.Code {
		case fiber.StatusBadRequest:
			return c.Status(400).JSON(fiber.Map{
				"code":        "BadRequest",
				"description": e.Message,
			})
		case fiber.StatusNotFound:
			return c.Status(404).JSON(fiber.Map{
				"code":        "NotFound",
				"description": e.Message,
			})
		case fiber.StatusUnauthorized:
			return c.Status(401).JSON(fiber.Map{
				"code":        "Unauthorized",
				"description": "Недействительный токен аутентификации",
			})
		case fiber.StatusMethodNotAllowed:
			return c.Status(405).JSON(fiber.Map{
				"code":        "Method Not Allowed",
				"description": "Метод не поддерживается",
			})
		case fiber.StatusConflict:
			return c.Status(409).JSON(fiber.Map{
				"code":        "Conflict",
				"description": e.Message,
			})
		}
	}

	return c.Status(500).JSON(fiber.Map{
		"code":        "InternalServerError",
		"description": err.Error(),
	})
}
