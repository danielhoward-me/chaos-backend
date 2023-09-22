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

	app.Get("/test", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT user_id, data FROM saves")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		d := []map[string]string{}
		for rows.Next() {
			var userId string
			var data string

			rows.Scan(&userId, &data)
			d = append(d, map[string]string{
				"userId": userId,
				"data":   data,
			})

		}

		return c.JSON(d)
	})

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
