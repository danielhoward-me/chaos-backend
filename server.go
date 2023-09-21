package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func createServer() {
	app := fiber.New()

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Redirect(getDevLink(c, "ssoDevPort", "sso.danielhoward.me", "/auth", map[string]string{
			"target": "chaos",
		}))
	})

	app.Listen(fmt.Sprintf(":%s", getPort()))
}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return
}
