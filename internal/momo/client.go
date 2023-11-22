package momo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/enzoforreal/momo-api/internal/api"
	"github.com/enzoforreal/momo-api/internal/config"
)

type Client struct {
	httpClient   *http.Client
	Config       *config.Config
	Token        string
	TokenExpires time.Time
}

type PaymentRequest struct {
	CorrelatorId          string                  `json:"correlatorId"`
	PaymentDate           string                  `json:"paymentDate"`
	Name                  string                  `json:"name"`
	CallingSystem         string                  `json:"callingSystem"`
	TransactionType       string                  `json:"transactionType"`
	TargetSystem          string                  `json:"targetSystem"`
	CallbackURL           string                  `json:"callbackURL"`
	QuoteId               string                  `json:"quoteId"`
	Channel               string                  `json:"channel"`
	Description           string                  `json:"description"`
	AuthorizationCode     *string                 `json:"authorizationCode"`
	FeeBearer             string                  `json:"feeBearer"`
	Amount                Money                   `json:"amount"`
	TaxAmount             Money                   `json:"taxAmount"`
	TotalAmount           Money                   `json:"totalAmount"`
	Payer                 Payer                   `json:"payer"`
	Payee                 []Payee                 `json:"payee"`
	PaymentMethod         PaymentMethod           `json:"paymentMethod"`
	Status                string                  `json:"status"`
	StatusDate            string                  `json:"statusDate"`
	AdditionalInformation []AdditionalInformation `json:"additionalInformation"`
	Segment               string                  `json:"segment"`
}

type PaymentResponse struct {
	StatusCode            string `json:"statusCode"`
	ProviderTransactionID string `json:"providerTransactionId"`
	StatusMessage         string `json:"statusMessage"`
	SupportMessage        string `json:"supportMessage"`
	SequenceNo            int    `json:"sequenceNo"`
	FulfillmentStatus     string `json:"fulfillmentStatus"`
	Data                  struct {
		ApprovalId            string `json:"approvalId"`
		TransactionFee        Money  `json:"transactionFee"`
		Discount              Money  `json:"discount"`
		NewBalance            Money  `json:"newBalance"`
		PayerNote             string `json:"payerNote"`
		Status                string `json:"status"`
		CorrelatorId          string `json:"correlatorId"`
		StatusDate            string `json:"statusDate"`
		AdditionalInformation struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"additionalInformation"`
		MetaData []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"metaData"`
		LoyaltyInformation struct {
			GeneratedAmount Money `json:"generatedAmount"`
			ConsumedAmount  Money `json:"consumedAmount"`
			NewBalance      Money `json:"newBalance"`
		} `json:"loyaltyInformation"`
		ExternalCode string `json:"externalCode"`
	} `json:"data"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

type Money struct {
	Amount float64 `json:"amount"`
	Units  string  `json:"units"`
}

type Payer struct {
	PayerIdType         string `json:"payerIdType"`
	PayerId             string `json:"payerId"`
	PayerNote           string `json:"payerNote"`
	PayerName           string `json:"payerName"`
	PayerEmail          string `json:"payerEmail"`
	PayerRef            string `json:"payerRef"`
	PayerSurname        string `json:"payerSurname"`
	IncludePayerCharges bool   `json:"includePayerCharges"`
}

type Payee struct {
	Amount      Money  `json:"amount"`
	TaxAmount   Money  `json:"taxAmount"`
	TotalAmount Money  `json:"totalAmount"`
	PayeeIdType string `json:"payeeIdType"`
	PayeeId     string `json:"payeeId"`
	PayeeNote   string `json:"payeeNote"`
	PayeeName   string `json:"payeeName"`
}

type PaymentMethod struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	ValidFrom   time.Time            `json:"validFrom"`
	ValidTo     time.Time            `json:"validTo"`
	Type        string               `json:"type"`
	Details     PaymentMethodDetails `json:"details"`
}

type PaymentMethodDetails struct {
	BankCard            *BankCardDetails            `json:"bankCard,omitempty"`
	TokenizedCard       *TokenizedCardDetails       `json:"tokenizedCard,omitempty"`
	BankAccountDebit    *BankAccountDebitDetails    `json:"bankAccountDebit,omitempty"`
	BankAccountTransfer *BankAccountTransferDetails `json:"bankAccountTransfer,omitempty"`
	Account             *AccountDetails             `json:"account,omitempty"`
	LoyaltyAccount      *LoyaltyAccountDetails      `json:"loyaltyAccount,omitempty"`
	Bucket              *BucketDetails              `json:"bucket,omitempty"`
	Voucher             *VoucherDetails             `json:"voucher,omitempty"`
	DigitalWallet       *DigitalWalletDetails       `json:"digitalWallet,omitempty"`
	Invoice             *InvoiceDetails             `json:"invoice,omitempty"`
}

type BankCardDetails struct {
}

type TokenizedCardDetails struct {
}

type BankAccountDebitDetails struct {
}

type BankAccountTransferDetails struct {
}

type AccountDetails struct {
}

type LoyaltyAccountDetails struct {
}

type BucketDetails struct {
}

type VoucherDetails struct {
}

type DigitalWalletDetails struct {
}

type InvoiceDetails struct {
}

type AdditionalInformation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) TokenValid() bool {
	return time.Now().Before(c.TokenExpires)
}

func (c *Client) RefreshToken() error {
	token, err := c.GetOAuthToken()
	if err != nil {
		return fmt.Errorf("error refreshing token: %v", err)
	}

	c.Token = token
	c.TokenExpires = time.Now().Add(45 * time.Minute)

	return nil
}

func (c *Client) EnsureValidToken() error {
	if !c.TokenValid() {
		if err := c.RefreshToken(); err != nil {
			return fmt.Errorf("error ensuring valid token: %v", err)
		}
	}

	return nil
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{},
		Config:     cfg,
	}
}

func NewPaymentMethod(name, description, typeValue string, validFrom, validTo time.Time) PaymentMethod {
	return PaymentMethod{
		Name:        name,
		Description: description,
		ValidFrom:   validFrom,
		ValidTo:     validTo,
		Type:        typeValue,
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

func (c *Client) createRequest(method, url string, body []byte) (*http.Request, error) {
	httpReq, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	httpReq.Header.Add("Authorization", "Bearer "+c.Token)
	httpReq.Header.Add("Content-Type", "application/json")
	return httpReq, nil
}

func (req *PaymentRequest) CalculateTotalAmount() {

	req.TotalAmount = Money{
		Amount: req.Amount.Amount + req.TaxAmount.Amount,
		Units:  req.Amount.Units,
	}
}

func (c *Client) SendPaymentRequest(req PaymentRequest) (*http.Response, error) {

	if err := c.EnsureValidToken(); err != nil {
		return nil, fmt.Errorf("could not ensure valid token: %v", err)
	}

	if req.CorrelatorId == "" {
		req.CorrelatorId = api.GenerateCorrelatorID()
	}

	if req.CallbackURL == "" {
		req.CallbackURL = c.Config.Momo.CallbackURL
	}

	if req.TransactionType == "" {
		req.TransactionType = "PAYMENT"
	}

	if req.PaymentMethod.Name == "" {
		defaultPaymentMethod := NewPaymentMethod(
			"Defautlt PaymentMethodName",
			"Default PaymentMethodDescription",
			"Mobile Money",
			time.Now(),
			time.Now().AddDate(0, 1, 0),
		)

		req.PaymentMethod = defaultPaymentMethod
	}

	req.CalculateTotalAmount()

	if len(req.AdditionalInformation) == 0 {
		req.AdditionalInformation = []AdditionalInformation{
			{
				Name:        "DefaultInfoName",
				Description: "DefaultInfoDescription",
			},
		}
	}

	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payment request: %v", err)
	}

	httpReq, err := c.createRequest("POST", c.Config.Momo.ApiEndpoint, requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error sending payment request: %v", err)
	}

	defer resp.Body.Close()

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("API responded with status code %d: %s", resp.StatusCode, paymentResp.StatusMessage)
	}

	return resp, nil

}
