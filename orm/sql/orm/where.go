package orm

type Op struct {
	Flag   string
	Column string
	Value  interface{}
}

func Eq(col string, val interface{}) *Op {
	return &Op{Flag: "eq", Column: col, Value: val}
}
func Nq(col string, val interface{}) *Op {
	return &Op{Flag: "nq", Column: col, Value: val}
}
func Lt(col string, val interface{}) *Op {
	return &Op{Flag: "lt", Column: col, Value: val}
}
func Le(col string, val interface{}) *Op {
	return &Op{Flag: "le", Column: col, Value: val}
}
func Gt(col string, val interface{}) *Op {
	return &Op{Flag: "gt", Column: col, Value: val}
}
func Ge(col string, val interface{}) *Op {
	return &Op{Flag: "ge", Column: col, Value: val}
}

func And(ops ...Op) *Link {
	return &Link
}

//==============================================================================

type Link struct {
	Flag   string
	Values []interface{}
}

type Or struct {
	Values []interface{}
}
