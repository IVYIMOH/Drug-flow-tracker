package models

type StockEntry struct {
	HospitalID     int     `json:"hospital_id
	DrugName       string  `json:"drug_name"`
	Source         string  `json:"source"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`

}


type Dispensation struct {
	HospitalID     int     `json:"hospital_id"`
	DrugName       string  `json:"drug_name"`
	Quantity       int     `json:"quantity"`
}