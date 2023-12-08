//** This file was code generated by GoT. DO NOT EDIT. ***

package panels

import (
	"context"
	"io"
)

// DrawTemplate draws the content of the matching control's template file.
func (ctrl *DefaultPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<p>
These pages demonstrate the various supported controls that are available in goradd. These
controls fall in to these categories:
</p>
<ul>
<li>Standard html tags, like `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `&lt;input&gt; `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` or `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `&lt;div&gt; `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` tags</li>
<li>Custom controls that let you edit database data types, like an Integer or Date field</li>
<li>Custom controls that let you create database relationships, like one-to-many relationships</li>
<li>Useful controls for dealing with common situations</li>
<li>Controls corresponding to widgets found in supported css/js libraries, like Bootstrap</li>
</ul>
<p>
The goal of these pages is to try to demonstrate many of the options available with these controls
and give you examples of how to use them. They are also used by our testing system to exercise these
controls during our unit tests.
</p>
<p>
Click on the items on the left to go to the various demonstration pages.
</p>
`); err != nil {
		return
	}

	return
}
