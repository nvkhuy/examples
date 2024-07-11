package mailer

import (
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/thaitanloi365/go-sendgrid"
)

var instance *Client

type Client struct {
	*sendgrid.Mailer
	logger *logger.Logger
}

func New(config *config.Configuration) *Client {
	var mailer = sendgrid.New(&sendgrid.Config{
		APIKey:       config.SendgridAPIKey,
		BccAddresses: config.BCCAddresses,
		SenderMail:   config.SenderMail,
		SenderName:   config.SenderName,
	})

	instance = &Client{
		Mailer: mailer,
		logger: logger.New("mailer"),
	}

	return instance
}

func (c *Client) Send(params sendgrid.SendMailParams) error {
	var err = c.SendMail(params)
	if err != nil {
		c.logger.Errorf("Send email error email=%s template_id=%s error=%v", params.Email, params.TemplateID, err)
		return err
	}
	c.logger.Debugf("Send email success email=%s template_id=%s", params.Email, params.TemplateID)
	return nil
}
func GetInstance() *Client {
	if instance == nil {
		panic("Must call New() first")
	}

	return instance
}
