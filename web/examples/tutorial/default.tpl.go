//** This file was code generated by GoT. DO NOT EDIT. ***

package tutorial

import (
	"context"
	"io"
)

func (ctrl *DefaultPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h1>Tutorial</h1>
<p>
Welcome to the GoRADD tutorial. These pages will help you learn about:
</p>
<ul>
<li>The GoRADD Object-Relation Model (ORM), and how to use the generated code to read and write to your database.</li>
<li>The form engine, and how to create dynamic html pages only using GO code.</li>
<li>Additional support libraries</li>
</ul>
<p>
Click on the items on the left to go to the various demonstration pages. To see the source code that generated a page,
click on the <strong>View Source</strong> button in the upper right corner of any page.
</p>
`); err != nil {
		return
	}

	return
}
