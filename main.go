// Serve Chronicle Fiber App
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/wesleyflorence/chronicle/notion"
)

func setupRoutes(app *fiber.App) {
	godotenv.Load()
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	digestionDbID := os.Getenv("DIGESTION_DB")
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))

	app.Get("/", hello)

	app.Post("/api/v1/dig", func(c *fiber.Ctx) error {
		type Payload struct {
			Message string `json:"message"`
		}

		var payload Payload
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		page, err := notion.AppendDigestionEntry(client, digestionDbID, 4, "Large", "Sent from chronicle")

		if err != nil {
			return err
		}

		return c.JSON(page)
	})
}

func main() {
	app := fiber.New()
	setupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
