package domain

import "time"

type ExporterRow struct {
	ID           string     `csv:"ID"`
	Amount       float64    `csv:"Monto"`
	Description  string     `csv:"Descripción"`
	WalletName   string     `csv:"Billetera"`
	CategoryName string     `csv:"Categoría"`
	CreatedAt    *time.Time `csv:"Fecha"`
}
