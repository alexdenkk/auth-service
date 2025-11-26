package handler

import (
	"alexdenkk/auth-service/internal/domain"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Handler struct
type Handler struct {
	service domain.Service
	logger  *zap.Logger
}

// New handler
func NewHandler(service domain.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Sign-in/sign-up request structure
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User authorization
func (handler *Handler) SignIn(c echo.Context) error {
	start := time.Now()

	handler.logger.Info("handling signin request",
		zap.String("layer", "http"),
	)

	// Parsing data from request
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		handler.logger.Warn("failed to bind signin request",
			zap.String("layer", "http"),
			zap.Error(err),
		)

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Calling service layer and checking error
	token, err := handler.service.SignIn(c.Request().Context(), req.Email, req.Password)

	if err != nil {
		handler.logger.Error("failed to authorize user",
			zap.String("layer", "http"),
			zap.String("email", req.Email),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	handler.logger.Info("user signed in successfully",
		zap.String("layer", "http"),
		zap.String("email", req.Email),
		zap.Duration("duration", time.Since(start)),
	)

	return c.JSON(http.StatusCreated, map[string]string{
		"token": token,
	})
}

// Register user
func (handler *Handler) SignUp(c echo.Context) error {
	start := time.Now()

	handler.logger.Info("handling signup request",
		zap.String("layer", "http"),
	)

	// Parsing data from request
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		handler.logger.Warn("failed to bind signup request",
			zap.String("layer", "http"),
			zap.Error(err),
		)

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Calling service layer and checking error
	if err := handler.service.SignUp(c.Request().Context(), req.Email, req.Password); err != nil {
		handler.logger.Error("failed to register user",
			zap.String("layer", "http"),
			zap.String("email", req.Email),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	handler.logger.Info("user signed up successfully",
		zap.String("layer", "http"),
		zap.String("email", req.Email),
		zap.Duration("duration", time.Since(start)),
	)

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "user registered",
	})
}

// Get user record by JWT token
func (handler *Handler) GetSelf(c echo.Context) error {
	start := time.Now()

	handler.logger.Info("handling get self request",
		zap.String("layer", "http"),
	)

	// Calling service layer
	user, err := handler.service.GetSelf(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)

	// Checking errors
	if err != nil {
		handler.logger.Error("failed to get user",
			zap.String("layer", "http"),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	handler.logger.Info("successfully retrieved self information",
		zap.String("layer", "http"),
		zap.String("email", user.Email),
		zap.Duration("duration", time.Since(start)),
	)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
