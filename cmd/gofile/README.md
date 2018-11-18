# gofile

This is a simple file and directory manipulation tool
primarily useful for build systems. By making a go version,
the build system can be cross-platform without having to know
what os you are running on. 

It is also module aware. Any directory can be represented as a module name,
followed by a subdirectory and the real path of the module will be substituted for 
the module name.

Usage:
```bash
gofile <command> <args...> [options] 
```

## Commands
### Copy
Copies a file or directory to another file or directory.

Usage:
```bash
gofile copy <src> <dest> [-x excludes...]
```
-x specifies names of files or directories you want to exclude from the source. This is useful when
expanding a directory using '*'.

### Generate

Runs go generate on the given file.

Usage:
```bash
gofile generate <sources...> [-x excludes...]

```

-x specifies names of files or directories you want to exclude from the source. This is useful when
expanding a directory using '*'.

### Mkdir

Creates the named directory if it does not exist. Sets it to be writable.

Usage:
```bash
gofile mkdir <dest>
```

### Remove

Deletes the named directories or files.

Usage:
```bash
gofile remove <dest...> [-x excludes...]
```

-x specifies names of files or directories you want to exclude from the destination. This is useful when
expanding a directory using '*'.


