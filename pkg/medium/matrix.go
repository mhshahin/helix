/*
To send messages to a channel in Matrix, following project has been
used to setup Slack-compatible webhook for Matrix.

https://github.com/turt2live/matrix-appservice-webhooks
*/

package medium

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mhshahin/helix/pkg/config"
	corev1 "k8s.io/api/core/v1"
)

type MatrixInterface interface {
	SendEvent(event *corev1.Event) error
}

// MatrixMessage contains required data to send to Matrix's webhook
type MatrixMessage struct {
	Text        string `json:"text"`
	Format      string `json:"format"`
	DisplayName string `json:"displayName"`
}

// MatrixClient holds everything that is required
// to create a new Matrix client
type MatrixClient struct {
	restyClient *resty.Client
	url         string
}

// NewMatrixClient returns a new instance of MatrixClient
func NewMatrixClient(restyClient *resty.Client) *MatrixClient {
	reqTimeout, err := time.ParseDuration(config.Cfg.Mediums.Matrix.Timeout)
	if err != nil {
		log.Panicln(err)
	}
	restyClient.SetTimeout(reqTimeout)

	return &MatrixClient{
		restyClient: restyClient,
		url:         fmt.Sprintf("%s/%s", config.Cfg.Mediums.Matrix.Address, config.Cfg.Mediums.Matrix.Token),
	}
}

// Send posts a new message to Matrix's webhook
func (mc *MatrixClient) SendEvent(event *corev1.Event) error {
	var ev = fmt.Sprintf("<blockquote data-mx-border-color=#D63232><h4><b>Name:</b> %s</h4><b>Kind</b><br>%s<br><b>Namespace</b><br>%s<br><b>Reason</b><br>%s<br><b>Message</b><br>%s</blockquote>",
		event.InvolvedObject.Name,
		event.InvolvedObject.Kind,
		event.InvolvedObject.Namespace,
		event.Reason,
		event.Message,
	)

	var msg = &MatrixMessage{
		Text:        ev,
		Format:      "html",
		DisplayName: config.Cfg.Mediums.Matrix.DisplayName,
	}

	resp, err := mc.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(msg).
		Post(mc.url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Println("Marix failed response:", resp.Result())
		return fmt.Errorf("NOT OK")
	}

	return nil
}
