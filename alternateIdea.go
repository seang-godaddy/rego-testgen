package main

import "fmt"

/*
so instead of main.go's method where it goes line by line in policy and populates testName/testConditions on the fly
it might be another idea to do this method instead where a permutation is calculated, then knowledge of each critical fail
condition (ie each node end) would be analyzed afterwards based on the permutation

ie

given
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

the program generates
	with common.customer_id_match as false
	with hbi.is_dry_run as false
	with hbi.required as true
	with sso.low_risk as false
	with sso.high_risk as true
	with hbi.authorized as false

it would then use analysis of
	low_risk - 1
	high_risk - 2
	authorized - 3

to figure out if authz or not authz would be appended as well as an in order generated array
to know which important failures (failures that lead to next policy) would be included in test name

probably would be quicker to code than the other main file
*/

var globalCounter = 0

func temp() {
	// just edit the tests
	tests := []string{
		"common.customer_id_match",
		"hbi.is_dry_run",
		"hbi.required",
		"sso.low_risk",
		"sso.high_risk",
		"hbi.authorized",
	}
	alternateValues(tests, 0, "")
}

func phrase(check string, mode bool) string {
	return fmt.Sprintf("with %s as %t \n\t", check, mode)
}

func overwriteTest(test, testCondition string, mode bool) string {
	return fmt.Sprintf("%s %s", test, phrase(testCondition, mode))
}

func alternateValues(authConditions []string, index int, test string) {
	modes := []bool{true, false}
	for _, mode := range modes {
		newTest := overwriteTest(test, authConditions[index], mode)
		if index < len(authConditions)-1 {
			alternateValues(authConditions, index+1, newTest)
		} else {
			fmt.Printf("test number %d: %s\n", globalCounter, newTest)
			globalCounter++
		}
	}
}
