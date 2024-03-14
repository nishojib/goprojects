package ecbank

import (
	"errors"
	"fmt"
	"goprojects/moneyconverter/money"
	"net/http"
	"net/url"
	"time"
)

const (
	ErrCallingServer      = ecbankError("error calling server")
	ErrTimeout            = ecbankError("timed out when waiting for response")
	ErrUnexpectedFormat   = ecbankError("unexpected response format")
	ErrChangeRateNotFound = ecbankError("couldn't find the exchange rate")
	ErrClientSide         = ecbankError("client side error when contacting ECB")
	ErrServerSide         = ecbankError("server side error when contacting ECB")
	ErrUnknownStatusCode  = ecbankError("unknown status code contacting ECB")
)

// Client can call the bank to retrieve exchange rates.
type Client struct {
	client *http.Client
	url    string
}

// NewClient builds a Client that can fetch exchange rates within a given timeout.
func NewClient(timeout time.Duration) Client {
	return Client{
		client: &http.Client{Timeout: timeout},
		url:    "",
	}
}

// FetchExchangeRate fetches the ExchangeRate for the day and returns it.
func (c Client) FetchExchangeRate(
	source, target money.Currency,
) (money.ExchangeRate, error) {
	const euroxrefURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	if c.url == "" {
		c.url = euroxrefURL
	}

	resp, err := c.client.Get(c.url)
	if err != nil {
		var urlErr *url.Error
		if ok := errors.As(err, &urlErr); ok && urlErr.Timeout() {
			return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrTimeout, err.Error())
		}
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrCallingServer, err.Error())
	}
	defer resp.Body.Close()

	if err = checkStatusCode(resp.StatusCode); err != nil {
		return money.ExchangeRate{}, err
	}

	rate, err := readRateFromResponse(source.Code(), target.Code(), resp.Body)
	if err != nil {
		return money.ExchangeRate{}, err
	}

	return rate, nil
}

const (
	clientErrorClass = 4
	serverErrorClass = 5
)

// checkStatusCode returns a different error depending on the returned status code.
func checkStatusCode(statusCode int) error {
	switch {
	case statusCode == http.StatusOK:
		return nil
	case httpStatusClass(statusCode) == clientErrorClass:
		// errors 4xx
		return fmt.Errorf("%w: %d", ErrClientSide, statusCode)
	case httpStatusClass(statusCode) == serverErrorClass:
		// errors 5xx
		return fmt.Errorf("%w: %d", ErrServerSide, statusCode)
	default:
		// any other use cases
		return fmt.Errorf("%w: %d", ErrUnknownStatusCode, statusCode)
	}
}

// httpStatusClass returns the class of a http status code.
func httpStatusClass(statusCode int) int {
	const httpErrorClassSize = 100
	return statusCode / httpErrorClassSize
}
