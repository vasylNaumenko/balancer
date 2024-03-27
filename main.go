package main

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"

	b "balancer/balancer"
)

func main() {
	totalAmount := decimal.NewFromFloat(10000.0)

	proportions := b.TokensMap{
		b.USDT: decimal.NewFromFloat(0.071),
		b.USDC: decimal.NewFromFloat(0.3834),
		b.DAI:  decimal.NewFromFloat(0.0639),
		b.TUSD: decimal.NewFromFloat(0.0639),
		b.FRAX: decimal.NewFromFloat(0.1278),
		b.PAXG: decimal.NewFromFloat(0.29),
	}

	sum := decimal.Zero
	for _, d := range proportions {
		sum = sum.Add(d)
	}
	fmt.Printf("> sum is %v \n", sum)

	balancer := b.NewBalancer(totalAmount, proportions)
	balancer.Amounts = balancer.FillAmounts()
	balancer.CurrentProportions = balancer.GetProportions()

	// pretty json output for balancer
	prettyJSON := PrettyJSON(balancer)
	fmt.Printf("Балансировщик: %v\n", prettyJSON)

	// Симулируем изменение баланса одного из стейблкоинов
	balancer.Amounts[b.USDT] = balancer.Amounts[b.USDT].Add(decimal.NewFromInt(2200))
	balancer.Amounts[b.USDC] = balancer.Amounts[b.USDC].Add(decimal.NewFromInt(300))
	balancer.Amounts[b.DAI] = balancer.Amounts[b.DAI].Add(decimal.NewFromInt(2500))
	balancer.TotalAmount = balancer.TotalAmount.Add(decimal.NewFromInt(5000))

	balancer.CurrentProportions = balancer.GetProportions()
	prettyJSON = PrettyJSON(balancer)
	fmt.Printf("Временный дисбаланс: %v\n", prettyJSON)

	// Перебалансируем
	fmt.Println("Список обменов:")
	exchanges := balancer.Balance()
	for i, e := range exchanges {
		fmt.Printf("%d > обмен %s на %s -> %v \n", i+1, e.From, e.To, e.Amount)
	}
	prettyJSON = PrettyJSON(balancer)
	fmt.Printf("Финальные пропорции: %v\n", prettyJSON)
}

func PrettyJSON(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}
