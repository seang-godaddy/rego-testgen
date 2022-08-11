package main

import (
	"fmt"
	"strings"
)

var testsAlreadyNamed = make(map[string]int)

func overwriteTestConditions(test, testCondition string, mode bool) string {
	return fmt.Sprintf("%s \n\t with %s as %t \n", test, testCondition, mode)
}

func overwriteTestName(nodeName, testName string, mode bool) string {
	newNodeName := strings.Split(nodeName, ".")
	if len(newNodeName) <= 1 {
		return ""
	}
	return fmt.Sprintf("_%s_%s_%t", testName, nodeName, mode)
}

func finishNamingTest(testName, finalNodeName string, mode bool) string {
	authzResult := "pass"
	authzExpected := "authz"

	if !mode {
		authzResult = "fail"
		authzExpected = "not_authz"
	}

	testName = fmt.Sprintf("test_%s%s_%s", authzExpected, overwriteTestName(finalNodeName, testName, mode), authzResult)
	newTestName := testName
	if testsAlreadyNamed[testName] != 0 {
		newTestName = fmt.Sprintf("%s_%d", testName, testsAlreadyNamed[testName])
	}
	testsAlreadyNamed[testName]++

	return newTestName
}
