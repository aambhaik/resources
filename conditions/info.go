package condition

// OperatorInfo is the information for an operation
type OperatorInfo struct {
	// Name of the operator
	Name string
	// Description of the operator
	Description string
}

// HasOperatorInfo is an interface for an object that
// has Operator Information
type HasOperatorInfo interface {
	OperatorInfo() *OperatorInfo
}
