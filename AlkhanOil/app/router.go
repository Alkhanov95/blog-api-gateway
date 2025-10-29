package app

import "github.com/gofiber/fiber"

func setupUserRouters(h *handler.Handle) *fiber.App {
	app := fiber.New()

	app.Get("/items", GetAllItems)
	app.Get("/items/:id", GetItemByID)
	app.Post("/items", CreateItem)
	app.Put("/items/:id/increase", IncreaseQuantity)
	app.Put("/items/:id/decrease", DecreaseQuantity)
	app.Delete("/items/:id", DeleteByID)

	return app

}
