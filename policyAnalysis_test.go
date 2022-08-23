package main

import (
	"testing"
)

func Test_Main(t *testing.T) {
	// todo: get input from caller
	// policyString := `authz {
	// 		common.customer_id_match
	// 		hbi.is_dry_run
	// 		hbi.required
	// 		sso.low_risk
	// 	} else {
	// 		common.customer_id_match
	// 		hbi.is_dry_run
	// 		not hbi.required
	// 		sso.high_risk
	// 	} else {
	// 		common.customer_id_match
	// 		hbi.authorized
	// 	}`
	policyString2 := `authz {
			common.customer_id_match
			hbi.is_dry_run
		}`
	p := ConvertPolicy(policyString2)
	//fmt.Println("\nnow printing whole policy")
	//printPolicy(p)
	StartPolicyWriting(p)
}
