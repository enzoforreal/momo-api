package momo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/enzoforreal/momo-api/internal/config"
)

type Client struct {
	httpClient *http.Client
	Config     *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{},
		Config:     cfg,
	}
}

func (c *Client) GetOAuthToken() (string, error) {

	credentials := base64.StdEncoding.EncodeToString([]byte(c.Config.Momo.ConsumerKey + ":" + c.Config.Momo.ConsumerSecret))

	requestBody := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest("POST", c.Config.Momo.TokenURL, requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Basic "+credentials)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to the token endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API MOBILE MONEY responded with status code %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return tokenResponse.AccessToken, nil
}
