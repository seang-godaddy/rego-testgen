package main

func main() {
	// todo: get input from caller
	policyString := `authz {
		common.customer_id_match
		hbi.is_dry_run
		hbi.required
		sso.low_risk
	} else {
		common.customer_id_match
		hbi.is_dry_run
		not hbi.required
		sso.high_risk
	} else {
		common.customer_id_match
		hbi.authorized
	}`
	p := ConvertPolicy(policyString)
	StartPolicyWriting(p)
}
