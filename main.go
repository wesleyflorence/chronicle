// Serve Chronicle Fiber App
package main

import (
	"fmt"
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
		log.Println("No .env file found")
	}
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	setupRoutes(app)
	log.Fatal(app.Listen(":8080"))
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

	app.Static("/public", "./public")
	app.Get("/serviceworker.js", func(c *fiber.Ctx) error {
		return c.SendFile("./public/serviceworker.js")
	})
	app.Get("/manifest.json", func(c *fiber.Ctx) error {
		return c.SendFile("./public/manifest.json")
	})
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte("secret")},
		TokenLookup: "cookie:authToken",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Next() // Pass to next handler if JWT fails, instead of erroring out.
		},
	}))

	app.Get("/", handleHome)
	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("authToken") // Clear the authToken cookie
		c.Set("HX-Redirect", "/")
		return c.SendString("Logged out")
	})
	app.Post("/api/v1/login", handleLogin(users))
	app.Post("/api/v1/dig", handleDigestionEntry(client, digestionDbID))
	app.Post("/api/v1/med", handleMedicineEntry(client, medicinePageID))
}
func handleHome(c *fiber.Ctx) error {
	user := c.Locals("user") // jwtware sets this if JWT is valid.
	isLoggedIn := user != nil
	var title string
	if jwtToken, ok := user.(*jwt.Token); ok {
		claimsName, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || !jwtToken.Valid {
			return fmt.Errorf("Unable to parse JWT Token")
		}
		name, ok := claimsName["name"]
		if ok {
			title = "Chronicle " + name.(string)
		}
	} else {
		title = "Chronicle - Login"
	}
	return c.Render("index", fiber.Map{
		"Title":      title,
		"IsLoggedIn": isLoggedIn,
	})
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

			c.Set("Content-Type", "text/html")
			return c.SendString(`<div class="text-red-500">Login failed. Please try again.</div>`)
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

		c.Set("HX-Redirect", "/")
		return c.SendString("Logged in!")
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

		_, err := notion.AppendMedicineEntry(client, medicinePageID, payload.Medicine, payload.Note)
		if err != nil {
			return err
		}

		return c.SendString(`<div id="medSuccessMessage" style="color:text-stone-300;">Data submitted successfully!</div>`)
		//return c.JSON(page)
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

		_, err := notion.AppendDigestionEntry(client, digestionDbID, payload.Bristol, payload.Size, payload.Note)
		if err != nil {
			return err
		}

		c.Set("Content-Type", "text/html")
		return c.SendString(`<div style="color:text-stone-300;">Data submitted successfully!</div>`)
		//return c.JSON(page)
	}
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Chronicle Home")
}
