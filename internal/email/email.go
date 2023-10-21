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
		log.Println(err)
	}
}

func SendValidationCode(to, code string) error {
	initDialer()

	m := gomail.NewMessage()
	m.SetAddressHeader("From", address, username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Pocket Purse 验证码")
	m.SetBody("text/html", "您好，您正在注册或登录 Pocket Purse！<br><br> 您的验证码是: <b>"+code+"</b>，有效期为5分钟。<br><br>如果不是您的操作，请忽略此邮件。")

	return d.DialAndSend(m)
}
