// Code generated by "stringer -type=Environment -linecomment -output=environment_string.go"; DO NOT EDIT.

package icloud

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Development-1]
	_ = x[Production-2]
}

const _Environment_name = "developmentproduction"

var _Environment_index = [...]uint8{0, 11, 21}

func (i Environment) String() string {
	i -= 1
	if i >= Environment(len(_Environment_index)-1) {
		return "Environment(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Environment_name[_Environment_index[i]:_Environment_index[i+1]]
}
