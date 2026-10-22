package domain

// MailService defines the contract for sending emails.
type MailService interface {
	SendOTP(to, subject, htmlBody string) error
}
