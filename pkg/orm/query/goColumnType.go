package query

// GoColumnType represents the GO type that corresponds to a database column
type GoColumnType int

const (
	ColTypeUnknown GoColumnType = iota
	ColTypeBytes
	ColTypeString
	ColTypeInteger
	ColTypeUnsigned
	ColTypeInteger64
	ColTypeUnsigned64
	ColTypeDateTime
	ColTypeFloat
	ColTypeDouble
	ColTypeBool
)

// String returns the constant type name as a string
func (g GoColumnType) String() string {
	switch g {
	case ColTypeUnknown:
		return "ColTypeUnknown"
	case ColTypeBytes:
		return "ColTypeBytes"
	case ColTypeString:
		return "ColTypeString"
	case ColTypeInteger:
		return "ColTypeInteger"
	case ColTypeUnsigned:
		return "ColTypeUnsigned"
	case ColTypeInteger64:
		return "ColTypeInteger64"
	case ColTypeUnsigned64:
		return "ColTypeUnsigned64"
	case ColTypeDateTime:
		return "ColTypeDateTime"
	case ColTypeFloat:
		return "ColTypeFloat" // always internally represent with max bits
	case ColTypeDouble:
		return "ColTypeDouble" // always internally represent with max bits
	case ColTypeBool:
		return "ColTypeBool" // always internally represent with max bits
	}
	return ""
}

// GoType returns the actual GO type as go code
func (g GoColumnType) GoType() string {
	switch g {
	case ColTypeUnknown:
		return "Unknown"
	case ColTypeBytes:
		return "[]byte"
	case ColTypeString:
		return "string"
	case ColTypeInteger:
		return "int"
	case ColTypeUnsigned:
		return "uint"
	case ColTypeInteger64:
		return "int64"
	case ColTypeUnsigned64:
		return "uint64"
	case ColTypeDateTime:
		return "datetime.DateTime"
	case ColTypeFloat:
		return "float32" // always internally represent with max bits
	case ColTypeDouble:
		return "float64" // always internally represent with max bits
	case ColTypeBool:
		return "bool" // always internally represent with max bits
	}
	return ""
}

// DefaultValue returns a string that represents the GO default value for the corresponding type
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
	return ""
}
