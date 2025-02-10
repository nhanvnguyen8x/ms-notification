package infrastructures

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"ms-notification/dtos"
	"net/http"
	"time"
)

type SupabaseClient struct {
	ApiUrl     string
	ApiKey     string
	httpClient *http.Client
}

func NewSupabaseClient(apiUrl string, apiKey string) *SupabaseClient {
	return &SupabaseClient{
		ApiUrl: apiUrl,
		ApiKey: apiKey,
	}
}

func (sc *SupabaseClient) SelectAll(table string) ([]byte, error) {
	var url = fmt.Sprintf("%s/%s?select=*", sc.ApiUrl, table)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Errorf("client: could not create request: %s\n", err)
		return nil, err
	}

	request.Header.Set("apikey", sc.ApiKey)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sc.ApiKey))
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		logrus.Errorf("client: error making http request: %s\n", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logrus.Errorf("client: unexpected status code: %d\n", response.StatusCode)
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	return responseBody, err
}

func (sc *SupabaseClient) Update(table string) ([]byte, error) {
	var url = fmt.Sprintf("%s/%s?select=*", sc.ApiUrl, table)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Errorf("client: could not create request: %s\n", err)
		return nil, err
	}

	request.Header.Set("apikey", sc.ApiKey)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sc.ApiKey))
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		logrus.Errorf("client: error making http request: %s\n", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logrus.Errorf("client: unexpected status code: %d\n", response.StatusCode)
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	return responseBody, err
}

func (sc *SupabaseClient) UpdateScheduleTiming(table string, updateSchedulerTimingRequest *dtos.UpdateSchedulerTimingRequest) ([]byte, error) {
	var url = fmt.Sprintf("%s/%s?id=eq.%d", sc.ApiUrl, table, updateSchedulerTimingRequest.ID)
	bytesData, err := json.Marshal(updateSchedulerTimingRequest)
	if err != nil {
		logrus.Errorf("SendPulseClient.SendTransactionalEmail: could not marshal request: %s\n", err)
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(bytesData))
	if err != nil {
		logrus.Errorf("client: could not create request: %s\n", err)
		return nil, err
	}

	request.Header.Set("apikey", sc.ApiKey)
	request.Header.Set("Prefer", "return=minimal")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sc.ApiKey))
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		logrus.Errorf("client: error making http request: %s\n", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logrus.Infof("client: unexpected status code: %d", response.StatusCode)
	}

	responseBody, err := io.ReadAll(response.Body)
	logrus.Infof("SupabaseClient UpdateSchedulerTimingResponse : %s", string(responseBody))
	return responseBody, err
}

func (sc *SupabaseClient) SelectTaskToday(table string, startDate, endDate string) ([]byte, error) {
	var url = fmt.Sprintf("%s/%s", sc.ApiUrl, table)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Errorf("client: could not create request: %s\n", err)
		return nil, err
	}

	request.Header.Set("apikey", sc.ApiKey)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sc.ApiKey))

	// Build Filter
	var queries = request.URL.Query()
	queries.Add("select", "*")
	queries.Add("created_at", fmt.Sprintf("gte.%s", startDate))
	queries.Add("created_at", fmt.Sprintf("lt.%s", endDate))

	request.URL.RawQuery = queries.Encode()
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		logrus.Errorf("client: error making http request: %s\n", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logrus.Errorf("client: unexpected status code: %d\n", response.StatusCode)
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	return responseBody, err
}
