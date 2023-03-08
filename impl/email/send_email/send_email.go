package send_email

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

var GmailService *gmail.Service

type SendEmail struct {
	c         echo.Context
	ctx       context.Context
	requestId string
}

func New(c echo.Context) *SendEmail {
	s := new(SendEmail)
	s.c = c
	s.ctx = logger.NewCtx(s.c)
	logger.Info(s.ctx, "SendEmail Initialized")
	return s
}

func (s *SendEmail) OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("GMAIL_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("GMAIL_REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		logger.ErrorMsg(s.ctx, "Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		logger.Info(s.ctx, "Email service is initialized")
	}
}

func (s *SendEmail) parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("impl/email/send_email/%s", templateFileName))
	if err != nil {
		return "", errors.New("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}

func (s *SendEmail) SendEmailOAUTH2(to string, data interface{}, template string) (bool, error) {
	emailBody, err := s.parseTemplate(template, data)
	if err != nil {
		return false, errors.New("unable to parse email template")
	}

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + "Your Upcoming Visit at HOME by Tales Of Paws" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SendEmail) Send(booking, date, slot, ticket, receiver string) {
	s.OAuthGmailService()
	data := struct {
		Booking string
		Date    string
		Slot    string
		Ticket  string
	}{
		Booking: booking,
		Date:    date,
		Slot:    slot,
		Ticket:  ticket,
	}

	status, err := s.SendEmailOAUTH2(receiver, data, "email_template.html")
	if err != nil {
		logger.Error(s.ctx, err)
		logger.ErrorMsg(s.ctx, "Failed: email not sent to %v", receiver)
	}
	if status {
		logger.Info(s.ctx, "Email sent successfully using OAUTH")
	}
}
