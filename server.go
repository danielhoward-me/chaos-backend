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
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.SendFile(screenshotUtils.Path(hash))
	})

	app.Post("/screenshot", func(c *fiber.Ctx) error {
		body := screenshot.Request{}
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		data := body.Data
		screenshot.Queue(data)

		return c.JSON(map[string]any{
			"hash": screenshotUtils.Hash(data),
		})
	})

	app.Get("/screenshot/status", func(c *fiber.Ctx) error {
		hash := c.Query("hash")

		if hash == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(map[string]any{
			"status": screenshot.GetStatus(hash),
		})
	})

	app.Get("/account", func(c *fiber.Ctx) error {
		account, failed, err := getAccount(c)
		if failed {
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
		account, failed, err := getAccount(c)
		if failed {
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
		account, failed, err := getAccount(c)
		if failed {
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

		screenshot.Queue(save.Data)

		return c.JSON(map[string]any{
			"save": save,
		})
	})

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func getAccount(c *fiber.Ctx) (sso.Account, bool, error) {
	authorisation := bearerRegex.ReplaceAllString(c.GetReqHeaders()["Authorization"], "")
	if authorisation == "" {
		return sso.Account{}, true, c.SendStatus(fiber.StatusUnauthorized)
	}

	account, exists, err := sso.Get(authorisation)
	if err != nil {
		fmt.Println(err)
		return sso.Account{}, true, c.SendStatus(fiber.StatusInternalServerError)
	}

	if !exists {
		return sso.Account{}, true, c.SendStatus(fiber.StatusUnauthorized)
	}

	return account, false, nil
}
