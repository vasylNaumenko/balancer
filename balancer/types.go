package balancer

import "github.com/shopspring/decimal"

const (
	USDT = "USDT"
	USDC = "USDC"
	DAI  = "DAI"
	TUSD = "TUSD"
	FRAX = "FRAX"
	PAXG = "PAXG"
)

type (
	TokensMap map[string]decimal.Decimal
	Exchange  struct {
		From   string
		To     string
		Amount decimal.Decimal
	}
)
