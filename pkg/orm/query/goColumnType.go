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
		return "ColTypeFloat"
	case ColTypeDouble:
		return "ColTypeDouble"
	case ColTypeBool:
		return "ColTypeBool"
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
		return "time.Time"
	case ColTypeFloat:
		return "float32" // always internally represent with max bits
	case ColTypeDouble:
		return "float64" // always internally represent with max bits
	case ColTypeBool:
		return "bool"
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
		return "time.Time{}"
	case ColTypeFloat:
		return "0.0" // always internally represent with max bits
	case ColTypeDouble:
		return "0.0" // always internally represent with max bits
	case ColTypeBool:
		return "false"
	}
	return ""
}

func ColTypeFromGoTypeString(name string) GoColumnType {
	switch name {
	case "Unknown": return ColTypeUnknown
	case "[]byte": return ColTypeBytes
	case "string": return ColTypeString
	case "int": return ColTypeInteger
	case "uint": return ColTypeUnsigned
	case "int64": return ColTypeInteger64
	case "uint64": return ColTypeUnsigned64
	case "time.Time": return ColTypeDateTime
	case "float32": return ColTypeFloat
	case "float64": return ColTypeDouble
	case "bool": return ColTypeBool
	default: panic("unknown column go type " + name)
}
}