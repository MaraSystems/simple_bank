package util

const (
	USD = "USD"
	EUR = "EUR"
	NGN = "NGN"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, NGN:
		return true
	}
	return false
}
