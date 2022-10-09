package main

import (
	"fmt"
	"strings"
)

var testsAlreadyNamed = make(map[string]int)

func overwriteTestConditions(test, testCondition string, mode bool) string {
	return fmt.Sprintf("%s with %s as %t", test, testCondition, mode)
}

func overwriteTestName(nodeName, testName string, mode bool) string {
	splitNodeName := strings.Split(nodeName, ".")
	if len(splitNodeName) > 1 {
		nodeName = splitNodeName[1]
	}
	return fmt.Sprintf("%s_%s_%t", testName, nodeName, mode)
}

func finishNamingTest(testName, finalNodeName string, mode bool) (string, string) {
	authzResult := "pass"
	authzExpected := "authz"

	if !mode {
		authzResult = "fail"
		authzExpected = "not_authz"
	}

	// todo fix naming of test period in hbi.isblah
	testName = fmt.Sprintf("test_%s%s_%s", authzExpected, overwriteTestName(finalNodeName, testName, mode), authzResult)
	newTestName := testName
	if testsAlreadyNamed[testName] != 0 {
		newTestName = fmt.Sprintf("%s_%d", testName, testsAlreadyNamed[testName])
	}
	testsAlreadyNamed[testName]++

	return newTestName, authzExpected
}
