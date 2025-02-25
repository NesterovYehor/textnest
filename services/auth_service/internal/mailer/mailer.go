package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/NesterovYehor/textnest/services/auth_service/config"
	"github.com/go-mail/mail/v2"
)

//go:embed templates/*
var templateFS embed.FS // Embeds the templates folder inside the smtp directory

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// NewMailer initializes a new Mailer instance with SMTP settings.
func NewMailer(opts *config.MailerConfig) *Mailer {
	dialer := mail.NewDialer(opts.Host, opts.Port, opts.Username, opts.Password)
	dialer.Timeout = 5 * time.Second

	return &Mailer{
		dialer: dialer,
		sender: opts.Sender,
	}
}

// Send compiles and sends an email using the provided template and data.
func (m *Mailer) Send(recipient, templateFile string, data any) error {
	// Parse the email template
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template file (%s): %w", templateFile, err)
	}

	// Extract subject, plain text body, and HTML body from the template
	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	plainBody := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(plainBody, "plainBody", data); err != nil {
		return fmt.Errorf("failed to execute plainBody template: %w", err)
	}

	htmlBody := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(htmlBody, "htmlBody", data); err != nil {
		return fmt.Errorf("failed to execute htmlBody template: %w", err)
	}

	// Construct the email message
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Send the email
	if err := m.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

