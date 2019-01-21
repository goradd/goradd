module goradd-test

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/goradd/goradd v0.0.1
	goradd-project v0.0.0
	goradd-tmp v0.0.0
)

replace goradd-project => ../goradd-project // Should be copied to this spot before testing

replace goradd-tmp => ../goradd-tmp
