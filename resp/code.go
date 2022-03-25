package resp

//go:generate stringer -type=Code -trimprefix=Code
type Code int

const (
	CodeSuccess Code = iota + 200
)

const (
	CodeBaseErr Code = iota + 500 // Basic error
	CodeParamErr
	CodeNotLogin // You are not login
)
