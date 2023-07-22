package controllers

import (
	"os"
	"strings"
	"time"

	"github.com/glbayk/gg-skeleton/models"
	"github.com/glbayk/gg-skeleton/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func (u *UserController) Me(c *fiber.Ctx) error {
	user := c.Locals("user")

	return c.Status(fiber.StatusOK).JSON(res.SuccessResponse("Profile data", user))
}

func (u *UserController) RefreshToken(c *fiber.Ctx) error {
	var token string
	refreshToken := c.Get("Authorization")
	if strings.HasPrefix(refreshToken, "Bearer ") {
		token = strings.TrimPrefix(refreshToken, "Bearer ")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(res.UnauthorizedResponse())
	}

	claims, err := utils.TokenValidate(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(res.UnauthorizedResponse())
	}

	claimsEmail, err := claims.GetSubject()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(res.UnauthorizedResponse())
	}

	atTtl, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	at, err := utils.TokenCreate(atTtl, claimsEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	rtTtl, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	rt, err := utils.TokenCreate(rtTtl, claimsEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	return c.Status(fiber.StatusCreated).JSON(res.SuccessResponse("Successfully refreshed token", fiber.Map{
		"access_token":  at,
		"refresh_token": rt,
	}))
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=255"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=255"`
}

func (u *UserController) ChangePassword(c *fiber.Ctx) error {
	var payload ChangePasswordDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.BadRequestResponse())
	}

	errors := validator.New().Struct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.ValidationErrorResponse(errors))
	}

	user := c.Locals("user").(models.User)
	if err := utils.IsHashOf(user.Password, payload.OldPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Incorrect password"))
	}

	if err := utils.IsHashOf(user.Password, payload.NewPassword); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("New password cannot be the same as old password"))
	}

	hash, err := utils.GetHash(payload.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot hash password"))
	}

	user.Password = hash
	if err := user.Update(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot update user"))
	}

	return c.Status(fiber.StatusOK).JSON(res.SuccessResponse("Password successfully changed", user))
}
