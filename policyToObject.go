package main

import (
	"fmt"
	"strings"
)

// for now, no empty lines or comments allowed

func ConvertPolicy(s string) regoPolicy {
	lines := strings.Split(s, "\n")
	return createPolicy(nil, lines)
}

func createPolicy(p *regoPolicy, lines []string) regoPolicy {
	if p == nil {
		p = &regoPolicy{}
	}

	for index := range lines {
		//fmt.Println("\nassessing line: ", strings.TrimSpace(lines[index]))
		if strings.Contains(lines[index], "else") {
			//fmt.Println("countered else")
			newPolicy := createPolicy(p.nextPolicy, lines[index+1:])
			p.nextPolicy = &newPolicy
			break
		}
		if strings.ContainsAny(lines[index], "{}") {
			//fmt.Println("encountered skip line")
			continue
		}
		testCondition := strings.TrimSpace(lines[index])
		p.firstNode = fillInNode(p.firstNode, testCondition)
	}

	return *p
}

func fillInNode(n *node, s string) *node {
	//fmt.Println("working on current node: ", n)
	if n == nil {
		//fmt.Println("empty node encountered")
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
		//fmt.Printf("filling new node: %+v\n", n)
		return n
	}
	//fmt.Println("filling in next node")
	n.nextNode = fillInNode(n.nextNode, s)
	return n
	// only checks first and second nodes currently
}

func printPolicy(p regoPolicy) {
	fmt.Printf("%+v\n", p)
	if p.firstNode != nil {
		printNextNodes(p.firstNode)
	}
	if p.nextPolicy != nil {
		printPolicy(*p.nextPolicy)
	}
}

func printNextNodes(n *node) {
	if n != nil {
		fmt.Printf("%+v\n", n)
	}
	if n.nextNode != nil {
		printNextNodes(n.nextNode)
	}
}
