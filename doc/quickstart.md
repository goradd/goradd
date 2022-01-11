# Quick Start
## Installation
### Setup Your Environment

1. Install go. Installation instructions for go are here: https://go.dev/doc/install
2. Create a new project directory *outside* of your go directory and change your working directory 
to the new directory. Note that go will create and use the "go" directory inside your home
directory for installing packages and go-based applications. 
Do NOT put your project in this directory.

For example, to create a "myproject" directory for your project in your HOME directory, do the following:

On Mac:
```
cd
mkdir myproject
cd myproject
```

On Windows:
```
cd %HOMEPATH%
mkdir myproject
cd myproject
```

### Install Goradd
1. Execute ```go install github.com/goradd/goradd@latest```
1. Execute ```goradd install```

You should now have a new directory in your current directory:
* goradd-project. This is where you will build out your project. You **can** put some
files outside of this path, but goradd will be placing its code generated files
inside of here.

You may also notice a number of executables that were installed 
that will be used by goradd to build your application.

If you have problems, see [Debugging Installation Problems](#debugging-installation-problems) below.

### Run the app
1. Change your working directory to the goradd-project directory that was created in the prior step. 
2. From the command line, run:
```go run -mod mod goradd-project/main```
You will see a number of messages about additional go packages being installed.
3. Once you see "Launching Server...", point your browser to the following URL. 
`http://localhost/goradd/`

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

#Debugging Installation Problems
Recent versions of the GO install process will create a "go" directory in your home
directory, and it should put the bin directory inside of that directory
into your PATH environment variable. This will allow go programs to be executed
from the command line.

If you have an an older version of GO, or something in this process fails, you may need to manually
set this process up.

First, make sure you can execute go by executing ```go version```. You should
see a go version string. If not, reinstall go.

Next, check your environment for a GOPATH variable. 
On Windows, execute ```set```. On a Mac or Linux, execute ```printenv```.
If you have a GOPATH variable, make sure the directory it points to exists
and has a bin directory. If it does not exist, make sure there is a go
directory inside you HOME directory (HOMEPATH on Windows), make sure there is
a bin directory inside it, and that both are writable, and the bin directory
has execute priviledges.

Finally, make sure the above go/bin directory is in your execution PATH.
Google "How to add to the execution Path for {Windows|Mac|Linux}" to get
info on that.

If everything is set up correctly, you should be able to run the following
commands:

```go install golang.org/x/example/hello```

and

```hello```

and then see the message:

```Hello, Go examples!```
