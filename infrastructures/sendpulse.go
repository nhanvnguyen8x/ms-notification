package infrastructures

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"ms-notification/configs"
	"ms-notification/dtos"
	"net/http"
	"time"
)

type SendPulseClient struct {
	sendPulseConfig *configs.SendPulseConfig
}

func NewSendPulseClient(sendPulseConfig *configs.SendPulseConfig) *SendPulseClient {
	return &SendPulseClient{
		sendPulseConfig: sendPulseConfig,
	}
}

func (sc *SendPulseClient) GetAccessToken() (string, error) {
	var request = &dtos.SendPulseAccessTokenRequest{
		GrantType:    "client_credentials",
		ClientID:     sc.sendPulseConfig.ClientID,
		ClientSecret: sc.sendPulseConfig.ClientSecret,
	}

	byteData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	result, err := http.Post(sc.sendPulseConfig.AuthUrl, "application/json", bytes.NewBuffer(byteData))
	if err != nil {
		return "", err
	}

	responseBody, err := io.ReadAll(result.Body)
	if err != nil {
		return "", err
	}

	var response = &dtos.SendPulseAccessTokenResponse{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

func (sc *SendPulseClient) SendTransactionalEmail(accessToken string, emailData interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(emailData)
	if err != nil {
		logrus.Errorf("SendPulseClient.SendTransactionalEmail: could not marshal request: %s\n", err)
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, sc.sendPulseConfig.EmailUrl, bytes.NewBuffer(bytesData))
	if err != nil {
		logrus.Errorf("SendPulseClient.SendTransactionalEmail: could not create request: %s\n", err)
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	result, err := client.Do(request)
	if err != nil {
		logrus.Errorf("SendPulseClient.SendTransactionalEmail: could not send request: %s\n", err)
		return nil, err
	}

	if result.StatusCode != http.StatusOK {
		logrus.Infof("SendPulseClient send email failed with status code: %d\n", result.StatusCode)
	}

	responseBody, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (sc *SendPulseClient) GetTemplate(accessToken string) (interface{}, error) {
	request, err := http.NewRequest(http.MethodGet, sc.sendPulseConfig.TemplateUrl, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil, err
	}

	return response, nil
}

func (sc *SendPulseClient) PrepareEmailData(toEmail string) *dtos.SendPulseEmailRequest {
	return &dtos.SendPulseEmailRequest{
		Email: &dtos.SendPulseEmailBody{
			Subject: "Insurance Renewing? Get Cheaper and Better Coverage With TFB",
			From: &dtos.SendPulseEmailAccount{
				Name:  "Growth HQ",
				Email: "tech@growthhq.io",
			},
			To: &dtos.SendPulseEmailAccount{
				Name:  "Steve",
				Email: "thio.thomasks@gmail.com",
			},
			Template: &dtos.SendPulseEmailTemplate{
				ID: 124258,
			},
		},
	}
}

func (sc *SendPulseClient) PrepareEmailHtml(schedulerTiming *dtos.SchedulerTiming, content string, toAccount *dtos.SendPulseEmailAccount) *dtos.SendPulseHtmlEmailRequest {
	var subject = fmt.Sprintf("%s - %s", "Insurance Renewing? Get Cheaper and Better Coverage With TFB", schedulerTiming.Timezone)
	var text = fmt.Sprintf("%s - %s", content, schedulerTiming.StartTime)

	return &dtos.SendPulseHtmlEmailRequest{
		Email: &dtos.SendPulseHtmlEmailBody{
			Html:    "<p>Example text</p>",
			Text:    text,
			Subject: subject,
			From: &dtos.SendPulseEmailAccount{
				Name:  "Growth HQ",
				Email: "tech@growthhq.io",
			},
			To: []*dtos.SendPulseEmailAccount{
				{
					Name:  toAccount.Name,
					Email: toAccount.Email,
				},
				{
					Name:  "Nhan Nguyen",
					Email: "nguyenvunhan00@gmail.com",
				},
			},
		},
	}
}
