package main

import (
	"fmt"
	"strings"
)

var glob int

type node struct {
	name        string
	excludeNode bool
	nextNode    *node
}
type regoPolicy struct {
	firstNode  *node
	nextPolicy *regoPolicy
}

func StartPolicyWriting(p regoPolicy) {
	var test, testName string
	var activeConditions = makeNodeMap(p)

	analyzePolicy(p, test, testName, activeConditions)
}

func analyzePolicy(p regoPolicy, test, testName string, activeConditions map[string]bool) {
	hasNextPolicy := p.nextPolicy != nil

	currentTest, currentTestName, newActiveConditions := analyzeNode(p.firstNode, test, testName, hasNextPolicy, activeConditions)
	if p.nextPolicy != nil {
		analyzePolicy(*p.nextPolicy, currentTest, currentTestName, newActiveConditions)
	}
}

// returns are test, and move to next policy
func analyzeNode(n *node, test, testName string, hasNextPolicy bool, activeConditions map[string]bool) (string, string, map[string]bool) {
	var modes = []bool{true, false}
	for _, mode := range modes {
		currentTest := test
		currentTestName := testName
		currentActiveConditions := newMapCopy(activeConditions)
		//fmt.Printf("\nThe current thing: '%s' : '%s' : '%+v'", currentTest, currentTestName, currentActiveConditions)

		// if condition has already been written in test, skip (done for backlog test filling)
		if truthy := currentActiveConditions[n.name]; !truthy {
			currentTest = overwriteTestConditions(currentTest, n.name, mode)
			currentActiveConditions[n.name] = true
			//fmt.Printf("\nnew map conditions %+v\n", currentActiveConditions)
		}

		if n.nextNode != nil {
			if !mode {
				// go to next policy
				if hasNextPolicy {
					currentTestName = overwriteTestName(n.name, currentTestName, mode)
					return currentTest, currentTestName, currentActiveConditions
				}
				currentTestName, expectedResult := finishNamingTest(currentTestName, n.name, mode)
				backfillNodes(currentTest, currentTestName, expectedResult, currentActiveConditions)
				continue
			}
			// go to next node
			analyzeNode(n.nextNode, currentTest, currentTestName, hasNextPolicy, currentActiveConditions)
		}
		// no more nodes. end of current policy
		currentTestName, expectedResult := finishNamingTest(currentTestName, n.name, mode)
		backfillNodes(currentTest, currentTestName, expectedResult, currentActiveConditions)
	}
	// nothing
	return "", "", nil
}

func backfillNodes(test, testName, expectedResult string, activeConditions map[string]bool) {
	// glob is a counter just for a test
	glob++
	// fmt.Println("\nbackfill called", glob)
	// for condition, added := range activeConditions {
	// 	if !added {
	// 		var modes = []bool{true, false}
	// 		activeConditions[condition] = true
	// 		for _, mode := range modes {
	// 			overwriteTest(test, condition, mode)
	// 			// backfillNodes(test, testName, expectedResult, activeConditions)
	// 		}
	// 	}
	// }
	notAuthzSplit := strings.Split(expectedResult, "_")
	resultRejoin := strings.Join(notAuthzSplit, " ")

	fmt.Println("")
	fmt.Printf(`%s {
		%s%s
	}`, testName, resultRejoin, test)
}

func newMapCopy(original map[string]bool) map[string]bool {
	target := make(map[string]bool, len(original))
	for key, value := range original {
		target[key] = value
	}
	return target
}
