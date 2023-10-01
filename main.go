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

func main() {
	app := fiber.New()
	setupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

func setupRoutes(app *fiber.App) {
	godotenv.Load()
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	digestionDbID := os.Getenv("DIGESTION_DB")
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))

	app.Get("/", hello)
	app.Post("/api/v1/dig", handleDigestionEntry(client, digestionDbID))
}

func handleDigestionEntry(client *notionapi.Client, digestionDbID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Payload struct {
			Bristol int
			Size    string
			Note    string
		}

		var payload Payload
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		page, err := notion.AppendDigestionEntry(client, digestionDbID, payload.Bristol, payload.Size, payload.Note)
		if err != nil {
			return err
		}

		return c.JSON(page)
	}
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Chronicle Home")
}
