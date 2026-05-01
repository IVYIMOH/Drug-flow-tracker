// models/models.go
package models

type SupplierSource string

const (
	SourceKEMSA   SupplierSource = "KEMSA"
	SourcePrivate SupplierSource = "PRIVATE"
)

func (s SupplierSource) IsValid() bool {
	return s == SourceKEMSA || s == SourcePrivate
}

type StockEntry struct {
	HospitalID int            `json:"hospital_id"`
	DrugName   string         `json:"drug_name"`
	Source     SupplierSource `json:"source"`
	Quantity   int            `json:"quantity"`
	UnitPrice  float64        `json:"unit_price"`
}

type Dispensation struct {
	HospitalID int    `json:"hospital_id"`
	DrugName   string `json:"drug_name"`
	Quantity   int    `json:"quantity"`
}
