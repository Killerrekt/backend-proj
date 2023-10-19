package utils

import (
	"fmt"
	"math/rand"

	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/database"
	"www.github.com/ic-ETITE-24/icetite-24-backend/internal/models"
)

type OTPPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Email   string `json:"email"`
}

func SendOTP(payload OTPPayload) (int, error) {
	otp := rand.Intn(900000) + 100000

	var user models.User
	database.DB.Find(&user, "email = ?", payload.Email)

	if user.ID == 0 {
		return 0, fmt.Errorf("User not found")
	}

	return otp, SendMail("OTP", payload.Body, payload.Email)
}
