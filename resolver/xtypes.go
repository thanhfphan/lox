package resolver

type FunctionType int

const (
	FT_NONE        FunctionType = 1
	FT_FUNCTION    FunctionType = 2
	FT_INITIALIZER FunctionType = 3
	FT_METHOD      FunctionType = 4
)

type ClassType int

const (
	CT_NONE     ClassType = 1
	CT_CLASS    ClassType = 2
	CT_SUBCLASS ClassType = 3
)
