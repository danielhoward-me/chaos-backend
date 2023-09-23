package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func createServer() {
	app := fiber.New()

	app.Get("/presets", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT data, screenshot FROM saves WHERE user_id IS NULL")
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
