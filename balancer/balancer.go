package balancer

import (
	"github.com/shopspring/decimal"
)

type (
	Balancer struct {
		// Сумма пула ликвидности в долларах
		TotalAmount decimal.Decimal
		// Проценты пула ликвидности в токенах
		ProportionsToKeep map[string]decimal.Decimal

		CurrentProportions map[string]decimal.Decimal
		// пул ликвидности в токенах
		Amounts TokensMap
	}
)

func NewBalancer(totalAmount decimal.Decimal, percentages map[string]decimal.Decimal) *Balancer {
	return &Balancer{
		TotalAmount:       totalAmount,
		ProportionsToKeep: percentages,
	}
}

func (b *Balancer) FillAmounts() TokensMap {
	result := make(TokensMap, len(b.ProportionsToKeep))
	for coin, proportion := range b.ProportionsToKeep {
		result[coin] = b.TotalAmount.Mul(proportion)
	}
	return result
}

func (b *Balancer) Balance() []Exchange {
	result := b.FillAmounts()

	// fill excess and shortages
	excess := make(TokensMap)
	shortage := make(TokensMap)
	for name, amount := range b.Amounts {
		diff := amount.Sub(result[name])
		if diff.LessThan(decimal.Zero) {
			shortage[name] = diff.Abs()
		} else {
			excess[name] = diff.Abs()
		}
	}

	// fill exchanges
	var exchanges []Exchange
	for shortageName, shortageAmount := range shortage {
		for excessName, excessAmount := range excess {
			if excessAmount.Equals(decimal.Zero) {
				continue
			}
			if excessAmount.Sub(shortageAmount).LessThan(decimal.Zero) {
				excess[excessName] = decimal.Zero
				shortage[shortageName].Sub(excessAmount)

				// exchange event
				exchanges = append(exchanges, Exchange{
					From:   excessName,
					To:     shortageName,
					Amount: excessAmount,
				})
			} else {
				excess[excessName] = excess[excessName].Sub(shortageAmount)
				shortage[shortageName].Sub(shortageAmount)

				// exchange event
				exchanges = append(exchanges, Exchange{
					From:   excessName,
					To:     shortageName,
					Amount: shortageAmount,
				})
			}

		}
	}

	b.Amounts = result
	b.CurrentProportions = b.GetProportions()

	return exchanges
}

func (b Balancer) GetProportions() map[string]decimal.Decimal {
	res := make(map[string]decimal.Decimal, len(b.Amounts))
	for coin, amount := range b.Amounts {
		res[coin] = amount.Div(b.TotalAmount)
	}

	return res
}
