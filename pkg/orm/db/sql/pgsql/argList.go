package pgsql

import "fmt"

type argLister interface {
	addArg(i any) string
	args() []any
}

// DB is the goradd driver for postgresql databases.
type argList struct {
	argList []any
}

func (a *argList) addArg(i any) string {
	a.argList = append(a.argList, i)
	return fmt.Sprintf(`$%d`, len(a.argList))
}

func (a *argList) args() []any {
	return a.argList
}
