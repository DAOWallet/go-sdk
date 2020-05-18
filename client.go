// Package daowallet implements client sdk to daowallet api
package daowallet

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	signatureHeader string = "X-Processing-Signature"
	keyHeader       string = "X-Processing-Key"

	addressesEndpoint     string = "addresses/take"
	withdrawalEndpoint    string = "withdrawal/crypto"
	invoiceNewEndpoint    string = "invoice/new"
	invoiceStatusEndpoint string = "invoice/status"

	contentTypeHeader string = "Content-Type"
	jsonContentType   string = "application/json"
)

// Client implements
type Client struct {
	api    string       // URL to daowallet server, i.e. https://b2b.daowallet.com/api/v2
	client *http.Client // http client injected by user
	apiKey string       // api key
	secret string       // secret for HMAC SHA512 signature
}

// Address represents user's crypto-address
type Address struct {
	ID        int64  `json:"id"`
	Address   string `json:"address"`
	Currency  string `json:"currency"`
	ForeignID string `json:"foreign_id"`
	Tag       string `json:"tag"`
}

// Withdrawal represents withdrawal operation info
type Withdrawal struct {
	ForeignID        string  `json:"foreign_id"`
	Type             string  `json:"type"`
	Amount           float64 `json:"amount"`
	SenderCurrency   string  `json:"sender_currency"`
	ReceiverCurrency string  `json:"receiver_currency"`
}

// Invoice represents issued invoice
type Invoice struct {
	ForeignID      string    `json:"foreign_id"`
	Status         string    `json:"status"`
	ExpiredAt      time.Time `json:"expired_at"`
	ClientAmount   float64   `json:"client_amount"`
	ClientCurrency string    `json:"client_currency"`
	Addresses      []struct {
		Address        string  `json:"address"`
		ExpectedAmount float64 `json:"expected_amount"`
		CryptoCurrency string  `json:"crypto_currency"`
		RateUSD        float64 `json:"rate_usd"`
		RateEUR        float64 `json:"rate_eur"`
	} `json:"addresses"`
}

// ErrorResp represents detailed info about api error
type ErrorResp struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Err        string `json:"error"`
	Status     string
}

func (e *ErrorResp) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Status)
}

// NewClient creates api client instance with custom daowallet server URL
func NewClient(c *http.Client, api, key, secret string) *Client {
	return &Client{
		client: c,
		api:    api,
		apiKey: key,
		secret: secret,
	}
}

// NewDefaultClient creates api client instance production daowallet server URL: https://b2b.daowallet.com/api/v2
func NewDefaultClient(key, secret string) *Client {
	return &Client{
		api:    "https://b2b.daowallet.com/api/v2",
		apiKey: key,
		secret: secret,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 4 * time.Second,
				ResponseHeaderTimeout: 3 * time.Second,
			},
			Timeout: 10 * time.Minute,
		},
	}
}

// Addresses obtains customer crypto-address
func (c *Client) Addresses(ctx context.Context, foreignID, currency string) (Address, error) {
	reqBody := map[string]string{
		"foreign_id": foreignID,
		"currency":   currency,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return Address{}, fmt.Errorf("request body marshaling error: %w", err)
	}

	addressesURL, err := joinURL(c.api, addressesEndpoint)
	if err != nil {
		return Address{}, fmt.Errorf("request url creating error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addressesURL.String(), bytes.NewBuffer(reqJSON))
	if err != nil {
		return Address{}, fmt.Errorf("request creating error: %w", err)
	}

	sig, err := createHmac(c.secret, reqJSON)
	if err != nil {
		return Address{}, fmt.Errorf("hmac signature creationg error: %w", err)
	}

	req.Header.Set(contentTypeHeader, jsonContentType)
	req.Header.Set(keyHeader, c.apiKey)
	req.Header.Set(signatureHeader, sig)

	resp, err := c.client.Do(req)
	if err != nil {
		return Address{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	err = ensureSuccessResponse(resp)
	if err != nil {
		return Address{}, fmt.Errorf("request failed: %w", err)
	}

	respBody := struct {
		Data Address `json:"data"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return Address{}, fmt.Errorf("response unmarshaling error: %w", err)
	}

	return respBody.Data, nil
}

// Withdraw withdraws cryptocurrency to the customer crypto address
func (c *Client) Withdraw(ctx context.Context, foreignID string, amount float64, currency, address string) (Withdrawal, error) {
	reqBody := map[string]interface{}{
		"foreign_id": foreignID,
		"amount":     amount,
		"currency":   currency,
		"address":    address,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return Withdrawal{}, fmt.Errorf("request body marshaling error: %w", err)
	}

	withdrawalURL, err := joinURL(c.api, withdrawalEndpoint)
	if err != nil {
		return Withdrawal{}, fmt.Errorf("request url creating error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, withdrawalURL.String(), bytes.NewBuffer(reqJSON))
	if err != nil {
		return Withdrawal{}, fmt.Errorf("request creating error: %w", err)
	}

	sig, err := createHmac(c.secret, reqJSON)
	if err != nil {
		return Withdrawal{}, fmt.Errorf("hmac signature creationg error: %w", err)
	}

	req.Header.Set(contentTypeHeader, jsonContentType)
	req.Header.Set(keyHeader, c.apiKey)
	req.Header.Set(signatureHeader, sig)

	resp, err := c.client.Do(req)
	if err != nil {
		return Withdrawal{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	err = ensureSuccessResponse(resp)
	if err != nil {
		return Withdrawal{}, fmt.Errorf("request failed: %w", err)
	}

	respBody := struct {
		Data Withdrawal `json:"data"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return Withdrawal{}, fmt.Errorf("response unmarshaling error: %w", err)
	}

	return respBody.Data, nil
}

// InvoiceNew issues an invoice to the customer
func (c *Client) InvoiceNew(ctx context.Context, amount float64, fiatCurrency string) (Invoice, error) {
	reqBody := map[string]interface{}{
		"amount":        amount,
		"fiat_currency": fiatCurrency,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return Invoice{}, fmt.Errorf("request body marshaling error: %w", err)
	}

	invoiceNewURL, err := joinURL(c.api, invoiceNewEndpoint)
	if err != nil {
		return Invoice{}, fmt.Errorf("request url creating error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, invoiceNewURL.String(), bytes.NewBuffer(reqJSON))
	if err != nil {
		return Invoice{}, fmt.Errorf("request creating error: %w", err)
	}

	sig, err := createHmac(c.secret, reqJSON)
	if err != nil {
		return Invoice{}, fmt.Errorf("hmac signature creationg error: %w", err)
	}

	req.Header.Set(contentTypeHeader, jsonContentType)
	req.Header.Set(keyHeader, c.apiKey)
	req.Header.Set(signatureHeader, sig)

	resp, err := c.client.Do(req)
	if err != nil {
		return Invoice{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	err = ensureSuccessResponse(resp)
	if err != nil {
		return Invoice{}, fmt.Errorf("request failed: %w", err)
	}

	inv := Invoice{}

	err = json.NewDecoder(resp.Body).Decode(&inv)
	if err != nil {
		return Invoice{}, fmt.Errorf("response unmarshaling error: %w", err)
	}

	return inv, nil
}

// InvoiceStatus obtains status of already issued invoice
func (c *Client) InvoiceStatus(ctx context.Context, id string) (Invoice, error) {
	invoiceStatusURL, err := joinURL(c.api, invoiceStatusEndpoint)
	if err != nil {
		return Invoice{}, fmt.Errorf("request url creating error: %w", err)
	}

	q := invoiceStatusURL.Query()
	q.Add("id", id)
	invoiceStatusURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, invoiceStatusURL.String(), nil)
	if err != nil {
		return Invoice{}, fmt.Errorf("request creating error: %w", err)
	}

	req.Header.Set(contentTypeHeader, jsonContentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return Invoice{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	err = ensureSuccessResponse(resp)
	if err != nil {
		return Invoice{}, fmt.Errorf("request failed: %w", err)
	}

	inv := Invoice{}

	err = json.NewDecoder(resp.Body).Decode(&inv)
	if err != nil {
		return Invoice{}, fmt.Errorf("response unmarshaling error: %w", err)
	}

	return inv, nil
}

func ensureSuccessResponse(resp *http.Response) error {
	statusOK := resp.StatusCode >= 200 && resp.StatusCode <= 299
	if !statusOK {
		errResp := ErrorResp{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return &errResp
	}
	return nil
}

func joinURL(o string, segments ...string) (*url.URL, error) {
	u, err := url.Parse(o)
	if err != nil {
		return nil, err
	}
	for _, s := range segments {
		u.Path = path.Join(u.Path, s)
	}
	return u, nil
}

func createHmac(key string, body []byte) (string, error) {
	h := hmac.New(sha512.New, []byte(key))
	_, err := h.Write(body)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
