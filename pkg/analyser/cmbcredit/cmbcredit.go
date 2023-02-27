package cmbcredit

import (
	"github.com/deb-sig/double-entry-generator/pkg/config"
	"github.com/deb-sig/double-entry-generator/pkg/ir"
	"github.com/deb-sig/double-entry-generator/pkg/util"
	"log"
)

type CmbCredit struct {
}

// GetAllCandidateAccounts returns all accounts defined in config.
func (c CmbCredit) GetAllCandidateAccounts(cfg *config.Config) map[string]bool {
	// uniqMap will be used to create the concepts.
	uniqMap := make(map[string]bool)
	if cfg.CmbCredit == nil || len(cfg.CmbCredit.Rules) == 0 {
		return uniqMap
	}
	for _, r := range cfg.CmbCredit.Rules {
		if r.MethodAccount != nil {
			uniqMap[*r.MethodAccount] = true
		}
		if r.TargetAccount != nil {
			uniqMap[*r.TargetAccount] = true
		}
	}
	uniqMap[cfg.DefaultPlusAccount] = true
	uniqMap[cfg.DefaultMinusAccount] = true
	return uniqMap
}

// GetAccounts returns minus and plus account.
func (c CmbCredit) GetAccountsAndTags(o *ir.Order, cfg *config.Config, target, provider string) (string, string, map[ir.Account]string, []string) {
	if cfg.CmbCredit == nil || len(cfg.CmbCredit.Rules) == 0 {
		return cfg.DefaultMinusAccount, cfg.DefaultPlusAccount, nil, nil
	}
	resMinus := cfg.DefaultMinusAccount
	resPlus := cfg.DefaultPlusAccount
	var extraAccounts map[ir.Account]string
	var tags = make([]string, 0)

	var err error
	for _, r := range cfg.CmbCredit.Rules {
		match := true
		// get separator
		sep := ","
		if r.Separator != nil {
			sep = *r.Separator
		}

		matchFunc := util.SplitFindContains
		if r.FullMatch {
			matchFunc = util.SplitFindEquals
		}

		if r.Peer != nil {
			match = matchFunc(*r.Peer, o.Peer, sep, match)
		}
		if r.Type != nil {
			match = matchFunc(*r.Type, o.TypeOriginal, sep, match)
		}
		if r.Item != nil {
			match = matchFunc(*r.Item, o.Item, sep, match)
		}
		if r.Method != nil {
			match = matchFunc(*r.Method, o.Method, sep, match)
		}
		if r.Time != nil {
			match, err = util.SplitFindTimeInterval(*r.Time, o.PayTime, match)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
		if r.TimestampRange != nil {
			match, err = util.SplitFindTimeStampInterval(*r.TimestampRange, o.PayTime, match)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}

		if match {
			// Support multiple matches, like one rule matches the
			// minus account, the other rule matches the plus account.
			if r.TargetAccount != nil {
				if o.Type == ir.TypeRecv {
					resMinus = *r.TargetAccount
				} else {
					resPlus = *r.TargetAccount
				}
			}
			if r.MethodAccount != nil {
				if o.Type == ir.TypeRecv {
					resPlus = *r.MethodAccount
				} else {
					resMinus = *r.MethodAccount
				}
			}
		}
	}
	return resMinus, resPlus, extraAccounts, tags
}
