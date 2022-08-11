package main

import (
	"fmt"
)

type node struct {
	name        string
	excludeNode bool
	nextNode    *node
}
type policy struct {
	firstNode  *node
	nextPolicy *policy
}

func startPolicyWriting(p policy) {
	var test, testName string
	var activeConditions = make(map[string]bool)
	analyzePolicy(p, test, testName, activeConditions)
}

func analyzePolicy(p policy, test, testName string, activeConditions map[string]bool) {
	currentTest, testName, nextPolicy, newActiveConditions := analyzeNode(p.firstNode, test, testName, activeConditions)
	if nextPolicy && p.nextPolicy != nil {
		analyzePolicy(*p.nextPolicy, currentTest, testName, newActiveConditions)
	}
	// todo: backfill all conditions that aren't there
	// ie if 4/6 conditions have been added by end of current test, add versions with the other 2/6 in every permutation of true/false
	// print test - move to backfill func or call at end of it
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
