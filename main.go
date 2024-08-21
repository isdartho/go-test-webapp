package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c fiber.Ctx) error {
		return c.Render("index", fiber.Map{"data": "test"})
	})

	app.Listen(":3000")
}
