# Quick Start
## Installation
### For Go 1.10 and below:
1. Create a new directory and set your GOPATH environment variable to it, if needed.
1. Make sure the GOPATH/bin directory is in your execution path, or execute commands from there.
1. Execute ```go get github.com/goradd/goradd```
1. Execute ```goradd install```

### For Go 1.11 and above using modules:
1. If you just installed go, make sure your GOPATH/bin directory is in your execution path.
1. Create a new directory *outside* of your GOPATH and cd to that new directory.
1. Execute ```go install github.com/goradd/goradd```
1. Execute ```goradd install```

You should now have two directories in you current directory:
* goradd-project. This is where you will build out your project. You **can** put some
files outside of this path, but goradd will be placing its code generated files
inside of here.
* goradd-tmp. This is a temporary directory goradd uses for code generation. You
will not check this in to your source control.

## Run the app
From the command line, run:
`go run goradd-project/main`

Now point your browser to the following URL. 
`http://localhost:8000/goradd/`

If everything is working fine, you should see the Goradd startup screen. It will lead 
you through some additional configuration steps. 


## Configuration
### Database
1) Goradd currently requires a Mysql database. Create a 
database schema to begin with. Don't worry about it being perfect, you
can change it as you understand your project more. Goradd is flexible enough
to handle your changes.

2) Open goradd-project/config/db.go in a text editor and follow the directions there
to input your database credentials for your development computer.

## Code Generation
From the command line, run:
`go generate goradd-project/codegen/cmd/build.go`

## Run the app
From the command line, run:
`go run goradd-project/main`

