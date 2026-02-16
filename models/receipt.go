package models

import (
	"correctum-agent/fptr10"
	"strconv"
	"time"
)

type TaxationSystem uint8
type ReceiptOperation uint8

type Receipt struct {
	DeviceId          string
	DriveId           string
	Document          uint32
	Session           uint32
	DocumentInSession uint32
	Created           time.Time
	Sign              uint32
	Type              ReceiptOperation
	BuyerContent      string
	TaxationSystem    TaxationSystem
	Items             []Position
	Payment           Payment
	TotalSum          uint64
}

func (r *Receipt) Do(fptr fptr10.IFptr) error {
	var err error
	var correctionInfo []byte
	if correctionInfo, err = generateCorrectionData(fptr, r.Created); err != nil {
		return err
	}
	fptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, r.Type)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_ELECTRONICALLY, true)
	fptr.SetParam(1173, 1)
	fptr.SetParam(1174, correctionInfo)
	fptr.SetParam(1055, r.TaxationSystem)
	fptr.SetParam(1192, strconv.FormatInt(int64(r.Sign), 10))

	if r.BuyerContent == "" {
		fptr.SetParam(1008, "none")
		fptr.SetParam(1117, "none")
	} else {
		fptr.SetParam(1008, r.BuyerContent)
	}

	if err = fptr.OpenReceipt(); err != nil {
		return err
	}
	for _, pos := range r.Items {
		if err = pos.Do(fptr); err != nil {
			fptr.CancelReceipt()
			return err
		}
	}
	if err = setTotal(fptr, r.TotalSum); err != nil {
		fptr.CancelReceipt()
		return err
	}
	if err = r.Payment.Do(fptr); err != nil {
		fptr.CancelReceipt()
		return err
	}
	if err = fptr.CloseReceipt(); err != nil {
		fptr.CancelReceipt()
		return err
	}
	return nil
}

func setTotal(fptr fptr10.IFptr, total uint64) error {
	fptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, total/100.0)
	return fptr.ReceiptTotal()
}

func generateCorrectionData(fptr fptr10.IFptr, date time.Time) (data []byte, err error) {
	fptr.SetParam(1178, date)
	fptr.SetParam(1179, " ")
	err = fptr.UtilFormTlv()

	data = fptr.GetParamByteArray(fptr10.LIBFPTR_PARAM_TAG_VALUE)
	return
}
