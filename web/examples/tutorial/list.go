package tutorial

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
)

type pageRecord struct {
	order     int
	id 	  string
	title string
	f     createFunction
	files []string
}
type pageRecordList []pageRecord

var pages = make(map[string]pageRecordList)

type createFunction func(ctx context.Context, parent page.ControlI) page.ControlI

func (p pageRecordList) Less(i, j int) bool {
	return p[i].order < p[j].order
}

func RegisterTutorialPage(category string, order int, id string, title string, f createFunction, files []string) {
	v, ok := pages[category]
	if !ok {
		pages[category] = pageRecordList{pageRecord{order, id, title, f, files}}
	} else {
		v = append(v, pageRecord{order, id, title, f, files})
		pages[category] = v
	}
}
