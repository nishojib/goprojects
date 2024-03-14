package ecbank

import (
	"encoding/xml"
	"fmt"
	"goprojects/moneyconverter/money"
	"io"
)

// envelope
type envelope struct {
	Rates []currencyRate `xml:"Cube>Cube>Cube"`
}

// currencyRate
type currencyRate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

const baseCurrencyCode = "EUR"

// exchangeRates builds a map of all the supported exchange rates.
func (e envelope) exchangeRates() map[string]float64 {
	rates := make(map[string]float64, len(e.Rates)+1)

	for _, c := range e.Rates {
		rates[c.Currency] = c.Rate
	}

	// add EUR to EUR rate
	rates[baseCurrencyCode] = 1.

	return rates
}

// exchangeRate reads the change rate form the Envelope's contents.
func (e envelope) exchangeRate(source, target string) (money.ExchangeRate, error) {
	if source == target {
		// No change rate for same source and target currencies.
		one, err := money.ParseDecimal("1")
		if err != nil {
			return money.ExchangeRate{}, fmt.Errorf("unable to create a rate of value 1: %w", err)
		}

		return money.ExchangeRate(one), nil
	}

	// rates stores the rates when Envelope is parsed
	rates := e.exchangeRates()

	sourceFactor, sourceFound := rates[source]
	if !sourceFound {
		return money.ExchangeRate{}, fmt.Errorf("failed to find the source currency %s", source)
	}

	targetFactor, targetFound := rates[target]
	if !targetFound {
		return money.ExchangeRate{}, fmt.Errorf("failed to find the target currency %s", target)
	}

	// use a precision of 8 digits after the decimal separator.
	// this should be enough, as most currencies only use 5 digits.
	rate, err := money.ParseDecimal(fmt.Sprintf("%.8f", targetFactor/sourceFactor))
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf(
			"unable to parse exchange rate from %s to %s: %w",
			source,
			target,
			err,
		)
	}

	return money.ExchangeRate(rate), nil
}

func readRateFromResponse(source, target string, respBody io.Reader) (money.ExchangeRate, error) {

	var ecbMessage envelope
	if err := xml.NewDecoder(respBody).Decode(&ecbMessage); err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrUnexpectedFormat, err)
	}

	rate, err := ecbMessage.exchangeRate(source, target)
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrUnexpectedFormat, err)
	}

	return rate, nil
}
