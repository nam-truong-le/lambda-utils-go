package ses

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/samber/lo"
	"gopkg.in/gomail.v2"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

const (
	headerFrom      = "From"
	headerTo        = "To"
	headerReplyTo   = "Reply-To"
	headerBCC       = "BCC"
	headerCC        = "CC"
	headerSubject   = "Subject"
	contentTypeHTML = "text/html"
)

type SendPropsAttachment struct {
	Name string
	Data []byte
}

type SendProps struct {
	From        string
	To          []string
	ReplyTo     string
	BCC         []string
	CC          []string
	Subject     string
	HTML        string
	Attachments []SendPropsAttachment
}

func Send(ctx context.Context, props SendProps) error {
	log := logger.FromContext(ctx)
	log.Infof("Send email [%s] to [%s]", props.Subject, props.To)
	msg := gomail.NewMessage()
	msg.SetHeader(headerFrom, props.From)
	msg.SetHeader(headerTo, props.To...)
	msg.SetHeader(headerReplyTo, props.ReplyTo)
	msg.SetHeader(headerBCC, props.BCC...)
	msg.SetHeader(headerCC, props.CC...)
	msg.SetHeader(headerSubject, props.Subject)
	msg.SetBody(contentTypeHTML, props.HTML)
	lo.ForEach(props.Attachments, func(a SendPropsAttachment, _ int) {
		msg.Attach(a.Name, gomail.SetCopyFunc(func(writer io.Writer) error {
			_, err := writer.Write(a.Data)
			return err
		}))
	})
	buf := new(bytes.Buffer)
	_, err := msg.WriteTo(buf)
	if err != nil {
		log.Errorf("Error write email to buffer: %v", err)
		return err
	}

	sesClient, err := NewClient(ctx)
	if err != nil {
		log.Errorf("Error create ses client: %v", err)
		return err
	}
	_, err = sesClient.SendRawEmail(ctx, &ses.SendRawEmailInput{
		RawMessage: &types.RawMessage{
			Data: buf.Bytes(),
		},
	})
	if err != nil {
		log.Errorf("Error send email: %v", err)
	}
	return err
}
