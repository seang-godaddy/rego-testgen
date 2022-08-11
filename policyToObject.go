package main

import "strings"

// for now, no empty lines or comments allowed

func convertPolicy(s string) policy {
	lines := strings.Split(s, "\n")
	return createPolicy(nil, lines)
}

func createPolicy(p *policy, lines []string) policy {
	if p == nil {
		p = &policy{}
	}

	for index := range lines {
		if strings.ContainsAny(lines[index], "else") {
			createPolicy(p.nextPolicy, lines[index+1:])
			break
		}
		if strings.ContainsAny(lines[index], "{}") {
			continue
		}
		testCondition := strings.TrimSpace(lines[index])
		fillInNode(p.firstNode, testCondition)
	}

	return *p
}

func fillInNode(n *node, s string) {
	if n == nil {
		var exclude bool
		if strings.Contains(s, "not") {
			splitCondition := strings.Split(s, " ")
			// todo: not safe
			s = splitCondition[1]
			exclude = true
		}
		n = &node{
			name:        s,
			excludeNode: exclude,
			nextNode:    nil,
		}
		return
	}
	if n.nextNode == nil {
		fillInNode(n.nextNode, s)
	}
}
