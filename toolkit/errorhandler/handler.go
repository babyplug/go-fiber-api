package errorhandler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func Handler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		// Default to 500
		code := "0000"
		status := fiber.StatusInternalServerError
		msg := utils.StatusMessage(status)

		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
			code = fmt.Sprint(e.Code)
			msg = e.Message
		}

		return ctx.
			Status(status).
			JSON(fiber.Map{
				"message": msg,
				"code":    code,
				"ok":      false,
			})
	}
}
