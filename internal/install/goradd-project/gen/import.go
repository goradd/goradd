// Package gen contains the code generated database model, forms and panels that make up the initial website
// generated from the database.
//
// Each database generated is placed in its own folder.
// Within each database folder, there will be the following directories:
//   - form: The top level form objects that control the routes used to get to the pages and the content
//     of each page. These files are meant to be copied to the goradd-project/web/form directory for modification
//     and use.
//   - panel: The panels and sub-panels that make up the content of each page. They are meant to be copied
//     to the goradd-project/web/panel directory for modification and use.
//   - panelbase: These are helper objects for the panels that connect the panels to the database. They should
//     be used in-place and not modified.
//   - model: The database model objects that facilitate accessing the database. Within the model directory
//     there are implementation files ending in .go, and base files ending in .base.go. The files are meant
//     to be used in-place, but you can change the implementation files. The base files should not be changed.
//
// To output the directories described above, run the code generated described in the goradd-project/codegen/cmd
// directory.
package gen
