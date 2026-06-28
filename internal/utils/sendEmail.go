package utils

import (
	"fmt"

	"github.com/HadeedTariq/go-production-grade-api/internal/utils/env"
	"gopkg.in/gomail.v2"
)

// SendEmailResponse matches your TypeScript return type structure
type SendEmailResponse struct {
	Success bool
	Error   error
}

func SendVerificationEmail(to string, magicLink string) SendEmailResponse {
	// 1. Pull environment variables (Replace these with your actual env logic)
	smtpUser := env.GetEnvString("NODE_MAILER_USER", "random")
	smtpPass := env.GetEnvString("NODE_MAILER_PASSWORD", "random")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	// 2. Define the HTML Body template using structural string formatting
	htmlBody := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; max-width: 620px; margin: 0 auto; background-color: #ffffff; border: 1px solid #ececec; border-radius: 12px; overflow: hidden;">
		<div style="padding: 32px;">
			<h1 style="margin: 0; font-size: 28px; color: #111827; font-weight: 700;">
				Welcome to Daily Dev 🚀
			</h1>
			<p style="margin-top: 18px; font-size: 15px; line-height: 1.7; color: #4b5563;">
				Thanks for joining Daily Dev — a place where developers share ideas,
				create squads, post content, collaborate, and grow together.
			</p>
			<p style="font-size: 15px; line-height: 1.7; color: #4b5563;">
				To activate your account and continue, please verify your email
				address using the button below.
			</p>
			<div style="margin: 36px 0; text-align: center;">
				<a href="%s" style="background-color: #111827; color: #ffffff; text-decoration: none; padding: 14px 28px; border-radius: 8px; display: inline-block; font-size: 15px; font-weight: 600;">
					Verify Email
				</a>
			</div>
			<p style="font-size: 14px; color: #6b7280; line-height: 1.6;">
				If the button above does not work, copy and paste this link into your browser:
			</p>
			<div style="background-color: #f9fafb; border: 1px solid #e5e7eb; border-radius: 8px; padding: 14px; word-break: break-word; font-size: 13px; color: #374151;">
				%s
			</div>
			<p style="margin-top: 24px; font-size: 14px; color: #6b7280;">
				This verification link will expire shortly for security reasons.
			</p>
			<hr style="margin: 32px 0; border: none; border-top: 1px solid #f3f4f6;" />
			<p style="font-size: 12px; color: #9ca3af; line-height: 1.6;">
				If you did not create an account on Daily Dev, you can safely ignore this email.
			</p>
			<p style="margin-top: 20px; font-size: 13px; color: #6b7280;">
				— Daily Dev Team <br />
				<span style="font-style: italic;">Connect. Share. Build.</span>
			</p>
		</div>
	</div>`, magicLink, magicLink)

	// 3. Compose the mail message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Daily Dev <%s>", smtpUser))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify Your Daily Dev Account")
	m.SetBody("text/html", htmlBody)

	// 4. Configure the SMTP Dialler (equivalent to nodemailer's transport)
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// 5. Attempt sending the email
	if err := d.DialAndSend(m); err != nil {
		return SendEmailResponse{
			Success: false,
			Error:   err,
		}
	}

	return SendEmailResponse{Success: true}
}
