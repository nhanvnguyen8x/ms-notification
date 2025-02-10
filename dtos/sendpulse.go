package dtos

type SendPulseAccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type SendPulseAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type SendPulseEmailRequest struct {
	Email *SendPulseEmailBody `json:"email"`
}

type SendPulseEmailBody struct {
	Subject  string                  `json:"subject"`
	From     *SendPulseEmailAccount  `json:"from"`
	To       *SendPulseEmailAccount  `json:"to"`
	Template *SendPulseEmailTemplate `json:"template"`
}

type SendPulseHtmlEmailRequest struct {
	Email *SendPulseHtmlEmailBody `json:"email"`
}

type SendPulseHtmlEmailBody struct {
	Html    string                   `json:"html"`
	Text    string                   `json:"text"`
	Subject string                   `json:"subject"`
	From    *SendPulseEmailAccount   `json:"from"`
	To      []*SendPulseEmailAccount `json:"to"`
}

type SendPulseEmailResponse struct {
}

type SendPulseEmailTemplate struct {
	ID int64 `json:"id"`
}

type SendPulseEmailAccount struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
