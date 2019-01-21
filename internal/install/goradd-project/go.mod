module goradd-project

require (
	github.com/alexedwards/scs v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/goradd/goradd v0.0.0-20190117090026-431482a06bb3
	github.com/gorilla/websocket v1.4.0
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/stretchr/testify v1.3.0
	github.com/trustelem/zxcvbn v0.0.0-20180404134528-5fa769e98b1e
	golang.org/x/crypto v0.0.0-20190103213133-ff983b9c42bc
	golang.org/x/net v0.0.0-20190119204137-ed066c81e75e
	gonum.org/v1/gonum v0.0.0-20190119014124-d54847ab4dca
	goradd-tmp v0.0.0

)

replace goradd-tmp => ../goradd-tmp
