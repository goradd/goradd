package templates

// This go file uses the .got files in this directory to build a variety of versions of the maps
// in this directory. You can use this as an example of how to use GoT to build your own custom
// versions of the maps here.

//go:generate gengen -c string_joinTreeItem.json  -o ../strslicejointreemap.go github.com/goradd/gengen/templates/map_src/slice_map.tmpl
//go:generate gengen -c string_joinTreeItem.json  -o ../strslicejointreemapi.go github.com/goradd/gengen/templates/map_src/mapi.tmpl
