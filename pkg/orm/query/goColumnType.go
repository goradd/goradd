package query

type GoColumnType string

const (
	ColTypeUnknown    GoColumnType = ""
	ColTypeBytes                   = "[]byte"
	ColTypeString                  = "string"
	ColTypeInteger                 = "int"
	ColTypeUnsigned                = "uint"
	ColTypeInteger64               = "int64"
	ColTypeUnsigned64              = "uint64"
	ColTypeDateTime                = "datetime.DateTime"
	ColTypeFloat                   = "float32" // always internally represent with max bits
	ColTypeDouble                  = "float64" // always internally represent with max bits
	ColTypeBool                    = "bool"
)

func (g GoColumnType) String() string {
	return string(g)
}

func (g GoColumnType) DefaultValue() string {
	switch g {
	case ColTypeUnknown:
		return ""
	case ColTypeBytes:
		return ""
	case ColTypeString:
		return "\"\""
	case ColTypeInteger:
		return "0"
	case ColTypeUnsigned:
		return "0"
	case ColTypeInteger64:
		return "0"
	case ColTypeUnsigned64:
		return "0"
	case ColTypeDateTime:
		return "datetime.DateTime{}"
		/*
			v, _ := goradd.DateTime{}.MarshalText()
			s := string(v[:])
			return fmt.Sprintf("%#v", s)*/

	case ColTypeFloat:
		return "0.0" // always internally represent with max bits
	case ColTypeDouble:
		return "0.0" // always internally represent with max bits
	case ColTypeBool:
		return "false" // always internally represent with max bits
	}
	return string(g)
}
