package models

import "correctum-agent/fptr10"

type PositionUnit uint8
type TaxRate uint8
type PaymentPosition uint8
type PositionType uint8

type Position struct {
	Name     string
	Type     PositionType
	Payment  PaymentPosition
	TaxRate  TaxRate
	Unit     PositionUnit
	Price    uint64
	Quantity uint64
	Cost     uint64
}

func (p *Position) Do(fptr fptr10.IFptr) error {
	fptr.SetParam(fptr10.LIBFPTR_PARAM_COMMODITY_NAME, p.Name)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_PRICE, p.Price/100.0)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_QUANTITY, p.Quantity/1000.0)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_TAX_TYPE, p.TaxRate)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_POSITION_SUM, p.Cost/100.0)
	fptr.SetParam(1212, p.Type)
	fptr.SetParam(1214, p.Payment)
	fptr.SetParam(2108, p.Unit)
	return fptr.Registration()
}
