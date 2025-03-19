package user

import (
	"sync"
)

var (
	h     *handlerImpl
	hOnce sync.Once
)

type Handler interface {
	// Get(c *fiber.Ctx) error
	// GetByID(c *fiber.Ctx) error
	// Create(c *fiber.Ctx) error
	// Update(c *fiber.Ctx) error
	// Delete(c *fiber.Ctx) error
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

func ResetProvideHandler() {
	hOnce = sync.Once{}
}

// func (c *handlerImpl) Create(ctx *fiber.Ctx) error {
// 	var dto *model.GameDTO
// 	if err := ctx.BodyParser(&dto); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, err.Error())
// 	}

// 	err := c.s.Create(ctx.Context(), dto)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
// 	}

// 	return ctx.JSON(&model.ResponseDTO{
// 		Message: "success",
// 		Data:    dto,
// 	})
// }

// func (c *handlerImpl) Get(ctx *fiber.Ctx) error {
// 	res, err := c.s.Get(ctx.Context())
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
// 	}

// 	return ctx.JSON(&model.ResponseDTO{
// 		Message: "success",
// 		Data:    res,
// 	})
// }

// func (c *handlerImpl) GetByID(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")

// 	data, err := c.s.GetBy(ctx.Context(), id)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusNotFound, err.Error())
// 	}
// 	return ctx.JSON(&model.ResponseDTO{
// 		Message: "success",
// 		Data:    data,
// 	})
// }

// func (c *handlerImpl) Update(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")

// 	var dto *model.GameDTO
// 	if err := ctx.BodyParser(&dto); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, err.Error())
// 	}

// 	dto.ID = id

// 	err := c.s.UpdateItem(ctx.Context(), dto)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
// 	}

// 	return ctx.JSON(&model.ResponseDTO{
// 		Message: "success",
// 		Data:    dto,
// 	})
// }

// func (c *handlerImpl) Delete(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	err := c.s.Delete(ctx.Context(), id)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
// 	}

// 	return ctx.JSON(&model.ResponseDTO{
// 		Message: "success",
// 	})
// }
