package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

type renderData struct {
	MJML string `json:"mjml"`
}

type renderResponse struct {
	HTML        string        `json:"html"`
	Errors      []interface{} `json:"errors"`
	MJML        string        `json:"mjml"`
	MJMLVersion string        `json:"mjml_version"`
}

func Render(ctx context.Context, tmpStr string, data any, user, pass string) (*string, error) {
	log := logger.FromContext(ctx)
	log.Infof("Render email with data: %+v", data)

	t, err := template.New("send-email-template").Parse(tmpStr)
	if err != nil {
		log.Errorf("Failed to parse template %s: %s", err, tmpStr)
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		log.Errorf("Failed to execute template %s: %s", err, tmpStr)
		return nil, err
	}

	finalMJTmp := buf.String()
	body, err := json.Marshal(renderData{MJML: finalMJTmp})
	if err != nil {
		log.Errorf("Failed to marshal data: %s", err)
		return nil, err
	}

	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodPost, "https://api.mjml.io/v1/render", bytes.NewBuffer(body))
	req.SetBasicAuth(user, pass)
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to send request: %s", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close body: %s", err)
		}
	}(res.Body)
	resBodyString, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Failed to read body: %s", err)
		return nil, err
	}

	resBody := new(renderResponse)
	err = json.Unmarshal(resBodyString, resBody)
	if err != nil {
		log.Errorf("Failed to unmarshal body: %s", err)
		return nil, err
	}

	return &resBody.HTML, nil
}
