package cmbcredit

import (
	"strings"
	"time"
)

const (
	//billTimeFmt 交易记录中的交易时间
	billTimeFmt = "20060102150405"
	//billDateFmt 交易记录中的交易日期、入账日期的格式
	billDateFmt = "20060102"
	//localTimeFmt set time format to utc+8
	localTimeFmt = "2006-01-02 15:04:05 -0700 CST"
)

// Statistics is the Statistics of the bill file.
type Statistics struct {
	UserID          string    `json:"user_id,omitempty"`
	Username        string    `json:"username,omitempty"`
	ParsedItems     int       `json:"parsed_items,omitempty"`
	Start           time.Time `json:"start,omitempty"`
	End             time.Time `json:"end,omitempty"`
	TotalInRecords  int       `json:"total_in_records,omitempty"`
	TotalInMoney    float64   `json:"total_in_money,omitempty"`
	TotalOutRecords int       `json:"total_out_records,omitempty"`
	TotalOutMoney   float64   `json:"total_out_money,omitempty"`
}

type Order struct {
	BillID             string          `json:"billId"`   // 单据序号，在账单列表中的顺序号
	BillType           BillType        `json:"billType"` // 目前只看到DetailBill
	BillTime           PayTime         `json:"billDate"` // 格式为 billTimeFmt
	BillMonth          string          `json:"billMonth"`
	Org                string          `json:"org"`
	TransactionAmount  float64         `json:"transactionAmount,string"` // 交易金额，有符号，表示交易币种金额，可能是外币
	Amount             float64         `json:"amount,string"`            // 交易金额对应的人民币金额
	Description        string          `json:"description"`              // 描述
	PostingDate        string          `json:"postingDate"`              //入账日期
	Location           string          `json:"location"`                 // 地点，目前见到是国家代码
	TotalStages        string          `json:"totalStages"`              // 总分期数
	CurrentStages      string          `json:"currentStages"`            // 当前分期期数
	RemainingStages    string          `json:"remainingStages"`          // 剩余分期期数
	TransactionType    TransactionType `json:"transactionType"`          // 交易类型：目前看到有消费、费用、分期、还款、其他（包含退款）
	CardNo             string          `json:"cardNo"`                   // 卡号后4位
	UniqueNo           string          `json:"uniqueNo"`                 // 唯一编号
	CommonDescFlag     Flag            `json:"commonDescFlag"`           // 目前看到: 还款是N,其他是Y
	RefundTimeHideFlag Flag            `json:"refundTimeHideFlag"`       // 目前看到: 退款是Y，其他是N
}

type BillType string

const (
	BillTypeDetailBill BillType = "DetailBill"
)

type TransactionType string

const (
	TransactionTypeExpense     TransactionType = "消费"
	TransactionTypeFee         TransactionType = "费用"
	TransactionTypeInstallment TransactionType = "分期"
	TransactionTypeRepay       TransactionType = "还款"
	TransactionTypeOther       TransactionType = "其他"
)

type Flag string

const (
	FlagNo  Flag = "N"
	FlagYes Flag = "Y"
)

type PayTime struct {
	time.Time
}

func (ct *PayTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(billTimeFmt, s)
	return
}
