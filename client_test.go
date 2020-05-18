package daowallet_test

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"daowallet"
)

const (
	testAPI    = "https://b2b.test.daowallet.com/api/v2"
	testKey    = "WytPv7tNE4RHtDbERU11AzamY82j4VUz"
	testSecret = "tvmQx3vRN1YFdmdexFQrDoB6lyNnCPpBuLl7kHEC"
)

var (
	mockServer = flag.Bool("mock", true, "run tests against mocked server or test api")
	httpProxy  = flag.String("proxy", "", "http proxy for debug purpose")
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	if *mockServer {
		s := httptest.NewTLSServer(handler)

		cli := &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
					return net.Dial(network, s.Listener.Addr().String())
				},
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		return cli, s.Close
	} else {
		t := http.Transport{}

		if httpProxy != nil && *httpProxy != "" {
			proxyURL, err := url.Parse(*httpProxy)
			if err != nil {
				panic(err)
			}
			t.Proxy = http.ProxyURL(proxyURL)
		}

		return &http.Client{
			Transport: &t,
			Timeout:   10 * time.Minute,
		}, func() {}
	}
}

func TestClient_Addresses(t *testing.T) {
	// arrange
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assertHeader(&r.Header, "Content-Type", "application/json")
		assertHeader(&r.Header, "X-Processing-Key", testKey)
		assertHeaderExist(&r.Header, "X-Processing-Signature")

		w.Write([]byte(`{"data":{"id":211,"address":"3Hg7gCcrjXYd6WoiV8BHw1MMrBCZY64say","currency":"BTC","foreign_id":"user-1250","tag":""}}`))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client := daowallet.NewClient(httpClient, testAPI, testKey, testSecret)

	// act
	res, err := client.Addresses(context.Background(), "user-1250", "BTC")

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedResult := daowallet.Address{
		Address:   "3Hg7gCcrjXYd6WoiV8BHw1MMrBCZY64say",
		Currency:  "BTC",
		ForeignID: "user-1250",
		ID:        211,
		Tag:       "",
	}
	if !reflect.DeepEqual(res, expectedResult) {
		t.Fatalf("got: %v, but expected: %v", res, expectedResult)
	}
}

func TestClient_Withdraw(t *testing.T) {
	// arrange
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assertHeader(&r.Header, "Content-Type", "application/json")
		assertHeader(&r.Header, "X-Processing-Key", testKey)
		assertHeaderExist(&r.Header, "X-Processing-Signature")

		w.Write([]byte(`{"data":{"foreign_id":"user-1250","type":"withdrawal","amount":0.01,"sender_currency":"BTC","receiver_currency":"BTC"}}`))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client := daowallet.NewClient(httpClient, testAPI, testKey, testSecret)

	// act
	res, err := client.Withdraw(context.Background(), "user-1250", 0.01, "BTC", "3AtTjKpqmD8Zr6rvcX9cvACTNxMz3praot")

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedResult := daowallet.Withdrawal{
		ForeignID:        "user-1250",
		Type:             "withdrawal",
		Amount:           0.01,
		SenderCurrency:   "BTC",
		ReceiverCurrency: "BTC",
	}
	if !reflect.DeepEqual(res, expectedResult) {
		t.Fatalf("got: %v, but expected: %v", res, expectedResult)
	}
}

func TestClient_InvoiceNew(t *testing.T) {
	// arrange
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assertHeader(&r.Header, "Content-Type", "application/json")
		assertHeader(&r.Header, "X-Processing-Key", testKey)
		assertHeaderExist(&r.Header, "X-Processing-Signature")

		w.Write([]byte(`{"foreign_id":"eif0Z2bfnkY6WU5mg7gIqTUQBgDs5zWI","status":"created",
		"expired_at":"2020-05-12T19:05:55.057Z","client_amount":1250,"client_currency":"USD",
		"addresses":[{"address":"3LAgvMFh11mvsjYxUzUNGaEWTzmC1nnzPZ","expected_amount":0.15498721,
		"crypto_currency":"BTC","rate_usd":8871.7,"rate_eur":8173.8},{"address":"0x0468bc919B99809155157C7aB101d2eeD84efb37",
		"expected_amount":7.21671128,"crypto_currency":"ETH","rate_usd":190.53,"rate_eur":175.53}]}`))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client := daowallet.NewClient(httpClient, testAPI, testKey, testSecret)

	// act
	res, err := client.InvoiceNew(context.Background(), 1250, "USD")

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedResult := daowallet.Invoice{
		ForeignID:      "eif0Z2bfnkY6WU5mg7gIqTUQBgDs5zWI",
		Status:         "created",
		ExpiredAt:      time.Date(2020, 05, 12, 19, 05, 55, 57*1000*1000, time.UTC),
		ClientAmount:   1250,
		ClientCurrency: "USD",
		Addresses: []struct {
			Address        string  `json:"address"`
			ExpectedAmount float64 `json:"expected_amount"`
			CryptoCurrency string  `json:"crypto_currency"`
			RateUSD        float64 `json:"rate_usd"`
			RateEUR        float64 `json:"rate_eur"`
		}{
			{
				Address:        "3LAgvMFh11mvsjYxUzUNGaEWTzmC1nnzPZ",
				ExpectedAmount: 0.15498721,
				CryptoCurrency: "BTC",
				RateUSD:        8871.7,
				RateEUR:        8173.8,
			},
			{
				Address:        "0x0468bc919B99809155157C7aB101d2eeD84efb37",
				ExpectedAmount: 7.21671128,
				CryptoCurrency: "ETH",
				RateUSD:        190.53,
				RateEUR:        175.53,
			},
		},
	}
	if !reflect.DeepEqual(res, expectedResult) {
		t.Fatalf("got: %v, but expected: %v", res, expectedResult)
	}
}

func TestClient_InvoiceStatus(t *testing.T) {
	// arrange
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assertHeader(&r.Header, "Content-Type", "application/json")

		w.Write([]byte(`{"foreign_id":"KMNTQCMWX8VSowpqGmnwYIDuchusB0B5","status":"created",
		"expired_at":"2020-05-12T19:05:55.057Z","client_amount":1250,"client_currency":"USD",
		"addresses":[{"address":"3LAgvMFh11mvsjYxUzUNGaEWTzmC1nnzPZ","expected_amount":0.15498721,
		"crypto_currency":"BTC","rate_usd":8871.7,"rate_eur":8173.8},{"address":"0x0468bc919B99809155157C7aB101d2eeD84efb37",
		"expected_amount":7.21671128,"crypto_currency":"ETH","rate_usd":190.53,"rate_eur":175.53}]}`))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client := daowallet.NewClient(httpClient, testAPI, testKey, testSecret)

	// act
	res, err := client.InvoiceStatus(context.Background(), "KMNTQCMWX8VSowpqGmnwYIDuchusB0B5")

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedResult := daowallet.Invoice{
		ForeignID:      "KMNTQCMWX8VSowpqGmnwYIDuchusB0B5",
		Status:         "created",
		ExpiredAt:      time.Date(2020, 05, 12, 19, 05, 55, 57*1000*1000, time.UTC),
		ClientAmount:   1250,
		ClientCurrency: "USD",
		Addresses: []struct {
			Address        string  `json:"address"`
			ExpectedAmount float64 `json:"expected_amount"`
			CryptoCurrency string  `json:"crypto_currency"`
			RateUSD        float64 `json:"rate_usd"`
			RateEUR        float64 `json:"rate_eur"`
		}{
			{
				Address:        "3LAgvMFh11mvsjYxUzUNGaEWTzmC1nnzPZ",
				ExpectedAmount: 0.15498721,
				CryptoCurrency: "BTC",
				RateUSD:        8871.7,
				RateEUR:        8173.8,
			},
			{
				Address:        "0x0468bc919B99809155157C7aB101d2eeD84efb37",
				ExpectedAmount: 7.21671128,
				CryptoCurrency: "ETH",
				RateUSD:        190.53,
				RateEUR:        175.53,
			},
		},
	}
	if !reflect.DeepEqual(res, expectedResult) {
		t.Fatalf("got: %v, but expected: %v", res, expectedResult)
	}
}

func assertHeader(h *http.Header, name, value string) {
	if h.Get(name) != value {
		panic(fmt.Errorf("got header: %v, but expected: %v", h.Get(name), value))
	}
}

func assertHeaderExist(h *http.Header, name string) {
	if h.Get(name) == "" {
		panic(fmt.Errorf("header `%v` not present or have empty value", name))
	}
}
