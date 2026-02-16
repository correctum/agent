package models

import "correctum-agent/fptr10"

type Payment struct {
	Cash  uint64
	ECash uint64
	Pre   uint64
	Post  uint64
	Other uint64
}

func (p *Payment) Do(fptr fptr10.IFptr) error {
	for _, payment := range []struct {
		pt  uint8
		sum uint64
	}{
		{fptr10.LIBFPTR_PT_CASH, p.Cash},
		{fptr10.LIBFPTR_PT_ELECTRONICALLY, p.ECash},
		{fptr10.LIBFPTR_PT_PREPAID, p.Pre},
		{fptr10.LIBFPTR_PT_CREDIT, p.Post},
		{fptr10.LIBFPTR_PT_OTHER, p.Other},
	} {
		if err := doPayment(fptr, payment.pt, payment.sum); err != nil {
			return err
		}
	}
	return nil
}

func doPayment(fptr fptr10.IFptr, typePayment uint8, sum uint64) error {
	if sum == 0 {
		return nil
	}
	fptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, typePayment)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, sum/100.0)
	return fptr.Payment()
}
