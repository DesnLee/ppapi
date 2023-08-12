package email

import (
	"log"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

var (
	d        *gomail.Dialer
	host     string
	port     int
	username string
	address  string
	password string

	isInit = false
)

func initDialer() {
	if isInit == true {
		return
	}
	host = viper.GetString("EMAIL.SMTP.HOST")
	port = viper.GetInt("EMAIL.SMTP.PORT")
	username = viper.GetString("EMAIL.SMTP.USERNAME")
	address = viper.GetString("EMAIL.SMTP.MAIL_ADDRESS")
	password = viper.GetString("EMAIL.SMTP.PASSWORD")
	d = gomail.NewDialer(host, port, address, password)

	isInit = true
}

func Send(to []string) {
	initDialer()

	m := gomail.NewMessage()
	m.SetAddressHeader("From", address, username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>LJK</b>")

	if err := d.DialAndSend(m); err != nil {
		log.Fatalln(err)
	}
}
