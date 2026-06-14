package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/example/user-age-api/internal/models"
	"github.com/example/user-age-api/internal/service"
)

type UserHandler struct {
	svc      service.UserService
	validate *validator.Validate
	log      *zap.Logger
}

func NewUserHandler(svc service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: validator.New(),
		log:      log,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warn("validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: err.Error()})
	}

	user, err := h.svc.CreateUser(c.Context(), req)
	if err != nil {
		h.log.Error("failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "failed to create user"})
	}

	h.log.Info("user created", zap.Int32("id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid user id"})
	}

	user, err := h.svc.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse{Error: "user not found"})
		}
		h.log.Error("failed to get user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "failed to get user"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid user id"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warn("validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: err.Error()})
	}

	user, err := h.svc.UpdateUser(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse{Error: "user not found"})
		}
		h.log.Error("failed to update user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "failed to update user"})
	}

	h.log.Info("user updated", zap.Int32("id", user.ID))
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid user id"})
	}

	if err := h.svc.DeleteUser(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse{Error: "user not found"})
		}
		h.log.Error("failed to delete user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "failed to delete user"})
	}

	h.log.Info("user deleted", zap.Int32("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 10))

	result, err := h.svc.ListUsers(c.Context(), page, pageSize)
	if err != nil {
		h.log.Error("failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "failed to list users"})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func parseID(c *fiber.Ctx) (int32, error) {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}
