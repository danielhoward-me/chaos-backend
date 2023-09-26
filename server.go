package main

import (
	"github.com/danielhoward-me/chaos-backend/saves"
	"github.com/danielhoward-me/chaos-backend/screenshot"
	screenshotUtils "github.com/danielhoward-me/chaos-backend/screenshot/utils"
	"github.com/danielhoward-me/chaos-backend/sso"

	"fmt"
	"os"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var bearerRegex = regexp.MustCompile("^Bearer ")

func createServer() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://chaos.danielhoward.me, http://local.danielhoward.me:3001",
	}))

	app.Get("/presets", func(c *fiber.Ctx) error {
		presets, err := saves.Get(db, "")
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(presets)
	})

	app.Get("/screenshot/:hash.jpg", func(c *fiber.Ctx) error {
		hash := c.Params("hash")

		if !screenshotUtils.Exists(hash) {
			c.Set(fiber.HeaderContentType, "image/jpeg")
			return c.Send(screenshot.PlaceholderImage)
		}

		return c.SendFile(screenshotUtils.Path(hash))
	})

	app.Post("/screenshot", func(c *fiber.Ctx) error {
		body := screenshot.Request{}
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		data := body.Data
		waitTime := screenshot.Queue(data)

		return c.JSON(map[string]any{
			"hash":           screenshotUtils.Hash(data),
			"screenshotTime": waitTime,
		})
	})

	app.Get("/account", func(c *fiber.Ctx) error {
		account, err := getAccount(c)
		if err != nil {
			return err
		}

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

	app.Get("/delete", func(c *fiber.Ctx) error {
		account, err := getAccount(c)
		if err != nil {
			return err
		}

		id := c.QueryInt("id")
		if id == 0 {
			return c.Status(fiber.StatusBadRequest).SendString("id should be an integer")
		}

		completed, err := saves.Delete(db, id, account.UserId)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if completed {
			return c.JSON(map[string]bool{"ok": true})
		} else {
			return c.SendStatus(fiber.StatusForbidden)
		}
	})

	app.Post("/create", func(c *fiber.Ctx) error {
		account, err := getAccount(c)
		if err != nil {
			return err
		}

		var body saves.RequestSave
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		save, err := saves.Create(db, body.Name, body.Data, account.UserId)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		waitTime := screenshot.Queue(save.Data)

		return c.JSON(map[string]any{
			"save":           save,
			"screenshotTime": waitTime,
		})
	})

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func getAccount(c *fiber.Ctx) (sso.Account, error) {
	authorisation := bearerRegex.ReplaceAllString(c.GetReqHeaders()["Authorization"], "")
	if authorisation == "" {
		return sso.Account{}, c.SendStatus(fiber.StatusUnauthorized)
	}

	account, exists, err := sso.Get(authorisation)
	if err != nil {
		fmt.Println(err)
		return sso.Account{}, c.SendStatus(fiber.StatusInternalServerError)
	}

	if !exists {
		return sso.Account{}, c.SendStatus(fiber.StatusUnauthorized)
	}

	return account, nil
}
