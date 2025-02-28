package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"time"

	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
	"github.com/NesterovYehor/textnest/services/auth_service/config"
	"github.com/go-mail/mail/v2"
)

//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer  *mail.Dialer
	sender  string
	breaker *middleware.CircuitBreakerMiddleware
}

func NewMailer(cfg *config.Config) *Mailer {
	dialer := mail.NewDialer(cfg.Mailer.Host, cfg.Mailer.Port, cfg.Mailer.Username, cfg.Mailer.Password)
	dialer.Timeout = 5 * time.Second

	return &Mailer{
		dialer:  dialer,
		sender:  cfg.Mailer.Sender,
		breaker: middleware.NewCircuitBreakerMiddleware(*cfg.CBConfig, "Mailer"),
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	operacion := func(ctx context.Context) (any, error) {
		tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template file (%s): %w", templateFile, err)
		}

		subject := new(bytes.Buffer)
		if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
			return nil, fmt.Errorf("failed to execute subject template: %w", err)
		}

		plainBody := new(bytes.Buffer)
		if err := tmpl.ExecuteTemplate(plainBody, "plainBody", data); err != nil {
			return nil, fmt.Errorf("failed to execute plainBody template: %w", err)
		}

		htmlBody := new(bytes.Buffer)
		if err := tmpl.ExecuteTemplate(htmlBody, "htmlBody", data); err != nil {
			return nil, fmt.Errorf("failed to execute htmlBody template: %w", err)
		}

		msg := mail.NewMessage()
		msg.SetHeader("To", recipient)
		msg.SetHeader("From", m.sender)
		msg.SetHeader("Subject", subject.String())
		msg.SetBody("text/plain", plainBody.String())
		msg.AddAlternative("text/html", htmlBody.String())

		if err := m.dialer.DialAndSend(msg); err != nil {
			return nil, fmt.Errorf("failed to send email: %w", err)
		}

		return nil, nil
	}
	if _, err := m.breaker.Execute(context.Background(), operacion); err != nil {
		return err
	}
	return nil
}
