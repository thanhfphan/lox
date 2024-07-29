package resolver

type FunctionType int

const (
	FT_NONE     FunctionType = 1
	FT_FUNCTION FunctionType = 2
	FT_METHOD   FunctionType = 3
)

type ClassType int

const (
	CT_NONE  ClassType = 1
	CT_CLASS ClassType = 2
)
