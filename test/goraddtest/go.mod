module goradd-test

require (
	goradd-project v0.0.0
	goradd-tmp v0.0.0
)

replace goradd-project => ../goradd-project // Should be copied to this spot before testing

replace goradd-tmp => ../goradd-tmp
