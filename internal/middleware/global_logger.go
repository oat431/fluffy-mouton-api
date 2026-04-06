package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func GlobalLogger(c fiber.Ctx) error {
	log.Info("", c.Method(), " ", c.Path())
	err := c.Next()
	if err != nil {
		log.Error("Error occurred while processing request: ", err)
	}
	return err
}
