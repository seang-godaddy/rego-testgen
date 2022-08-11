package main

import (
	"fmt"
	"strings"
)

// TODO: write func that translates rego policy into policy/node object

func main() {
	// Example of policy to generate tests for
	/*
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
	*/
	// Example of what that policy would turn into as an object
	/*
		// example of what would become populated
			lfsPolicyDryRunCVRequired := policy{
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
}

type node struct {
	name        string
	excludeNode bool
	nextNode    *node
}
type policy struct {
	firstNode  node
	nextPolicy *policy
}

// todo: try to get rid of global var
var testsAlreadyNamed = make(map[string]int)

func addTestPhrase(check string, mode bool) string {
	return fmt.Sprintf("with %s as %t \n\t", check, mode)
}

func overwriteTestConditions(test, testCondition string, mode bool) string {
	return fmt.Sprintf("%s %s", test, addTestPhrase(testCondition, mode))
}

func overwriteTestName(nodeName, testName string, mode bool) string {
	newNodeName := strings.Split(nodeName, ".")
	if len(newNodeName) <= 1 {
		return ""
	}
	return fmt.Sprintf("%s_%s_%t", testName, nodeName, mode)
}

func analyzePolicy(p policy, test, testName string, activeConditions map[string]bool) {
	currentTest, testName, nextPolicy, newActiveConditions := analyzeNode(&p.firstNode, test, testName, activeConditions)
	if nextPolicy && p.nextPolicy != nil {
		analyzePolicy(*p.nextPolicy, currentTest, testName, newActiveConditions)
	}
	// todo: backfill all conditions that aren't there

	// print test
	fmt.Printf(`%s {
		%s
	}\n`, testName, currentTest)
}

// returns are test, and move to next policy
func analyzeNode(n *node, test, testName string, activeConditions map[string]bool) (string, string, bool, map[string]bool) {
	var modes = []bool{true, false}
	for _, mode := range modes {
		currentTest := test
		currentTestName := testName
		currentActiveConditions := activeConditions

		// if condition has already been written in test, skip (done for backlog test filling)
		if !currentActiveConditions[n.name] {
			currentTest = overwriteTestConditions(test, n.name, mode)
			currentActiveConditions[n.name] = true
		}

		if n.nextNode != nil {
			if !mode {
				// go to next policy
				currentTestName = overwriteTestName(n.name, testName, mode)
				return currentTest, currentTestName, true, currentActiveConditions
			}
			// go to next node
			analyzeNode(n, test, testName, activeConditions)
		}
		// no more nodes. end of current policy

		// finish naming test
		currentTestName = finishNamingTest(currentTestName, n.name, mode)
		return currentTest, currentTestName, false, currentActiveConditions
	}
	// nothing
	return "", "", false, nil
}

func finishNamingTest(testName, finalNodeName string, mode bool) string {
	authzResult := "pass"
	authzExpected := "authz"

	if !mode {
		authzResult = "fail"
		authzExpected = "not_authz"
	}

	testName = fmt.Sprintf("test_%s_%s_%s", authzExpected, overwriteTestName(finalNodeName, testName, mode), authzResult)
	newTestName := testName
	if testsAlreadyNamed[testName] != 0 {
		newTestName = fmt.Sprintf("%s_%d", testName, testsAlreadyNamed[testName])
	}
	testsAlreadyNamed[testName]++

	return newTestName
}
