package api

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	sendgridKey    string
	sendgridClient *sendgrid.Client
	defaultSender  = mail.NewEmail("CalChat", "noreply@calchat.com")
	defaultAPI     = ""
	defaultScheme  = "http"
)

func init() {
	// initialize environmental variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	sendgridKey = os.Getenv("SENDGRID_KEY")
	sendgridClient = sendgrid.NewSendClient(sendgridKey)
}

//SendEmail sends an email to the recipient with the specified subject
func SendEmail(recipient string, subject string, templatePath string, data map[string]interface{}) error {
	// Parse template file and execute with data.
	var html bytes.Buffer
	tmpl, err := template.ParseFiles("./api/templates/" + templatePath)
	if err != nil {
		return err
	}
	tmpl.Execute(&html, data)

	recipientEmail := mail.NewEmail("recipient", recipient)
	plainTextContent := html.String()

	// Construct and send email via Sendgrid.
	message := mail.NewSingleEmail(defaultSender, subject, recipientEmail, plainTextContent, html.String())
	response, err := sendgridClient.Send(message)
	if err != nil {
		return err
	}

	log.Println(response.StatusCode)

	return nil
}
