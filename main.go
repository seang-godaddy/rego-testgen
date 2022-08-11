package main

// todo delete comment when 1.0 is out
/*
// example of what would become populated
	p := policy{
		firstNode: node{
			name: "common.customer_id_match",
			nextNode: &node{
				name: "hbi.is_dry_run",
				nextNode: &node{
					name:        "hbi.required",
					excludeNode: true,
					nextNode: &node{
						name:     "sso.high_risk",
						nextNode: nil,
					},
				},
			},
		},
		nextPolicy: &policy{
			firstNode: node{
				name: "common.customer_id_match",
				nextNode: &node{
					name: "hbi.is_dry_run",
					nextNode: &node{
						name: "sso.low_risk",
						nextNode: &node{
							name:     "sso.low_risk",
							nextNode: nil,
						},
					},
				},
			},
			nextPolicy: &policy{
				firstNode: node{
					name: "common.customer_id_match",
					nextNode: &node{
						name:     "hbi.authorized",
						nextNode: nil,
					},
				},
				nextPolicy: nil,
			},
		},
	}
}
*/

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
	policy := convertPolicy(policyString)
	startPolicyWriting(policy)
}
