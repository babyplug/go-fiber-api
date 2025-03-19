package apikey

import (
	"go-fiber-api/internal/core/model"
	"go-fiber-api/internal/core/response"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	h     *handlerImpl
	hOnce sync.Once
)

type Handler interface {
	Create(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	DeleteByID(c *fiber.Ctx) error
}

type handlerImpl struct {
	s Service
}

func ProvideHandler(s Service) Handler {
	hOnce.Do(func() {
		h = &handlerImpl{
			s: s,
		}
	})

	return h
}

func ResetHandler() {
	hOnce = sync.Once{}
}

func (c *handlerImpl) Create(ctx *fiber.Ctx) error {
	var dto *model.APIKeyDTO
	if err := ctx.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := c.s.Create(ctx.Context(), dto)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&response.ResponseDTO{
		Message: "success",
		Data:    token,
	})
}

func (c *handlerImpl) FindAll(ctx *fiber.Ctx) error {
	res, err := c.s.FindAll(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&response.ResponseDTO{
		Message: "success",
		Data:    res,
	})
}

func (c *handlerImpl) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.s.FindByID(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return ctx.JSON(&response.ResponseDTO{
		Message: "success",
		Data:    data,
	})
}

func (c *handlerImpl) DeleteByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.s.DeleteByID(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(&response.ResponseDTO{
		Message: "success",
	})
}
