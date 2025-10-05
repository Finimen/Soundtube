package auth

type IEmailSener interface {
	SendVerificationEmail(email, verifyToken string) error
}
