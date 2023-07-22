package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/glbayk/gg-skeleton/globals"
	"github.com/glbayk/gg-skeleton/models"
	"github.com/glbayk/gg-skeleton/services"
	"github.com/glbayk/gg-skeleton/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct{}

type RegisterDTO struct {
	Email           string `json:"email" validate:"required,email,min=4,max=255"`
	Password        string `json:"password" validate:"required,min=8,max=255,eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8,max=255"`
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	var payload RegisterDTO

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.BadRequestResponse())
	}

	errors := validator.New().Struct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.ValidationErrorResponse(errors))
	}

	user := models.User{
		Email: payload.Email,
	}

	if user.Find() == nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("User already exists"))
	}

	hash, err := utils.GetHash(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot hash password"))
	}

	user.Password = hash
	user.ActivationToken = utils.RandomString(32)

	err = user.Create()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create user"))
	}

	link := os.Getenv("BASE_URL") + "/api/auth/activate/" + user.ActivationToken
	templatePath := "./services/mailer-templates/confirm123.html"
	m := services.Mailer{
		To:           []string{user.Email},
		Subject:      "GG-Skeleton Confirmation Email",
		TemplatePath: templatePath,
		Variables: map[string]string{
			"link": link,
		},
	}

	globals.GetWaitGroup().Add(1)
	go func() {
		defer globals.GetWaitGroup().Done()
		err = m.SendWithTemplate()
		if err != nil {
			message := services.AmqpConfirmationMessage{
				Email:     user.Email,
				Token:     user.ActivationToken,
				CreatedAt: time.Now(),
			}

			messageToJSON, _ := json.Marshal(message)

			services.AmqpClient.SendMessage(os.Getenv("RABBITMQ_QUEUE_CONFIRMATIONS"), string(messageToJSON))
			lm := services.LogMessageType{
				Message:   string(messageToJSON),
				CreatedAt: time.Now().UTC(),
			}
			err := services.Create(lm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(res.SuccessResponse("User successfully created", user))
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email,min=4,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var payload LoginDTO

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.BadRequestResponse())
	}

	user := models.User{
		Email: payload.Email,
	}

	if err := user.Find(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.NotFoundResponse())
	}

	if user.ActivatedAt == nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("User not activated"))
	}

	if err := utils.IsHashOf(user.Password, payload.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Incorrect password"))
	}

	atTtl, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	at, err := utils.TokenCreate(atTtl, user.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	rtTtl, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	rt, err := utils.TokenCreate(rtTtl, user.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot create token"))
	}

	return c.Status(fiber.StatusCreated).JSON(res.SuccessResponse("User successfully logged in", fiber.Map{
		"access_token":  at,
		"refresh_token": rt,
	}))
}

func (a *AuthController) Activate(c *fiber.Ctx) error {
	token := c.Params("token", "")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Token is required"))
	}

	user := models.User{
		ActivationToken: token,
	}
	if err := user.FindByActivationToken(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.NotFoundResponse())
	}

	now := time.Now().UTC()
	user.ActivatedAt = &now
	user.ActivationToken = ""
	if err := user.Update(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot update user"))
	}

	return c.Status(fiber.StatusOK).JSON(res.SuccessResponse("User successfully activated", user))
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email,min=4,max=255"`
}

func (a *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var payload ForgotPasswordDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.BadRequestResponse())
	}

	user := models.User{
		Email: payload.Email,
	}

	newPass := utils.RandomString(32)
	hash, err := utils.GetHash(newPass)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot hash password"))
	}

	user.Password = hash
	if err := user.Update(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(res.CustomErrorResponse("Cannot update user"))
	}

	templatePath := "./services/mailer-templates/forgot.html"
	m := services.Mailer{
		To:           []string{user.Email},
		Subject:      "GG-Skeleton Forgot Password Email",
		TemplatePath: templatePath,
		Variables: map[string]string{
			"password": newPass,
		},
	}

	globals.GetWaitGroup().Add(1)
	go func() {
		defer globals.GetWaitGroup().Done()
		err = m.SendWithTemplate()
		if err != nil {
			message := services.AmqpForgotPasswordMessage{
				Email:     user.Email,
				CreatedAt: time.Now(),
			}

			messageToJSON, _ := json.Marshal(message)

			services.AmqpClient.SendMessage(os.Getenv("RABBITMQ_QUEUE_FORGOT_PASSWORD"), string(messageToJSON))
			services.MongoClient.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_CONFIRMATIONS")).InsertOne(context.TODO(), message)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(res.SuccessResponse("User successfully updated", user))
}
