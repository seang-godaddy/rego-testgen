package main

import (
	"fmt"
	"strings"
)

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
	var activeConditions = make(map[string]bool, 0)
	analyzePolicy(p, test, testName, activeConditions)
}

func analyzePolicy(p regoPolicy, test, testName string, activeConditions map[string]bool) {
	currentTest, testName, nextPolicy, newActiveConditions := analyzeNode(p.firstNode, test, testName, activeConditions)
	if nextPolicy && p.nextPolicy != nil {
		analyzePolicy(*p.nextPolicy, currentTest, testName, newActiveConditions)
	}
}

// returns are test, and move to next policy
func analyzeNode(n *node, test, testName string, activeConditions map[string]bool) (string, string, bool, map[string]bool) {
	var modes = []bool{true, false}
	for _, mode := range modes {
		currentTest := test
		currentTestName := testName
		currentActiveConditions := newMapCopy(activeConditions)
		//fmt.Printf("\nThe current thing: '%s' : '%s' : '%+v'", currentTest, currentTestName, currentActiveConditions)

		// if condition has already been written in test, skip (done for backlog test filling)
		if _, found := currentActiveConditions[n.name]; !found {
			// todo: map doesnt get duplicated
			currentTest = overwriteTestConditions(test, n.name, mode)
			currentActiveConditions[n.name] = true
		}

		if n.nextNode != nil {
			if !mode {
				// go to next policy
				currentTestName = overwriteTestName(n.name, currentTestName, mode)
				return currentTest, currentTestName, true, currentActiveConditions
			}
			// go to next node
			analyzeNode(n.nextNode, test, testName, activeConditions)
		}
		// no more nodes. end of current policy
		// finish naming test
		currentTestName, expectedResult := finishNamingTest(currentTestName, n.name, mode)
		backfillNodes(currentTest, currentTestName, expectedResult, activeConditions)
	}
	// nothing
	return "", "", false, nil
}

func backfillNodes(test, testName, expectedResult string, activeConditions map[string]bool) {
	glob++
	fmt.Println("\nbackfill called", glob)
	var visited int
	for condition, added := range activeConditions {
		if !added {
			var modes = []bool{true, false}
			activeConditions[condition] = true
			for _, mode := range modes {
				overwriteTest(test, condition, mode)
				backfillNodes(test, testName, expectedResult, activeConditions)
			}
		} else {
			visited++
		}
	}
	notAuthzSplit := strings.Split(expectedResult, "_")
	resultRejoin := strings.Join(notAuthzSplit, " ")
	if visited == len(activeConditions) {
		fmt.Println("")
		fmt.Printf(`%s {
			%s%s
		}`, testName, resultRejoin, test)
	}
}

var glob int

func newMapCopy(original map[string]bool) map[string]bool {
	target := make(map[string]bool, len(original))
	for key, value := range original {
		target[key] = value
	}
	return target
}
