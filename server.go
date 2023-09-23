package main

import (
	"github.com/danielhoward-me/chaos-backend/saves"
	"github.com/danielhoward-me/chaos-backend/sso"

	"fmt"
	"os"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var bearerRegex = regexp.MustCompile("^Bearer ")

func createServer() {
	app := fiber.New()

	app.Get("/presets", func(c *fiber.Ctx) error {
		presets, err := saves.Get(db, "")
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(presets)
	})

	app.Use(func(c *fiber.Ctx) error {
		authorisation := bearerRegex.ReplaceAllString(c.GetReqHeaders()["Authorization"], "")
		if authorisation == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		ssoDevPort := c.QueryInt("ssodevport")
		if ssoDevPort != 0 && os.Getenv("NODE_ENV") == "production" {
			return c.Status(fiber.StatusBadRequest).SendString("ssodevport can only be used in development")
		}

		account, exists, err := sso.Get(authorisation, ssoDevPort)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if !exists {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("account", account)
		return c.Next()

	})

	app.Get("/account", func(c *fiber.Ctx) error {
		account := c.Locals("account").(sso.Account)
		userSaves, err := saves.Get(db, account.UserId)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(map[string]interface{}{
			"account": account,
			"saves":   userSaves,
		})
	})

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
