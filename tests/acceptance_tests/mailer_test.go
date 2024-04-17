package acceptance_tests

import (
	"context"

	"github.com/way11229/simple_merchant/domain"
)

type testMailerClient struct {
	verificationCode string
}

var mailerClient *testMailerClient

func newTestMailerClient() domain.MailerClient {
	mailerClient = &testMailerClient{}
	return mailerClient
}

func (t *testMailerClient) Send(ctx context.Context, input *domain.MailerClientSendParams) error {
	t.verificationCode = input.HtmlContent
	return nil
}

func (t *testMailerClient) GetVerificationCode() string {
	return t.verificationCode
}
