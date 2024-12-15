package util

const (
	USD = "USD"
	RMB = "RMB"
	GBP = "GBP"
	EUR = "EUR"
)

func IsSupportCurrency(currency string) bool {
	switch currency {
	case USD, RMB, GBP, EUR:
		return true
	}
	return false
}
