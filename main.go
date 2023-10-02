// Serve Chronicle Fiber App
package main

import (
	"log"
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/wesleyflorence/chronicle/notion"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Could not load env vars")
	}
}

func main() {
	app := fiber.New()
	setupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

func setupUsers() map[string]string {
	username := os.Getenv("CHRONICLE_USERNAME")
	password := os.Getenv("CHRONICLE_PW")
	users := map[string]string{
		username: password,
	}
	return users
}

func setupRoutes(app *fiber.App) {
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	digestionDbID := os.Getenv("DIGESTION_DB")
	medicinePageID := os.Getenv("MEDICINE_PAGE")
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))

	app.Static("/", "./public")
	app.Post("/api/v1/login", handleLogin())
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))
	app.Get("/api/v1/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	app.Post("/api/v1/dig", handleDigestionEntry(client, digestionDbID))
	app.Post("/api/v1/med", handleMedicineEntry(client, medicinePageID))
}

func handleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Payload struct {
			Password string
		}
		var payload Payload
		if err := c.BodyParser(&payload); err != nil {

			return err
		}

		if payload.Password != "test" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"name":  "Wesley",
			"admin": true,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t, "success": true})
	}
}

func handleMedicineEntry(client *notionapi.Client, medicinePageID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Payload struct {
			Medicine string
			Note     string
		}

		var payload Payload
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		page, err := notion.AppendMedicineEntry(client, medicinePageID, payload.Medicine, payload.Note)
		if err != nil {
			return err
		}

		return c.JSON(page)
	}
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
