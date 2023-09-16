package utils

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/noirbizarre/gonja"
)

func SendEmail(subject string, from string, recipients []string, ctx gonja.Context, template string) error {

	ctx["gosass_base_url"] = os.Getenv("gosass_URL")

	view, err := gonja.Must(gonja.FromFile("templates/emails/" + template + ".jinja")).Execute(ctx)

	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx["view"] = view
	master := gonja.Must(gonja.FromFile("templates/emails/master.jinja"))
	tpl, err := master.Execute(ctx)

	if err != nil {
		fmt.Println(err)
		return err
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if from == "" {
		from = os.Getenv("SMTP_FROM_EMAIL")
	}

	message := "From: " + from + "\n"
	message += "To: " + recipients[0] + "\n"
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-version: 1.0;\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\";\r\n"
	message += "Content-Transfer-Encoding: 7bit;\r\n"
	message += "\r\n"
	message += tpl

	if from == "" {
		from = os.Getenv("SMTP_FROM")
	}

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	address := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	fmt.Println(recipients)
	err = smtp.SendMail(address, auth, from, recipients, []byte(message))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
