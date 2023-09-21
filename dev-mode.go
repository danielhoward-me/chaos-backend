package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func isDevMode() bool {
	return os.Getenv("NODE_ENV") != "production"
}

func getDevLink(c *fiber.Ctx, devPortQueryName string, domain string, path string, queryParams map[string]string) string {
	if isDevMode() {
		queryParams["testingPort"] = getPort()
	}

	queryString := constructQueryString(queryParams)
	devPort := c.QueryInt(devPortQueryName)

	if devPort == 0 {
		return fmt.Sprintf("https://%s%s%s", domain, path, queryString)
	} else {
		return fmt.Sprintf("http://local.danielhoward.me:%d%s%s", devPort, path, queryString)
	}
}

func constructQueryString(values map[string]string) string {
	if len(values) == 0 {
		return ""
	}

	queryParts := []string{}
	for k, v := range values {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("?%s", strings.Join(queryParts, "&"))
}
