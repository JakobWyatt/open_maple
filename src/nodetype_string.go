// Code generated by "stringer -type=nodeType"; DO NOT EDIT.

package main

import "strconv"

const _nodeType_name = "statementassignoperatevariableconstantroot"

var _nodeType_index = [...]uint8{0, 9, 15, 22, 30, 38, 42}

func (i nodeType) String() string {
	if i < 0 || i >= nodeType(len(_nodeType_index)-1) {
		return "nodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _nodeType_name[_nodeType_index[i]:_nodeType_index[i+1]]
}
