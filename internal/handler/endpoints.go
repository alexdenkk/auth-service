package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Endpoints registration
func (handler *Handler) RegisterEndpoints(e *echo.Echo) {
	// Healthcheck endpoint
	e.GET("/health/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// Auth group
	auth := e.Group("/auth")

	auth.POST("/sign-in/", handler.SignIn)
	auth.POST("/sign-up/", handler.SignUp)
	auth.GET("/self/", handler.GetSelf)
}
