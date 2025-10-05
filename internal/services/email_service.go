package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"soundtube/pkg"

	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	logger *pkg.CustomLogger
	dialer *gomail.Dialer
	addr   string
	from   string
}

func NewEmailService(dialer *gomail.Dialer, addr, from string, logger *pkg.CustomLogger) *EmailService {
	return &EmailService{dialer: dialer, addr: addr, from: from, logger: logger}
}

func (s *EmailService) SendVerificationEmail(ctx context.Context, email, verifyToken string) error {
	ctx, span := s.logger.GetTracer().Start(ctx, "EmailService.SendVerificationEmail")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("token", verifyToken),
	)

	verifyLink := fmt.Sprintf(s.addr+"verify-email?token=%s", email)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Verify Your Email</title>
		</head>
		<body>
			<h2>Email Verification</h2>
			<p>Hello,</p>
			<p>Please verify your email address by clicking the button below:</p>
			<p>
				<a href="%s" style="
					background-color: #007bff; 
					color: white; 
					padding: 12px 24px; 
					text-decoration: none; 
					border-radius: 4px; 
					display: inline-block;
				">Verify Email</a>
			</p>
			<p>Or copy and paste this link in your browser:</p>
			<p>%s</p>
			<p>If you didn't create an account, please ignore this email.</p>
			<br>
			<p>Best regards,<br>Your App Team</p>
		</body>
		</html>
	`, verifyLink, verifyLink)

	textBody := fmt.Sprintf(`
		Verify Your Email Address
		
		Please verify your email address by visiting the following link:
		%s
		
		If you didn't create an account, please ignore this email.
		
		Best regards,
		Your App Team
	`, verifyLink)

	messege := gomail.NewMessage()
	messege.SetHeader("From", s.from)
	messege.SetHeader("To", email)
	messege.SetHeader("Subject", "Verify your email address")

	messege.SetBody("text/html", htmlBody)

	messege.AddAlternative("text/plain", textBody)

	if err := s.dialer.DialAndSend(messege); err != nil {
		s.logger.Error("failed to send verification email", err).WithTrace(ctx)
		return err
	}

	s.logger.Info("sending verify email", "email", email, "link", verifyLink).WithTrace(ctx)
	return nil
}

func generateVerifyToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
