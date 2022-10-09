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
	policyString := `authz {
			a
			b
			c
		}`
	p := ConvertPolicy(policyString)
	// fmt.Printf("\nnow printing whole policy")
	// printPolicy(p)
	StartPolicyWriting(p)
}

/*
- todo should add fail condition on name of not authz test

Current results are:
test_authz_c_true_pass {
                authz with a as true with b as true with c as true
        }
test_not_authz_c_false_fail {
                not authz with a as true with b as true with c as false
        }
test_authz_b_true_pass {
                authz with a as true with b as true
        }
test_not_authz_b_false_fail {
                not authz with a as true with b as false
        }
test_authz_a_true_pass {
                authz with a as true
        }
test_not_authz_a_false_fail {
                not authz with a as false
        }
*/
