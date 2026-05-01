// utils/price.go
package utils

type PriceCheckResult struct {
	IsHigh  bool
	IsKnown bool
}

// kemsaPrices should eventually be moved to the database
var kemsaPrices = map[string]float64{
	"Amoxicillin":   10,
	"Paracetamol":   5,
	"Metformin":     8,
	"Ciprofloxacin": 15,
}

func CheckHighPrice(drug string, price float64) PriceCheckResult {
	base, exists := kemsaPrices[drug]
	if !exists {
		return PriceCheckResult{IsHigh: false, IsKnown: false}
	}
	return PriceCheckResult{IsHigh: price > base*1.5, IsKnown: true}
}
