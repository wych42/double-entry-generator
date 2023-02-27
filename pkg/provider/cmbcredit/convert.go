package cmbcredit

import (
	"github.com/deb-sig/double-entry-generator/pkg/ir"
	"math"
)

// convertToIR convert cmbcredit bills to IR
func (c *CmbCredit) convertToIR() *ir.IR {
	i := ir.New()
	for _, o := range c.Orders {
		irO := ir.Order{
			Peer:           o.Description,
			Item:           o.Description,
			PayTime:        o.BillTime.Time,
			Money:          math.Abs(o.Amount),
			OrderID:        &o.UniqueNo,
			Type:           convertType(o),
			TypeOriginal:   string(o.TransactionType),
			TxTypeOriginal: string(o.TransactionType),
			Method:         o.CardNo,
			Commission:     0,
		}
		irO.Metadata = getMetadata(o)
		i.Orders = append(i.Orders, irO)

	}
	return i
}

// convertType convert to ir.Type
func convertType(o Order) ir.Type {
	var typ ir.Type
	if o.Amount < 0 {
		typ = ir.TypeRecv
	} else if o.Amount > 0 {
		typ = ir.TypeSend
	} else {
		typ = ir.TypeUnknown
	}
	return typ
}

// getMetadata get the metadata (e.g. status, method, category and so on.)	from order.
func getMetadata(o Order) map[string]string {
	// FIXME: hard-coded
	source := "招商银行信用卡"
	return map[string]string{
		"source":      source,
		"payTime":     o.BillTime.Format(localTimeFmt),
		"postingDate": o.PostingDate,
		"orderId":     o.UniqueNo,
		"txType":      string(o.TransactionType),
		"method":      o.CardNo,
	}
}
