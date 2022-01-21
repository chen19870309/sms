package utils

import (
	gomail "github.com/go-mail/gomail"
)

var m *gomail.Message

func InitEmailSender() {

}

func SendMail(email, subject, body string) error {
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", subject)
	r := gomail.NewPlainDialer("", 0, "", "")
	return r.DialAndSend(m)
}

func SendMailWithCode(email, code string) error {
	return SendMail(email, "账号邮箱激活码", code)
}
