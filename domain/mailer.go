package domain

import (
	"context"
)

type MailerClientSendParams struct {
	Sender           *MailInfo
	Receiver         *MailInfo
	Subject          string
	PlainTextContent string
	HtmlContent      string
}

type MailInfo struct {
	Name    string
	Address string
}

type MailerClient interface {
	Send(ctx context.Context, input *MailerClientSendParams) error
}
