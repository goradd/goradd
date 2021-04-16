# Quick Start
## Installation
### Setup Your Environment
The below instructions mention "modules". If you are new to go or you don't know that
modules are, we recommend installing
the latest version of go, and following the instructions below for go when not using modules.
For more information, see [More On Modules](#more-on-modules) further below.

#### For Go 1.10 and below, or Go 1.11+ not using modules:
1. Create a new directory for your project, set the GOPATH environment variable to it and 
then change your working directory to the new directory.
1. The following steps will create a `bin` directory in your new directory. Before doing this
though, make sure the `GOPATH/bin` directory is in your execution path. On Mac OS or Linux 
you would do that by editing your .bash_profile directory or .profile 
directory and on Windows you would use the System utility. Note that you will need to
do this for every new project.

#### For Go 1.11+ using modules:
1. Make sure your GOPATH/bin directory is in your execution path. With go modules, you
only need to do this once and it will work for all of your projects.
1. Create a new project directory *outside* of your GOPATH and change your working directory 
to the new directory.

### Install Goradd
1. Execute ```go get -u github.com/goradd/goradd```
1. Execute ```goradd install```

You should now have a new directory in you current directory:
* goradd-project. This is where you will build out your project. You **can** put some
files outside of this path, but goradd will be placing its code generated files
inside of here.

You will also notice a number of executables that were installed in your GOPATH/bin directory
that will be used by goradd to build your application.

### Run the app
1. Change your working directory to the goradd-project directory (`cd goradd-project`) and run: `go mod tidy`. This will download all the dependencies of the generated project. 
2. From the command line, run:
```go run goradd-project/main```
You will see a number of messages about additional go packages being installed.
3. Once you see "Launching Server...", point your browser to the URL shown in the output. 
e.g. `http://192.168.29.22/` when the output was `Launching server on 192.168.29.22`.

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

# More on Modules
Since the introduction of go modules in version 1.10, the go build environment has
been in flux. The transition to go modules for some has not been very smooth, and
this is compounded by the go team's insistence that they are eventually going to make go modules
required (they said this would happen in 1.12, and then didn't)
before they have worked out the kinks.

That said, go modules brings a couple of nice features to the go build environment:
1) Only one GOPATH directory. You don't need to change your environment variables 
every time you change projects.
2) Reproducible builds. This is the primary goal of modules, and generally has been
successful.

By default, go tries to detect whether to use go modules using
the following heuristic:
1. Is the current working directory inside the GOPATH environment variable? If so,
we are definitely NOT using go modules.
2. Else, does the current working directory, or a directory above it, have a go.mod file in
it. If so, we definitely ARE using go modules.
3. Otherwise, we are in limbo. Go version 1.11 handled this badly by just complaining.
Go 1.12+ handles it a little better, and allows you to install things with `go get`,
but you can't really do anything else.

Goradd is module aware, and will work whether you are using modules or not. Because
of the above behavior, the main thing you should be aware of is that whenever you
are building your application, or doing anything with the go command line tool,
you should do it from within the goradd-project directory. That way, the go tool
will be able to correctly figure out whether its in go module mode, and will be
able to find all the other parts of your application.

See the following for even more:
* [Go wiki on modules](https://github.com/golang/go/wiki/Modules)
* [Dave Cheney's Go Modules Article](https://dave.cheney.net/2018/07/14/taking-go-modules-for-a-spin)
