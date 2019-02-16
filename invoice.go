package goswiftpos

import (
	"strings"
	"time"
)

//Sale defines an Invoice from SwiftPOS POS
type Sale struct {
	ID              int    `json:"Id"`
	TransactionType int    `json:"TransactionType"`
	TransactionDate string `json:"TransactionDate"`
	Items           []Item `json:"Items"`
}

//GetSaleDate will return the sale date, resolving format issues
func (obj *Sale) GetSaleDate() time.Time {
	d := obj.TransactionDate

	if !strings.Contains(d, "T") {
		d = strings.Replace(d, " ", "T", 1)
	}

	if !strings.Contains(d, "Z") {
		d = d + "Z"
	}

	t1, _ := time.Parse(time.RFC3339, d)
	return t1
}

//GetTotal will return the total for the sale
func (obj *Sale) GetTotal() float64 {
	t := 0.00
	for _, item := range obj.Items {
		t += item.TotalInc
	}
	return t
}

//GetTotalTax will return the total tax for the sale
func (obj *Sale) GetTotalTax() float64 {
	t := 0.00
	for _, item := range obj.Items {
		t += item.Tax
	}
	return t
}

//GetTotalTax will return the total tax for the sale
func (obj *Sale) GetTotalDiscount() float64 {
	t := 0.00
	//for _, item := range obj.Items {
	//	t += 0
	//}
	return t
}

//Sales defines a list of sales from SwiftPOS
type Sales struct {
	ID    int    `json:"Id"`
	Name  string `json:"Name"`
	Sales []Sale `json:"Sales"`
}

//Item defines an Item from SwiftPOS Sale
type Item struct {
	InventoryCode string   `json:"InventoryCode"`
	PLU           int      `json:"Plu"`
	Name          string   `json:"Name"`
	Category      ItemName `json:"Category"`
	Group         ItemName `json:"Group"`
	Quantity      float64  `json:"Quantity"`
	Tax           float64  `json:"Tax"`
	TotalInc      float64  `json:"TotalInc"`
	NormalPrice   float64  `json:"NormalPrice"`
}

//ItemName defines an ItemName from SwiftPOS Item
type ItemName struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}
