package main

import (
	"reflect"
)

type nodeType int

const (
	assign nodeType = iota
	operate
	variable
	constant
	rootNode
)

type tree struct {
	nodes []*tree
	group nodeType
	value string
}

//acceptTokenValue checks if a token's value matches any strings in an array
//if the array is empty
//	it is assumed there are no constraints on the tokens value,
//	and this function returns true.
func acceptTokenValue(val string, array []string) bool {
	if len(array) == 0 {
		return true
	}
	return inArray(val, array)
}

//splitParse gets an array of tokens
//	splits it into two segments seperated by the first token that satisfies both toSplit and tokenValues
//	generates a new subnode for both segments (using subNodeType)
//	and runs recursiveDescentParser on both of these segments
//splitParse returns root, after subnodes have been added
func splitParse(toSplit tokenType, subNodeType nodeType, tokens []token, root tree, tokenValues ...string) tree {
	//the reason we iterate backwards is so that mathematical statements
	//	are parsed correctly. If there are operators of equal prededence,
	//	then leftmost operators are evaluated first, meaning they are
	//	deeper in the parse tree and are parsed last
	//Despite this, the left-most tokens will still be parsed as a LH subnode
	for i := len(tokens) - 1; i >= 0; i-- {
		currentToken := tokens[i]
		if currentToken.group == toSplit && acceptTokenValue(currentToken.value, tokenValues) {
			//create a subnode, build its subtrees, and add it to root
			subNode := tree{group: subNodeType, value: currentToken.value}
			//build left subtree
			subNode = recursiveDescentParser(tokens[0:i], subNode)
			//build right subtree
			subNode = recursiveDescentParser(tokens[i+1:len(tokens)], subNode)

			root.nodes = append(root.nodes, &subNode)
			return root
		}
	}
	return root
}

//parser takes a slice of tokens, generated by tokenizer,
//	and splits the code into statements,
//	which are then parsed by recursiveDescentParser
func parser(tokens []token) (tree, error) {
	//seperate into statements
	root := tree{group: rootNode}
	lastStatementSplit := -1

	for i, currentToken := range tokens {
		if currentToken.group == statementDelim {
			root = recursiveDescentParser(tokens[lastStatementSplit+1:i], root)
			lastStatementSplit = i
		}
	}
	return root, checkAST(root)
}

//recursiveDescentParser parses an array of tokens, given to it by parser, into an AST.
//This AST can then be used to either generate machine code or direct execution.
//This function recursively calls itself to build the AST.
//The second argument to recursiveDescentParser specifies the initial node
func recursiveDescentParser(tokens []token, root tree) tree {
	//this function can be optimized, as it currently iterates through tokens when not needed,
	//	as a side effect of recursion
	//do not parallelize statement parsing,
	//	as currently recursiveDescentParser does not have exclusive access to root

	//parse assignment
	assignmentTree := splitParse(assignment, assign, tokens, root)
	if !reflect.DeepEqual(assignmentTree, root) {
		return assignmentTree
	}

	//at this point, we could optionally check for a well formed statement tree

	//parse mathematical operators
	addsubTree := splitParse(operator, operate, tokens, root, "+", "-")
	if !reflect.DeepEqual(addsubTree, root) {
		return addsubTree
	}
	muldivTree := splitParse(operator, operate, tokens, root, "*", "/")
	if !reflect.DeepEqual(muldivTree, root) {
		return muldivTree
	}

	//parse names and numbers
	nameTree := splitParse(name, variable, tokens, root)
	if !reflect.DeepEqual(nameTree, root) {
		return nameTree
	}
	numTree := splitParse(number, constant, tokens, root)
	if !reflect.DeepEqual(numTree, root) {
		return numTree
	}

	return root
}
