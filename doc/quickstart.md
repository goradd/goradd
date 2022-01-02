# Quick Start
## Installation
### Setup Your Environment

1. Make sure go is installed and the bin directory of your go installation is in your execution path. 
To check this, execute ```go version``` on the command line.
If go is correctly installed, you should see a response with the current
version of your go executable.
2. Create a new project directory *outside* of your go installation directory and change your working directory 
to the new directory.

### Install Goradd
1. Execute ```go get -u github.com/goradd/goradd```
1. Execute ```goradd install```

You should now have a new directory in your current directory:
* goradd-project. This is where you will build out your project. You **can** put some
files outside of this path, but goradd will be placing its code generated files
inside of here.

You may also notice a number of executables that were installed in your go /bin directory
that will be used by goradd to build your application.

### Run the app
1. Change your working directory to the goradd-project directory that was created in the prior step. 
2. From the command line, run:
```go run -mod mod goradd-project/main```
You will see a number of messages about additional go packages being installed.
3. Once you see "Launching Server...", point your browser to the following URL. 
`http://localhost:8000/goradd/`

If everything is working fine, you should see the Goradd startup screen. It will lead 
you through some additional configuration steps and get you started building your
application. 


## Configuration
### Database
1. Goradd currently requires a Mysql database. Create a 
database schema to begin with. Don't worry about it being perfect, you
can change it as you understand your project more. Goradd is flexible enough
to handle your changes.

2. Open goradd-project/config/db.go in a text editor and follow the directions there
to input your database credentials for your development computer.

3. Restart your application.

## Code Generation
From the command line, run:
`go generate goradd-project/codegen/cmd/build.go`

## Run Your Application
Whenever you want to run your application locally, change to the goradd-project directory and run:
```go run goradd-project/main```

## Install an IDE
There are a few IDEs and go-friendly editors out there. Here is a quick overview:

* Goland by JetBrains. This commercial editor is very powerful, and is the one we prefer. It is
free for writing open-source software, and there are some discounts for students. The
cost for a commercial development version is very reasonable. Its built-in
source-level debugger is very easy to use.
* Atom with go-plus plugin
* Visual Studio Code
* Eclipse with GoClipse

# Modules
Goradd is module aware. Whenever you run goradd tools, it will look
in the nearest go.mod file to read the current module environment.

## GO 1.16+
GO 1.16 adds a new wrinkle to the module problem. Before this version, GO would automatically
update the go.mod and go.sum files with any missing packages. You could also tell it to 
automatically update to the latest version of everything.

However, in GO 1.16, they made the go.mod file read-only by default. This is great for people
who are trying to carefully control their builds, but for most development, it made life more difficult.
The good news is that there are many ways to deal with the problem:
1. Manually update using `go get` for every single dependency (Ugh).
2. Add the `-mod mod` build flag whenever you are building
3. Add `-mod=mod` to your GOFLAGS environment variable. As in `GOFLAGS=-mod=mod`
4. Add `-mod=mod` to the private GOFLAGS environment that is only for Go. To do this, 
   run `go env -w GOFLAGS=-mod=mod` on your command line. To undo this, run `go env -u GOFLAGS`

See the following for more info:
* [Go wiki on modules](https://github.com/golang/go/wiki/Modules)
