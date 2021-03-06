// Code generated by "stringer -type=OperationType -linecomment -output=records_string.go"; DO NOT EDIT.

package icloud

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Create-1]
	_ = x[Update-2]
	_ = x[ForceUpdate-3]
	_ = x[Replace-4]
	_ = x[ForceReplace-5]
	_ = x[Delete-6]
	_ = x[ForceDelete-7]
}

const _OperationType_name = "createupdateforceUpdatereplaceforceReplacedeleteforceDelete"

var _OperationType_index = [...]uint8{0, 6, 12, 23, 30, 42, 48, 59}

func (i OperationType) String() string {
	i -= 1
	if i >= OperationType(len(_OperationType_index)-1) {
		return "OperationType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _OperationType_name[_OperationType_index[i]:_OperationType_index[i+1]]
}
