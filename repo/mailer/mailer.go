package mailer

import (
	"context"

	"github.com/way11229/simple_merchant/domain"
)

type Mailer struct {
}

func NewMailer() domain.MailerClient {
	return &Mailer{}
}

func (m *Mailer) Send(ctx context.Context, input *domain.MailerClientSendParams) error {
	return nil
}
