package resp

import (
	"fmt"
)

// Size returns the number of bytes the value will occupy when encoded
func Size(v Value) int {
	switch v := v.(type) {
	case *SimpleString:
		// +<string>\r\n
		return 1 + len(v.Bytes) + 2
	case *BulkString:
		if v.Null {
			// $-1\r\n
			return 5
		}
		// $<len>\r\n<data>\r\n
		lenStr := fmt.Sprintf("%d", len(v.Bytes))
		return 1 + len(lenStr) + 2 + len(v.Bytes) + 2
	case *Integer:
		// :<number>\r\n
		numStr := fmt.Sprintf("%d", v.Number)
		return 1 + len(numStr) + 2
	case *Error:
		// -<msg>\r\n
		return 1 + len(v.Msg) + 2
	case *Array:
		if v.Null {
			// *-1\r\n
			return 5
		}
		if len(v.Elements) == 0 {
			// *0\r\n
			return 4
		}
		// *<count>\r\n + elements
		countStr := fmt.Sprintf("%d", len(v.Elements))
		size := 1 + len(countStr) + 2
		for _, elem := range v.Elements {
			size += Size(elem)
		}
		return size
	default:
		return 0
	}
}
