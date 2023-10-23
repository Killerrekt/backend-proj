package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/utils"
)

func CreateUser(c *fiber.Ctx) error {
	var createUser models.CreateUser

	if err := c.BodyParser(&createUser); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": false, "message": "Error parsing JSON"})
	}

	validator := validator.New()

	if err := validator.Struct(createUser); err != nil {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in all the required fields"})
	}

	dob, _ := time.Parse("2006-01-02", createUser.DateOfBirth)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(createUser.Password), 10)

	user := models.User{
		FirstName:   createUser.FirstName,
		LastName:    createUser.LastName,
		Email:       createUser.Email,
		Password:    string(hashedPassword),
		Gender:      createUser.Gender,
		DateOfBirth: dob,
		Bio:         createUser.Bio,
		TeamID:      0,
		IsLeader:    false,
		IsApproved:  false,
		PhoneNumber: createUser.PhoneNumber,
		College:     createUser.College,
		Github:      createUser.Github,
		Country:     createUser.Country,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status": true, "message": "Successfully created user",
			"verification_status": true, "roasted": false,
		})
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Find(&users)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Return Successful", "users": users,
	})
}

// func ForgotPassword(c *fiber.Ctx) error {
// 	email := c.Params("email")

// 	var check models.User
// 	database.DB.Find(&check, "email = ?", email)
// 	if check.ID == 0 {
// 		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
// 			"Status": false,
// 			"Error":  "The email address given doesn't exist",
// 		})
// 	}

// 	payload := utils.TokenPayload{
// 		Email:   email,
// 		Role:    "",
// 		Version: 0,
// 	}

// 	resetToken, err := utils.CreateToken(
// 		time.Minute*2,
// 		payload,
// 		utils.REFRESH_TOKEN,
// 		viper.GetString("RESET_SECRET_KEY"),
// 	)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
// 			"Status": false,
// 			"Error":  "Failed to create an JWT token",
// 		})
// 	}

// 	url := fmt.Sprintf("%s%s", viper.GetString("RESET_PASSWORD_URL"), resetToken)
// 	message := fmt.Sprintf("%s\n%s %s\n%s",
// 		"Click the link below to reset your password",
// 		url,
// 		"If this request was not sent by you please report to the concerned authorities",
// 		"This is an auto generated email.",
// 	)

// 	err = utils.SendMail("Password Reset", email, message)

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
// 			"Status": false,
// 			"Error":  "Something went wrong while sending the email",
// 		})
// 	}

// 	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
// 		"Status": true,
// 		"data":   resetToken,
// 	})
// }

// func ResetPassword(c *fiber.Ctx) error {
// 	token := c.Params("Token", "")
//
// 	if token == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid token"})
// 	}
//
// 	type Password struct {
// 		Password     string `json:"password"`
// 		Confirm_pass string `json:"confirm_pass"`
// 	}
//
// 	Token, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("invalid signing method")
// 		}
//
// 		return []byte(viper.GetString("RESET_SECRET_KEY")), nil
// 	})
//
// 	if decoded, ok := Token.Claims.(jwt.MapClaims); ok {
// 		if float64(time.Now().Unix()) > decoded["exp"].(float64) {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"Error": "Token Expired",
// 			})
// 		}
//
// 		email := decoded["email"]
// 		var user models.User
// 		database.DB.Find(&user, "email = ?", email)
// 		if user.ID == 0 {
// 			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
// 				"Error": "The email doesn't exist",
// 			})
// 		}
//
// 		req := new(Password)
// 		if err := c.BodyParser(&req); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
// 				"Error": "Error rose while parsing through the body",
// 			})
// 		}
//
// 		if req.Password != req.Confirm_pass {
// 			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
// 				"Error": "Password and confirm password are not the same",
// 			})
// 		}
//
// 		hashed_password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
// 		user.Password = string(hashed_password)
// 		database.DB.Save(user)
// 		return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
// 			"Message": "The password has been updated",
// 		})
// 	}
// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
// }

func SendResetPasswordOTP(c *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Error parsing JSON",
		})
	}

	if request.Email == "" {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in valid data"})
	}

	otp := rand.Intn(900000) + 100000

	var user models.User
	database.DB.Find(&user, "email = ?", request.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "The user does not exist"})
	}

	if err := database.RedisClient.Set(fmt.Sprintf("reset_password:%s", request.Email),
		fmt.Sprint(otp), time.Minute*10); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Could not set otp"})
	}

	body := fmt.Sprintf("Your otp for reset password is: %d", otp)

	if err := utils.SendMail("Reset Password", body, request.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Error sending email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Email Sent"})
}

func VerifyResetPasswordOTP(c *fiber.Ctx) error {
	var request struct {
		Email       string `json:"email" validate:"required"`
		OTP         int    `json:"otp" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Error parsing JSON",
		})
	}

	validator := validator.New()

	if err := validator.Struct(request); err != nil {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in an email, otp and a new password"})
	}

	var user models.User
	database.DB.Find(&user, "email = ?", request.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "User Not Found"})
	}

	otpStr, err := database.RedisClient.Get(fmt.Sprintf("reset_password:%s", request.Email))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Could not read otp from redis"})
	}

	otp, _ := strconv.Atoi(otpStr)

	if otp != request.OTP {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": false, "message": "Invalid OTP"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Some error occured while hashing"})
	}

	user.Password = string(hashed)
	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Password reset successful",
	})
}

func SendVerifyUserOTP(c *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Error parsing JSON",
		})
	}

	if request.Email == "" {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in an email"})
	}

	otp := rand.Intn(900000) + 100000

	var user models.User
	database.DB.Find(&user, "email = ?", request.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "The user does not exist"})
	}

	if user.IsVerified {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": false, "message": "User already verified"})
	}

	if err := database.RedisClient.Set(fmt.Sprintf("verification_otp:%s", request.Email),
		fmt.Sprint(otp), time.Minute*10); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Some error occured"})
	}

	body := fmt.Sprintf("You otp for verification is: %d", otp)

	if err := utils.SendMail("Verification Request", body, request.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Could not send mail"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Email Sent"})
}

func VerifyUserOTP(c *fiber.Ctx) error {
	var request struct {
		Email string `json:"email" validate:"required"`
		OTP   int    `json:"otp" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": false, "message": "Error parsing JSON"})
	}

	validator := validator.New()
	if err := validator.Struct(request); err != nil {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{"status": false, "message": "Please pass in an email and otp"})
	}

	if request.Email == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": false, "message": "Please give a valid email"})
	}

	var user models.User
	database.DB.Find(&user, "email = ?", request.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "User Not Found"})
	}

	otpStr, err := database.RedisClient.Get(fmt.Sprintf("verification_otp:%s", request.Email))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": false, "message": "Could not read otp from redis"})
	}

	otp, _ := strconv.Atoi(otpStr)

	if otp != request.OTP {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": false, "message": "Invalid OTP"})
	}

	user.IsVerified = true
	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Verified User"})
}

func RoastUser(c *fiber.Ctx) error {
	requestID := c.Params("id")

	userId, err := strconv.Atoi(requestID)
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status": false, "message": "Invalid ID",
		})
	}

	var user models.User
	database.DB.Find(&user, "id = ?", userId)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "User Not Found"})
	}

	user.IsRoasted = true
	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "User roasted"})
}

func RevokeRoast(c *fiber.Ctx) error {
	requestID := c.Params("id")

	userId, err := strconv.Atoi(requestID)
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status": false, "message": "Invalid ID",
		})
	}

	var user models.User
	database.DB.Find(&user, "id = ?", userId)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "User Not Found"})
	}

	user.IsRoasted = false
	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "User roast revoked"})
}
