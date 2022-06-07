package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// IsSupported Currency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool  {
	switch currency {
	case USD , EUR ,CAD :
		return true
	}
	return false
}
