package handler

import (
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, status int, err error) error {
	return c.JSON(status, model.NewErrorResponse[any](status, err.Error()))
}
