package query

type GoColumnType string

const (
	COL_TYPE_UNKNOWN    GoColumnType = ""
	COL_TYPE_BYTES                   = "[]byte"
	COL_TYPE_STRING                  = "string"
	COL_TYPE_INTEGER                 = "int"
	COL_TYPE_UNSIGNED                = "uint"
	COL_TYPE_INTEGER64               = "int64"
	COL_TYPE_UNSIGNED64              = "uint64"
	COL_TYPE_DATETIME                = "datetime.DateTime"
	COL_TYPE_FLOAT                   = "float32" // always internally represent with max bits
	COL_TYPE_DOUBLE                  = "float64" // always internally represent with max bits
	COL_TYPE_BOOL                    = "bool"
)

func (g GoColumnType) String() string {
	return string(g)
}

func (g GoColumnType) DefaultValue() string {
	switch g {
	case COL_TYPE_UNKNOWN:
		return ""
	case COL_TYPE_BYTES:
		return ""
	case COL_TYPE_STRING:
		return "\"\""
	case COL_TYPE_INTEGER:
		return "0"
	case COL_TYPE_UNSIGNED:
		return "0"
	case COL_TYPE_INTEGER64:
		return "0"
	case COL_TYPE_UNSIGNED64:
		return "0"
	case COL_TYPE_DATETIME:
		return "datetime.DateTime{}"
		/*
			v, _ := goradd.DateTime{}.MarshalText()
			s := string(v[:])
			return fmt.Sprintf("%#v", s)*/

	case COL_TYPE_FLOAT:
		return "0.0" // always internally represent with max bits
	case COL_TYPE_DOUBLE:
		return "0.0" // always internally represent with max bits
	case COL_TYPE_BOOL:
		return "false" // always internally represent with max bits
	}
	return string(g)
}
