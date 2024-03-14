package money

import "regexp"

// Currency defines the code of a currency and its decimal precision.
type Currency struct {
	code      string
	precision byte
}

// ErrInvalidCurrencyCode is returned when the currency to parse is not a standard 3-letter code.
const ErrInvalidCurrencyCode = Error("invalid currency code")

// ParseCurrency returns the currency associated to a name and may return ErrInvalidCurrencyCode.
func ParseCurrency(code string) (Currency, error) {
	if len(code) != 3 {
		return Currency{}, ErrInvalidCurrencyCode
	}

	valid, err := validateUppercase(code)
	if err != nil || !valid {
		return Currency{}, ErrInvalidCurrencyCode
	}

	switch code {
	case "IRR":
		return Currency{code: code, precision: 0}, nil
	case "CNY", "VND":
		return Currency{code: code, precision: 1}, nil
	case "BHD", "IQD", "KWD", "LYD", "OMR", "TND":
		return Currency{code: code, precision: 3}, nil
	default:
		return Currency{code: code, precision: 2}, nil
	}
}

// String implements Stringer.
func (c Currency) String() string {
	return c.code
}

// Code returns the ISO code for the currency
func (c Currency) Code() string {
	return c.code
}

// validateUppercase validates if a code is all uppercase
func validateUppercase(code string) (bool, error) {
	uppercasePattern, err := regexp.Compile("^[A-Z]+$")
	if err != nil {
		return false, err
	}

	return uppercasePattern.MatchString(code), nil
}
