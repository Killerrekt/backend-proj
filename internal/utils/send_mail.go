package utils

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendMail(subject, body, to string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", viper.GetString("SENDER_MAIL"))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, viper.GetString("SENDER_MAIL"), viper.GetString("SENDER_PASS"))
	//dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} //need to be check later

	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
