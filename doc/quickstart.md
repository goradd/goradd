# Quick Start
## Installation
### Setup Your Environment

If you have problems with the following directions, 
see [Debugging Installation Problems](#debugging-installation-problems) below.

1. Install go. Installation instructions for go are here: https://go.dev/doc/install. Note that
the installation process will put your local Go directory in your path so that the command-line Go programs
you create and install will be easily executable. The rest of these installation steps rely on that happening.
Be sure to restart any command-line or terminal windows after you install Go.
2. Create a new project directory *outside* of your go directory and change your working directory 
to the new directory. Note that go will create and use the "go" directory inside your home
directory for installing packages and go-based applications. Do NOT put your project in that directory.

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
2. Execute ```goradd install```

The _goradd install_ step will execute a series of commands to install required packages.
It will display what commands it is executing, and it may pause at times. Allow it to complete. 
You should now have a new directory in your current directory called _goradd-project_. 
This is where you will build out your go-based web project.

### Run the app
1. Change your working directory to the goradd-project directory that was created in the prior step. 
2. From the command line, run:
```go run goradd-project/main```
3. Once you see "Launching Server...", point your browser to the following URL. 
`http://localhost/goradd/`

If everything is working fine, you should see the Goradd startup screen. It will lead 
you through some additional configuration steps and get you started building your
application.

If you get an error message that looks something like this:
```
listen tcp :80: bind: address already in use
```
It means your computer already has a webserver running at port 80. If instead the message is a permission denied error, it means
your operating system has reserved low numbered ports for system use. By default, goradd applications run on the standard HTTP port 80.

In either case, you can use a different port as follows. Try using port 8000 to start with.
1. Execute:
```go run goradd-project/main -port XXXX```
Where XXXX is the port number you would like to use.
2. Go to the following address in your browser:
`http://localhost:XXXX/goradd/`
Where XXXX is the same number used in the prior step.

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

## Setting Up Debugging
When debugging, be sure to execute the main directory as a package, rather than just executing the main.go file. The reason for this is that by executing as a package, *all* the .go file in the main directory will be included. Executing just the main.go file will exclude the other .go files, and these other .go files have useful includes for the development process.

Here is an example launch.json configuration file for setting up debugging in VS Code. It assumes you have installed the VS Code Go extensions.
```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name":"Debug Goradd Program",
            "type":"go",
            "request":"launch",
            "mode":"debug",
            "program": "${workspaceFolder}/goradd-project/main"
        }
    ]
}
```

If you attempt to debug and you get an error about CGO and missing stdlib.h file, it means that the C to Go interface is enabled, but Go cannot find a C compiler. Either install a C compiler like gcc, or turn off CGO using this command:
```sh
go env -w CGO_ENABLED=0
```

# Modules
Goradd is module aware. Whenever you run goradd tools, it will look
in the nearest go.mod file to read the current module environment.

See the following for more info:
* [Go wiki on modules](https://github.com/golang/go/wiki/Modules)

# Debugging Installation Problems
Recent versions of the GO install process will create a "go" directory in your home
directory, and it should put the "bin" directory inside of that directory
into your PATH environment variable. This will allow go programs to be executed
from the command line.

If you have an an older version of GO, or something in this process fails, you may need to manually
set this process up.

First, make sure you can execute go by executing 
```
go version
```
You should
see a go version string. If not, reinstall go.

Next, check your GOPATH by executing 
```
go env GOPATH
```

Make sure there is a bin directory inside it, and that both are writable, 
and the bin directory has execute priviledges.

Finally, make sure the above go/bin directory is in your execution PATH.
Google "How to add to the execution Path for {Windows|Mac|Linux}" to get
info on that.

If everything is set up correctly, you should be able to run the following
commands:

```
go install golang.org/x/example/hello
```

and

```
hello
```

and then see the message:

```
Hello, Go examples!
```
