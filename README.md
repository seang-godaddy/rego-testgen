# Rego Policy Unit Test Generator

## Still a WIP!!!

Something out there probably already exists but I'm here to make a worse version of it.

## TODO LIST:
- get working the base version
- function to convert policy to object


## Usage
You have a policy as such
```
authz {
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
}
```

put it into this program
then it generates unit tests for it like such

```
test_authz_customer_id_match_false_fail {
	not authz with common.customer_id_match as false
		 with hbi.is_dry_run as false
		 with hbi.required as false
		 with sso.low_risk as false
		 with sso.high_risk as false
		 with hbi.authorized as false
}
```

but with every possible input combination of rego conditions you have. For example
there are 6 unique checks here, so it would generate every 2^6 tests (64 tests), every possible permutation
then PRINT IT OUT