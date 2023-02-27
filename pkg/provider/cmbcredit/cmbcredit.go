package cmbcredit

import (
	"encoding/json"
	"fmt"
	"github.com/deb-sig/double-entry-generator/pkg/io/reader"
	"github.com/deb-sig/double-entry-generator/pkg/ir"
	"log"
)

// CmbCredit is provider for CmbCredit.
type CmbCredit struct {
	Statistics Statistics `json:"statistics,omitempty"`
	LineNum    int        `json:"line_num,omitempty"`
	Orders     []Order    `json:"orders,omitempty"`
}

// New creates a new cmbcredit provider.
func New() *CmbCredit {
	return &CmbCredit{
		Statistics: Statistics{},
		LineNum:    0,
		Orders:     make([]Order, 0),
	}
}

type billData struct {
	Data struct {
		Detail []Order `json:"detail"`
	} `json:"data"`
}

// Translate translates the alipay bill records to IR.
func (c *CmbCredit) Translate(filename string) (*ir.IR, error) {
	log.SetPrefix("[Provider-CmbCredit] ")
	billReader, err := reader.GetReader(filename)
	if err != nil {
		return nil, fmt.Errorf("can't get bill reader, err: %v", err)
	}
	decoder := json.NewDecoder(billReader)
	var billData billData
	err = decoder.Decode(&billData)
	if err != nil {
		return nil, fmt.Errorf("can't parse bill data, err: %v", err)
	}
	c.Orders = billData.Data.Detail
	log.Printf("Finished to parse the file %s", filename)
	return c.convertToIR(), nil
}
