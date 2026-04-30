package utils

func CheckHighPrice(drug string, price float64) bool {
	kemsaPrices := map[string]float64{
		"Amoxicillin": 10,
	}

	base, exists := kemsaPrices[drug]
	if !exists {
		return false
	}
	return price > base*1.5
}
