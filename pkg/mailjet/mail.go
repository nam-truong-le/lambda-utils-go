package mailjet

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/pkg/errors"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/pkg/retry"
)

// Send sends email
func Send(ctx context.Context, m mailjet.InfoMessagesV31) error {
	log := logger.FromContext(ctx)

	username, err := ssm.GetParameter(ctx, "/mailjet/username", false)
	if err != nil {
		return errors.Wrap(err, "failed to get mailjet username")
	}
	password, err := ssm.GetParameter(ctx, "/mailjet/password", true)
	if err != nil {
		return errors.Wrap(err, "failed to get mailjet password")
	}

	httpClient := &http.Client{Timeout: 300 * time.Second}
	client := mailjet.NewMailjetClient(username, password)
	client.SetClient(httpClient)
	messages := mailjet.MessagesV31{
		Info: []mailjet.InfoMessagesV31{m},
	}

	return retry.Do(ctx, func() error {
		res, err := client.SendMailV31(&messages)
		if err != nil {
			log.Errorf("failed to send email [%s]: %s", m.Subject, err)
			return errors.Wrap(err, "failed to send email")
		}
		log.Infof("[%d]sent status: %s", len(res.ResultsV31), res.ResultsV31[0].Status)
		if res.ResultsV31[0].Status != "success" {
			log.Errorf("failed to send email [%s] to [%+v]: %s", m.Subject, m.To, res.ResultsV31[0].Status)
			return fmt.Errorf("failed to send email [%s] to [%+v]: %s", m.Subject, m.To, res.ResultsV31[0].Status)
		}
		log.Infof("email [%s] sent to [%+v]", m.Subject, m.To)
		return nil
	}, 3, 100*time.Millisecond)
}
