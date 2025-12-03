package mailer

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

type Mailer struct {
	From     string
	Password string
}

func NewMailer(from, password string) *Mailer {
	return &Mailer{
		From: from,
		Password: password,
	}
}

func buildMessage(from, to, subject, body string) string {
	return "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		body
}

func smtpConnect(host, port string) (*smtp.Client, net.Conn, error) {
	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.Dial("tcp", host+":"+port)
	if err != nil {
		return nil, nil, fmt.Errorf("smtp dial timeout error: %w", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, nil, fmt.Errorf("smtp client error: %w", err)
	}

	return client, conn, nil
}

func smtpStartTLS(client *smtp.Client, host string) error {
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: host}

		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("smtp starttls error: %w", err)
		}
	}

	return nil
}

func smtpAuthenticate(client *smtp.Client, from, password, host string) error {
	auth := smtp.PlainAuth("", from, password, host)

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	return nil
}

func smtpSendMessage(client *smtp.Client, from, to, msg string) error {
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp MAIL FROM error: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO error: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA error: %w", err)
	}

	if _, err := writer.Write([]byte(msg)); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("smtp close error: %w", err)
	}

	return nil
}

func (m *Mailer) Send(to, subject, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := buildMessage(m.From, to, subject, body)

	client, conn, err := smtpConnect(smtpHost, smtpPort)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer client.Close()

	if err := smtpStartTLS(client, smtpHost); err != nil {
		return err
	}

	if err := smtpAuthenticate(client, m.From, m.Password, smtpHost); err != nil {
		return err
	}

	if err := smtpSendMessage(client, m.From, to, msg); err != nil {
		return err
	}

	return nil
}