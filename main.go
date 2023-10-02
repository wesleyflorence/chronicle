// Serve Chronicle Fiber App
package main

import (
	"log"
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/wesleyflorence/chronicle/notion"
)

const admin = "wesley"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Could not load env vars")
	}
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
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
	users := setupUsers()
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte("secret")},
		TokenLookup: "cookie:authToken",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Next() // Pass to next handler if JWT fails, instead of erroring out.
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		user := c.Locals("user") // jwtware sets this if JWT is valid.
		isLoggedIn := user != nil
		return c.Render("index", fiber.Map{
			"Title":      "Hello, World!",
			"IsLoggedIn": isLoggedIn,
		})
	})
	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("authToken") // Clear the authToken cookie
		return c.Redirect("/")     // Optionally redirect to homepage or login page
	})
	app.Post("/api/v1/login", handleLogin(users))
	app.Post("/api/v1/dig", handleDigestionEntry(client, digestionDbID))
	app.Post("/api/v1/med", handleMedicineEntry(client, medicinePageID))
}

func handleLogin(users map[string]string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Payload struct {
			Password string
		}
		var payload Payload
		if err := c.BodyParser(&payload); err != nil {

			return err
		}

		if payload.Password != users[admin] {
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

		c.Cookie(&fiber.Cookie{
			Name:     "authToken",
			Value:    t,
			Expires:  time.Now().Add(72 * time.Hour),
			HTTPOnly: true,  // This means the cookie is not accessible via JavaScript
			Secure:   true,  // Use this if your app is served over HTTPS
			SameSite: "Lax", // CSRF protection. You can also consider "Strict" based on your needs.
		})
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
