package models

type TreatmentSummary struct {
	Name      string  `json:"name"`
	Opno      string  `json:"opno"`
	OtherID   string  `json:"other_id"`
	Telno     string  `json:"telno"`
	Diagnosis string  `json:"diagnosis"`
	Age       int     `json:"age"`
	Cycle     int     `json:"cycle"`
	Height    int     `json:"height"`
	Weight    int     `json:"weight"`
	Bsa       float64 `json:"bsa"`
	Sex       string  `json:"sex"`
	Regimen   string  `json:"regimen"`
	Goal      string  `json:"goal"`
	StartDate string  `json:"start_date"`
	Drugs     []struct {
		Name   string `json:"name"`
		Dosage string `json:"dosage"`
	} `json:"drugs"`
	Provider string `json:"provider"`
}

type GeneratedPdf struct {
	Message string `json:"message"`
	Url     string `json:"url"`
}
